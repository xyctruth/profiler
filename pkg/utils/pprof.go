package utils

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/felixge/fgprof"
)

func RegisterPProf() {
	go func() {
		http.DefaultServeMux.Handle("/debug/fgprof", fgprof.Handler())
		err := http.ListenAndServe(":9000", nil)
		if err != nil {
			panic(err)
		}
	}()
}
