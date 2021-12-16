package apiserver

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/xyctruth/profiler/pkg/utils"

	"github.com/xyctruth/profiler/pkg/storage"

	"github.com/google/pprof/driver"
)

type pprofServer struct {
	exits   map[string]struct{}
	mux     *http.ServeMux
	mu      sync.Mutex
	webPath string
	store   storage.Store
}

func newPprofServer(webPath string, store storage.Store) *pprofServer {
	s := &pprofServer{
		mux:     http.NewServeMux(),
		exits:   make(map[string]struct{}),
		webPath: webPath,
		store:   store,
	}
	s.mux.HandleFunc("/", s.register)
	return s
}

func (s *pprofServer) web(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *pprofServer) register(w http.ResponseWriter, r *http.Request) {
	sampleType := utils.RemoveSampleTypePrefix(r.URL.Query().Get("si"))

	id := r.URL.Query().Get(("id"))
	if id == "" {

	}

	data, err := s.store.GetProfile(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	curPath := path.Join(s.webPath, id) + "/"
	if _, ok := s.exits[id]; ok {
		http.Redirect(w, r, curPath+"?si="+sampleType, http.StatusSeeOther)
		return
	}

	s.exits[id] = struct{}{}

	filepath := path.Join(os.TempDir(), id)
	err = ioutil.WriteFile(filepath, data, 0600)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	flags := &pprofFlags{
		args: []string{"-http=localhost:0", "-no_browser", filepath},
	}

	options := &driver.Options{
		Flagset: flags,
		HTTPServer: func(args *driver.HTTPServerArgs) error {
			for pattern, handler := range args.Handlers {
				var joinedPattern string
				if pattern == "/" {
					joinedPattern = curPath
				} else {
					joinedPattern = path.Join(curPath, pattern)
				}
				s.mux.Handle(joinedPattern, handler)
			}
			return nil
		},
	}
	if err := driver.PProf(options); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, curPath+"?si="+sampleType, http.StatusSeeOther)
}
