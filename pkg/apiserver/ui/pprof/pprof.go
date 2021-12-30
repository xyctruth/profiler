package pprof

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/google/pprof/driver"
)

func Driver(basePath string, mux *http.ServeMux, id string, data []byte) error {
	filepath := path.Join(os.TempDir(), id)
	if err := ioutil.WriteFile(filepath, data, 0600); err != nil {
		return err
	}

	flags := &flags{
		args: []string{"-http=localhost:0", "-no_browser", filepath},
	}

	curPath := path.Join(basePath, id) + "/"
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
				mux.Handle(joinedPattern, handler)
			}
			return nil
		},
	}
	if err := driver.PProf(options); err != nil {
		return err
	}
	return nil
}
