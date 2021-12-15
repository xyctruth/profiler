package api_server

import (
	"bytes"
	"context"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/google/pprof/driver"

	"github.com/xyctruth/profiler/pkg/storage"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/google/pprof/profile"
)

type ApiServer struct {
	store       storage.Store
	router      *gin.Engine
	host        string
	srv         *http.Server
	pprofExits  map[string]struct{}
	pprofMux    *http.ServeMux
	pprofPrefix string
	pprofMu     sync.Mutex
}

func NewApiServer(addr string, store storage.Store) *ApiServer {
	apiServer := &ApiServer{
		store:       store,
		pprofPrefix: "/web/pprof",
		pprofMux:    http.NewServeMux(),
		pprofExits:  make(map[string]struct{}, 0),
	}

	router := gin.Default()
	router.Use(HandleCors).GET("/api/targets", apiServer.listTarget)
	router.Use(HandleCors).GET("/api/sample_types", apiServer.listSampleTypes)
	router.Use(HandleCors).GET("/api/group_sample_types", apiServer.listGroupSampleTypes)
	router.Use(HandleCors).GET("/api/profile/:id", apiServer.getProfile)
	router.Use(HandleCors).GET("/api/profile_meta/:sample_type", apiServer.listProfileMeta)
	router.Use(HandleCors).GET("/web/profile/:id", apiServer.webProfile)
	router.Use(HandleCors).GET(apiServer.pprofPrefix+"/*any", apiServer.webPprof)

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	apiServer.router = router
	apiServer.srv = srv
	return apiServer
}

func (apiServer *ApiServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := apiServer.srv.Shutdown(ctx); err != nil {
		log.Fatal("api server forced to shutdown:", err)
	}
	log.Info("api server exit ")
}

func (apiServer *ApiServer) Run() {
	if err := apiServer.srv.ListenAndServe(); err != nil {
		log.Fatal("api server listen: %s\n", err)
	}
}

func (apiServer *ApiServer) listTarget(c *gin.Context) {
	jobs, err := apiServer.store.ListTarget()
	if err != nil {
		c.Error(err)
	}
	c.JSON(200, jobs)
}

func (apiServer *ApiServer) listProfileMeta(c *gin.Context) {
	sampleType := c.Param("sample_type")

	startTime, err := time.Parse(time.RFC3339, c.Query("start_time"))
	if err != nil {
		c.Error(err)
		return
	}
	endTime, err := time.Parse(time.RFC3339, c.Query("end_time"))
	if err != nil {
		c.Error(err)
		return
	}

	req := struct {
		Targets []string `json:"targets" form:"targets"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		c.Error(err)
	}

	metas, err := apiServer.store.ListProfileMeta(sampleType, req.Targets, startTime, endTime)
	if err != nil {
		c.Error(err)
	}
	c.JSON(200, metas)
}

func (apiServer *ApiServer) listSampleTypes(c *gin.Context) {
	jobs, err := apiServer.store.ListSampleType()
	if err != nil {
		c.Error(err)
	}
	c.JSON(200, jobs)
}

func (apiServer *ApiServer) listGroupSampleTypes(c *gin.Context) {
	jobs, err := apiServer.store.ListGroupSampleType()
	if err != nil {
		c.Error(err)
	}
	c.JSON(200, jobs)
}

func (apiServer *ApiServer) getProfile(c *gin.Context) {
	id := c.Param("id")
	data, err := apiServer.store.GetProfile(id)
	if err != nil {
		c.Error(err)
	}

	p, _ := profile.ParseData(data)
	c.Writer.Header().Set("Content-Type", "application/vnd.google.protobuf+gzip")
	c.Writer.Header().Set("Content-Disposition", "attachment;filename=profile.pb.gz")
	p.Write(c.Writer)
}

func (apiServer *ApiServer) webProfile(c *gin.Context) {
	apiServer.pprofMu.Lock()
	defer apiServer.pprofMu.Unlock()
	id := c.Param("id")
	curPath := path.Join(apiServer.pprofPrefix, id) + "/"

	if _, ok := apiServer.pprofExits[id]; ok {
		http.Redirect(c.Writer, c.Request, curPath, http.StatusSeeOther)
		return
	}

	apiServer.pprofExits[id] = struct{}{}

	data, err := apiServer.store.GetProfile(id)
	if err != nil {
		c.Error(err)
		return
	}
	p, _ := profile.ParseData(data)
	b := &bytes.Buffer{}
	p.Write(b)
	filepath := path.Join(os.TempDir(), id)
	err = ioutil.WriteFile(filepath, b.Bytes(), 0600)
	if err != nil {
		c.Error(err)
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
				apiServer.pprofMux.Handle(joinedPattern, handler)
			}
			return nil
		},
	}
	if err := driver.PProf(options); err != nil {
		c.Error(err)
		return
	}
	http.Redirect(c.Writer, c.Request, curPath, http.StatusSeeOther)
}

func (apiServer *ApiServer) webPprof(c *gin.Context) {
	if apiServer.pprofMux == nil {
		http.Error(c.Writer, "must upload profile first", http.StatusInternalServerError)
		return
	}
	apiServer.pprofMux.ServeHTTP(c.Writer, c.Request)
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
