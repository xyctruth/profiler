package pprof

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	"github.com/google/pprof/driver"
	"github.com/xyctruth/profiler/pkg/storage"
	"github.com/xyctruth/profiler/pkg/utils"
)

type Server struct {
	mux      *http.ServeMux
	mu       sync.Mutex
	basePath string
	store    storage.Store
	exitChan chan struct{}
}

func NewServer(basePath string, store storage.Store) *Server {
	s := &Server{
		mux:      http.NewServeMux(),
		basePath: basePath,
		store:    store,
		exitChan: make(chan struct{}),
	}
	s.mux.HandleFunc("/", s.register)

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-s.exitChan:
				return
			case <-ticker.C:
				s.gc()
			}
		}
	}()
	return s
}

func (s *Server) Exit() {
	close(s.exitChan)
}

func (s *Server) gc() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/", s.register)
}

func (s *Server) Web(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractProfileID(r.URL.Path)
	if id == "" {
		http.Error(w, "Invalid parameter", http.StatusBadRequest)
		return
	}
	data, err := s.store.GetProfile(id)
	if err != nil {
		if errors.Is(err, storage.ErrProfileNotFound) {
			http.Error(w, "Profile not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	filepath := path.Join(os.TempDir(), id)
	if err = ioutil.WriteFile(filepath, data, 0600); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	flags := &flags{
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
	if err = driver.PProf(options); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, r.URL.Path+"?"+r.URL.RawQuery, http.StatusSeeOther)
}
