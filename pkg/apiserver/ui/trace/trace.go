package trace

import (
	"net/http"
	"path"

	"github.com/xyctruth/profiler/pkg/internal/v1175/traceui"
)

func Driver(basePath string, mux *http.ServeMux, id string, data []byte) error {
	ui, err := traceui.NewUI(data)
	if err != nil {
		return err
	}

	curPath := path.Join(basePath, id) + "/"
	for pattern, handler := range ui.Handlers {
		var joinedPattern string
		if pattern == "/" {
			joinedPattern = curPath
		} else {
			joinedPattern = path.Join(curPath, pattern)
		}
		mux.Handle(joinedPattern, handler)
	}
	return nil
}
