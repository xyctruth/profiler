package apiserver

import (
	"time"

	"github.com/xyctruth/profiler/pkg/storage"
)

type Options struct {
	Addr       string
	GCInternal time.Duration
	Store      storage.Store
}

func DefaultOptions(store storage.Store) Options {
	return Options{
		Store:      store,
		Addr:       ":8080",
		GCInternal: 2 * time.Minute,
	}
}

func (opt Options) WithAddr(addr string) Options {
	opt.Addr = addr
	return opt
}

func (opt Options) WithGCInternal(internal time.Duration) Options {
	opt.GCInternal = internal
	return opt
}
