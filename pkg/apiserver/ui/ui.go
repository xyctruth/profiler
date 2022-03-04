package ui

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/xyctruth/profiler/pkg/storage"
	"github.com/xyctruth/profiler/pkg/utils"
)

type Driver func(basePath string, mux *http.ServeMux, id string, data []byte) error

type Server struct {
	cache    map[string]struct{}
	mux      *http.ServeMux
	mu       sync.Mutex
	basePath string
	store    storage.Store
	exitChan chan struct{}
	drive    Driver
}

func NewServer(basePath string, store storage.Store, gcInternal time.Duration, drive Driver) *Server {
	s := &Server{
		mux:      http.NewServeMux(),
		basePath: basePath,
		store:    store,
		exitChan: make(chan struct{}),
		cache:    make(map[string]struct{}),
		drive:    drive,
	}
	s.mux.HandleFunc("/", s.register)

	go func() {
		ticker := time.NewTicker(gcInternal)
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
	s.cache = make(map[string]struct{})
}

func (s *Server) Web(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := utils.ExtractProfileID(r.URL.Path)
	if id == "" {
		http.Error(w, "Invalid parameter", http.StatusBadRequest)
		return
	}

	if _, ok := s.cache[id]; ok {
		http.Redirect(w, r, r.URL.Path+"?"+r.URL.RawQuery, http.StatusSeeOther)
		return
	}

	_, data, err := s.store.GetProfile(id)
	if err != nil {
		if errors.Is(err, storage.ErrProfileNotFound) {
			http.Error(w, "Profile not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.drive(s.basePath, s.mux, id, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.cache[id] = struct{}{}
	http.Redirect(w, r, r.URL.Path+"?"+r.URL.RawQuery, http.StatusSeeOther)
}
