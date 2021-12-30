package badger

import "time"

type Options struct {
	Path       string
	GCInternal time.Duration
}

func DefaultOptions(path string) Options {
	return Options{
		Path:       path,
		GCInternal: 5 * time.Minute,
	}
}

func (opt Options) WithGCInternal(internal time.Duration) Options {
	opt.GCInternal = internal
	return opt
}
