package api_server

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/google/pprof/driver"
)

type pprofServer struct {
	exits        map[string]struct{}
	mux          *http.ServeMux
	mu           sync.Mutex
	registerPath string
	webPath      string
}

func newPprofServer(registerPath, webPath string) *pprofServer {
	s := &pprofServer{
		mux:          http.NewServeMux(),
		exits:        make(map[string]struct{}, 0),
		registerPath: registerPath,
		webPath:      webPath,
	}

	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		reg, _ := regexp.Compile(`(` + webPath + `|/[a-zA-Z]+)`)
		path := registerPath + reg.ReplaceAllString(r.URL.Path, "")
		sampleType := r.URL.Query().Get("si")
		s.redirectProfile(w, r, path, sampleType)
	})

	return s
}

func (s *pprofServer) redirectProfile(w http.ResponseWriter, r *http.Request, path, sampleType string) {
	http.Redirect(w, r, path+"?si="+sampleType, http.StatusSeeOther)
}

func (s *pprofServer) register(w http.ResponseWriter, r *http.Request, data []byte, id, sampleType string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	curPath := path.Join(s.webPath, id) + "/"
	if _, ok := s.exits[id]; ok {
		s.redirectProfile(w, r, curPath, sampleType)
		return
	}

	s.exits[id] = struct{}{}

	filepath := path.Join(os.TempDir(), id)
	err := ioutil.WriteFile(filepath, data, 0600)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	flags := &pprofFlags{
		args: []string{"-http=localhost:0", "-no_browser", filepath},
	}

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
	if err := driver.PProf(options); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.redirectProfile(w, r, curPath, sampleType)
}

func (s *pprofServer) web(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mux.ServeHTTP(w, r)
}

type pprofFlags struct {
	args  []string
	s     flag.FlagSet
	usage []string
}

// Bool implements the plugin.FlagSet interface.
func (p *pprofFlags) Bool(o string, d bool, c string) *bool {
	return p.s.Bool(o, d, c)
}

// Int implements the plugin.FlagSet interface.
func (p *pprofFlags) Int(o string, d int, c string) *int {
	return p.s.Int(o, d, c)
}

// Float64 implements the plugin.FlagSet interface.
func (p *pprofFlags) Float64(o string, d float64, c string) *float64 {
	return p.s.Float64(o, d, c)
}

// String implements the plugin.FlagSet interface.
func (p *pprofFlags) String(o, d, c string) *string {
	return p.s.String(o, d, c)
}

// BoolVar implements the plugin.FlagSet interface.
func (p *pprofFlags) BoolVar(b *bool, o string, d bool, c string) {
	p.s.BoolVar(b, o, d, c)
}

// IntVar implements the plugin.FlagSet interface.
func (p *pprofFlags) IntVar(i *int, o string, d int, c string) {
	p.s.IntVar(i, o, d, c)
}

// Float64Var implements the plugin.FlagSet interface.
// the value of the flag.
func (p *pprofFlags) Float64Var(f *float64, o string, d float64, c string) {
	p.s.Float64Var(f, o, d, c)
}

// StringVar implements the plugin.FlagSet interface.
func (p *pprofFlags) StringVar(s *string, o, d, c string) {
	p.s.StringVar(s, o, d, c)
}

// StringList implements the plugin.FlagSet interface.
func (p *pprofFlags) StringList(o, d, c string) *[]*string {
	return &[]*string{p.s.String(o, d, c)}
}

// AddExtraUsage implements the plugin.FlagSet interface.
func (p *pprofFlags) AddExtraUsage(eu string) {
	p.usage = append(p.usage, eu)
}

// ExtraUsage implements the plugin.FlagSet interface.
func (p *pprofFlags) ExtraUsage() string {
	return strings.Join(p.usage, "\n")
}

// Parse implements the plugin.FlagSet interface.
func (p *pprofFlags) Parse(usage func()) []string {
	p.s.Usage = usage
	p.s.Parse(p.args)
	args := p.s.Args()
	if len(args) == 0 {
		usage()
	}
	return args
}
