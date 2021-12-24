package trace

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/xyctruth/profiler/pkg/internal/v1175/traceui"
	"github.com/xyctruth/profiler/pkg/storage"
	"github.com/xyctruth/profiler/pkg/utils"
)

type Server struct {
	mux      *http.ServeMux
	mu       sync.Mutex
	basePath string
	store    storage.Store
}

func NewServer(basePath string, store storage.Store) *Server {
	s := &Server{
		mux:      http.NewServeMux(),
		basePath: basePath,
		store:    store,
	}
	s.mux.HandleFunc("/", s.register)
	return s
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

	ui := traceui.NewTraceUI(filepath)

	curPath := path.Join(s.basePath, id) + "/"
	for pattern, handler := range ui.Handlers {
		var joinedPattern string
		if pattern == "/" {
			joinedPattern = curPath
		} else {
			joinedPattern = path.Join(curPath, pattern)
		}
		s.mux.Handle(joinedPattern, handler)
	}

	http.Redirect(w, r, r.URL.Path+"?"+r.URL.RawQuery, http.StatusSeeOther)
}
