package trace

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/xyctruth/profiler/pkg/internal/v1175/traceui"
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

	buf := bytes.NewBuffer(data)
	gzipReader, err := gzip.NewReader(buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer gzipReader.Close()
	b, err := ioutil.ReadAll(gzipReader)
	if err != nil && !strings.Contains(err.Error(), "unexpected EOF") {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	ui := traceui.NewTraceUI(b)

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
