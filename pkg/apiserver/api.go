package apiserver

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/xyctruth/profiler/pkg/apiserver/ui"
	"github.com/xyctruth/profiler/pkg/apiserver/ui/pprof"
	"github.com/xyctruth/profiler/pkg/apiserver/ui/trace"
	"github.com/xyctruth/profiler/pkg/storage"
	"github.com/xyctruth/profiler/pkg/utils"
)

type APIServer struct {
	opt    Options
	store  storage.Store
	router *gin.Engine
	srv    *http.Server
	pprof  *ui.Server
	trace  *ui.Server
}

func NewAPIServer(opt Options) *APIServer {
	pprofPath := "/api/pprof/ui"
	tracePath := "/api/trace/ui"

	apiServer := &APIServer{
		opt:   opt,
		store: opt.Store,
		pprof: ui.NewServer(pprofPath, opt.Store, opt.GCInternal, pprof.Driver),
		trace: ui.NewServer(tracePath, opt.Store, opt.GCInternal, trace.Driver),
	}

	router := gin.Default()
	router.GET("/api/healthz", func(c *gin.Context) {
		c.String(200, "I'm fine")
	})
	router.Use(HandleCors).GET("/api/targets", apiServer.listTarget)
	router.Use(HandleCors).GET("/api/sample_types", apiServer.listSampleTypes)
	router.Use(HandleCors).GET("/api/group_sample_types", apiServer.listGroupSampleTypes)
	router.Use(HandleCors).GET("/api/profile_meta/:sample_type", apiServer.listProfileMeta)
	router.Use(HandleCors).GET("/api/profile/:id", apiServer.getProfile)
	router.Use(HandleCors).GET("/api/trace/:id", apiServer.getTrace)

	// register pprof page
	router.Use(HandleCors).GET(pprofPath+"/*any", apiServer.webPProf)
	// register trace page
	router.Use(HandleCors).GET(tracePath+"/*any", apiServer.webTrace)

	srv := &http.Server{
		Addr:    opt.Addr,
		Handler: router,
	}

	apiServer.router = router
	apiServer.srv = srv
	return apiServer
}

func (s *APIServer) Stop() {
	s.pprof.Exit()
	s.trace.Exit()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal("api server forced to shutdown:", err)
	}
	log.Info("api server exit ")
}

func (s *APIServer) Run() {
	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Info("api server close")
				return
			}
			log.Fatal("api server listen: ", err)
		}
	}()
}

func (s *APIServer) listTarget(c *gin.Context) {
	jobs, err := s.store.ListTarget()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, jobs)
}

func (s *APIServer) listSampleTypes(c *gin.Context) {
	jobs, err := s.store.ListSampleType()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, jobs)
}

func (s *APIServer) listGroupSampleTypes(c *gin.Context) {
	jobs, err := s.store.ListGroupSampleType()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, jobs)
}

func (s *APIServer) listProfileMeta(c *gin.Context) {
	var startTime, endTime time.Time
	var err error

	sampleType := c.Param("sample_type")

	if c.Query("start_time") == "" || c.Query("end_time") == "" {
		c.String(http.StatusBadRequest, "start_time or end_time is empty")
		return
	}

	if startTime, err = time.Parse(time.RFC3339, c.Query("start_time")); err != nil {
		c.String(http.StatusBadRequest, "%s ,%s", "The time format must be RFC3339", err.Error())
		return
	}
	if endTime, err = time.Parse(time.RFC3339, c.Query("end_time")); err != nil {
		c.String(http.StatusBadRequest, "%s ,%s", "The time format must be RFC3339", err.Error())
		return
	}

	req := struct {
		Targets []string        `json:"targets" form:"targets"`
		Labels  []storage.Label `json:"labels" form:"labels"`
	}{}

	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	req.Targets = utils.RemoveDuplicateElement(req.Targets)
	metas, err := s.store.ListProfileMeta(sampleType, req.Targets, req.Labels, startTime, endTime)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, metas)
}

func (s *APIServer) getProfile(c *gin.Context) {
	id := c.Param("id")
	data, err := s.store.GetProfile(id)
	if err != nil {
		if errors.Is(err, storage.ErrProfileNotFound) {
			c.String(http.StatusNotFound, "Profile not found")
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Writer.Header().Set("Content-Encoding", "gzip")
	c.Writer.Header().Set("Content-Disposition", "attachment;filename=profile_"+id+".pb.gz")
	c.Data(200, "application/vnd.google.protobuf+gzip", data)
}

func (s *APIServer) getTrace(c *gin.Context) {
	id := c.Param("id")
	data, err := s.store.GetProfile(id)
	if err != nil {
		if errors.Is(err, storage.ErrProfileNotFound) {
			c.String(http.StatusNotFound, "Profile not found")
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	buf := bytes.NewBuffer(data)
	gzipReader, err := gzip.NewReader(buf)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer gzipReader.Close()
	b, err := ioutil.ReadAll(gzipReader)
	if err != nil && !strings.Contains(err.Error(), "unexpected EOF") {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Writer.Header().Set("Content-Disposition", "attachment;filename=trace_"+id+".out")
	c.Data(200, "application/octet-stream", b)
}

func (s *APIServer) webPProf(c *gin.Context) {
	c.Request.URL.RawQuery = utils.RemovePrefixSampleType(c.Request.URL.RawQuery)
	s.pprof.Web(c.Writer, c.Request)
}

func (s *APIServer) webTrace(c *gin.Context) {
	c.Request.URL.RawQuery = utils.RemovePrefixSampleType(c.Request.URL.RawQuery)
	s.trace.Web(c.Writer, c.Request)
}
