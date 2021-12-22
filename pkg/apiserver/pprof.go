package apiserver

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/google/pprof/driver"
	"github.com/xyctruth/profiler/pkg/storage"
)

type pprofServer struct {
	mux      *http.ServeMux
	mu       sync.Mutex
	basePath string
	store    storage.Store
}

func newPprofServer(basePath string, store storage.Store) *pprofServer {
	s := &pprofServer{
		mux:      http.NewServeMux(),
		basePath: basePath,
		store:    store,
	}
	s.mux.HandleFunc("/", s.register)
	return s
}

func (s *pprofServer) web(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *pprofServer) register(w http.ResponseWriter, r *http.Request) {
	id := extractProfileID(r.URL.Path)
	if id == "" {
		http.Error(w, "Invalid parameter", http.StatusBadRequest)
		return
	}
	data, err := s.store.GetProfile(id)
	if err != nil {
		if err == storage.ErrProfileNotFound {
			http.Error(w, "Profile not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	filepath := path.Join(os.TempDir(), id)
	err = ioutil.WriteFile(filepath, data, 0600)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	flags := &pprofFlags{
		args: []string{"-http=localhost:0", "-no_browser", filepath},
	}

	curPath := path.Join(s.basePath, id) + "/"
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
	http.Redirect(w, r, r.URL.Path+"?"+r.URL.RawQuery, http.StatusSeeOther)
}
