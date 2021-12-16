package apiserver

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/pprof/profile"
	log "github.com/sirupsen/logrus"
	"github.com/xyctruth/profiler/pkg/storage"
	"github.com/xyctruth/profiler/pkg/utils"
)

type APIServer struct {
	store  storage.Store
	router *gin.Engine
	srv    *http.Server
	pprof  *pprofServer
}

func NewAPIServer(addr string, store storage.Store) *APIServer {
	registerPprofPath, webPprofPath := "/api/pprof/register", "/api/pprof/ui"

	apiServer := &APIServer{
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

func (s *APIServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal("api server forced to shutdown:", err)
	}
	log.Info("api server exit ")
}

func (s *APIServer) Run() {
	if err := s.srv.ListenAndServe(); err != nil {
		log.Fatal("api server listen: ", err)
	}
}

func (s *APIServer) listTarget(c *gin.Context) {
	jobs, err := s.store.ListTarget()
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, jobs)
}

func (s *APIServer) listProfileMeta(c *gin.Context) {
	sampleType := c.Param("sample_type")

	startTime, err := time.Parse(time.RFC3339, c.Query("start_time"))
	if err != nil {
		c.JSON(500, err)
		return
	}
	endTime, err := time.Parse(time.RFC3339, c.Query("end_time"))
	if err != nil {
		c.JSON(500, err)
		return
	}

	req := struct {
		Targets []string `json:"targets" form:"targets"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(500, err)
		return
	}

	metas, err := s.store.ListProfileMeta(sampleType, req.Targets, startTime, endTime)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, metas)
}

func (s *APIServer) listSampleTypes(c *gin.Context) {
	jobs, err := s.store.ListSampleType()
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, jobs)
}

func (s *APIServer) listGroupSampleTypes(c *gin.Context) {
	jobs, err := s.store.ListGroupSampleType()
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, jobs)
}

func (s *APIServer) getProfile(c *gin.Context) {
	id := c.Param("id")
	data, err := s.store.GetProfile(id)
	if err != nil {
		c.JSON(500, err)
		return
	}

	p, _ := profile.ParseData(data)
	c.Writer.Header().Set("Content-Type", "application/vnd.google.protobuf+gzip")
	c.Writer.Header().Set("Content-Disposition", "attachment;filename=profile.pb.gz")
	err = p.Write(c.Writer)
	if err != nil {
		c.JSON(500, err)
	}
}
func (s *APIServer) registerPprof(c *gin.Context) {
	sampleType := utils.RemoveSampleTypePrefix(c.Query("si"))
	id := c.Param("id")
	data, err := s.store.GetProfile(id)
	if err != nil {
		c.JSON(500, err)
		return
	}
	s.pprof.register(c.Writer, c.Request, data, id, sampleType)
}

func (s *APIServer) webPprof(c *gin.Context) {
	s.pprof.web(c.Writer, c.Request)
}
