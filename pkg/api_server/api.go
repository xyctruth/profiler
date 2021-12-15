package api_server

import (
	"context"
	"net/http"
	"time"

	"github.com/xyctruth/profiler/pkg/utils"

	"github.com/xyctruth/profiler/pkg/storage"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/google/pprof/profile"
)

type ApiServer struct {
	store  storage.Store
	router *gin.Engine
	host   string
	srv    *http.Server
	pprof  *pprofServer
}

func NewApiServer(addr string, store storage.Store) *ApiServer {
	registerPprofPath, webPprofPath := "/pprof/register", "/pprof/web"

	apiServer := &ApiServer{
		store: store,
		pprof: newPprofServer(registerPprofPath, webPprofPath),
	}

	router := gin.Default()
	router.Use(HandleCors).GET("/api/targets", apiServer.listTarget)
	router.Use(HandleCors).GET("/api/sample_types", apiServer.listSampleTypes)
	router.Use(HandleCors).GET("/api/group_sample_types", apiServer.listGroupSampleTypes)
	router.Use(HandleCors).GET("/api/profile/:id", apiServer.getProfile)
	router.Use(HandleCors).GET("/api/profile_meta/:sample_type", apiServer.listProfileMeta)

	// show pprof page
	router.Use(HandleCors).GET(registerPprofPath+"/:id", apiServer.registerPprof)
	router.Use(HandleCors).GET(webPprofPath+"/*any", apiServer.webPprof)

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	apiServer.router = router
	apiServer.srv = srv
	return apiServer
}

func (s *ApiServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal("api server forced to shutdown:", err)
	}
	log.Info("api server exit ")
}

func (s *ApiServer) Run() {
	if err := s.srv.ListenAndServe(); err != nil {
		log.Fatal("api server listen: %s\n", err)
	}
}

func (s *ApiServer) listTarget(c *gin.Context) {
	jobs, err := s.store.ListTarget()
	if err != nil {
		c.Error(err)
	}
	c.JSON(200, jobs)
}

func (s *ApiServer) listProfileMeta(c *gin.Context) {
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

	metas, err := s.store.ListProfileMeta(sampleType, req.Targets, startTime, endTime)
	if err != nil {
		c.Error(err)
	}
	c.JSON(200, metas)
}

func (s *ApiServer) listSampleTypes(c *gin.Context) {
	jobs, err := s.store.ListSampleType()
	if err != nil {
		c.Error(err)
	}
	c.JSON(200, jobs)
}

func (s *ApiServer) listGroupSampleTypes(c *gin.Context) {
	jobs, err := s.store.ListGroupSampleType()
	if err != nil {
		c.Error(err)
	}
	c.JSON(200, jobs)
}

func (s *ApiServer) getProfile(c *gin.Context) {
	id := c.Param("id")
	data, err := s.store.GetProfile(id)
	if err != nil {
		c.Error(err)
	}

	p, _ := profile.ParseData(data)
	c.Writer.Header().Set("Content-Type", "application/vnd.google.protobuf+gzip")
	c.Writer.Header().Set("Content-Disposition", "attachment;filename=profile.pb.gz")
	p.Write(c.Writer)
}
func (s *ApiServer) registerPprof(c *gin.Context) {
	sampleType := utils.RemoveSampleTypePrefix(c.Query("si"))
	id := c.Param("id")
	data, err := s.store.GetProfile(id)
	if err != nil {
		c.Error(err)
		return
	}
	s.pprof.register(c.Writer, c.Request, data, id, sampleType)
}

func (s *ApiServer) webPprof(c *gin.Context) {
	s.pprof.web(c.Writer, c.Request)
}
