package trace

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/xyctruth/profiler/pkg/internal/v1175/traceui"
)

func Driver(basePath string, mux *http.ServeMux, id string, data []byte) error {
	buf := bytes.NewBuffer(data)
	gzipReader, err := gzip.NewReader(buf)
	if err != nil {
		return err
	}
	defer gzipReader.Close()
	b, err := ioutil.ReadAll(gzipReader)
	if err != nil && !strings.Contains(err.Error(), "unexpected EOF") {
		return err
	}

	ui, err := traceui.NewUI(b)
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
