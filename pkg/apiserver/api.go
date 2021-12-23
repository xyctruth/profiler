package apiserver

import (
	"context"
	"errors"
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
	pprofPath := "/api/pprof/ui"

	apiServer := &APIServer{
		store: store,
		pprof: newPprofServer(pprofPath, store),
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

	// register pprof page
	router.Use(HandleCors).GET(pprofPath+"/*any", apiServer.webPprof)

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
	c.JSON(200, jobs)
}

func (s *APIServer) listGroupSampleTypes(c *gin.Context) {
	jobs, err := s.store.ListGroupSampleType()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(200, jobs)
}

func (s *APIServer) listProfileMeta(c *gin.Context) {
	sampleType := c.Param("sample_type")

	if c.Query("start_time") == "" {
		c.String(http.StatusBadRequest, "start_time is empty")
		return
	}

	if c.Query("end_time") == "" {
		c.String(http.StatusBadRequest, "end_time is empty")
		return
	}

	startTime, err := time.Parse(time.RFC3339, c.Query("start_time"))
	if err != nil {
		c.String(http.StatusBadRequest, "%s ,%s", "The time format must be RFC3339", err.Error())
		return
	}
	endTime, err := time.Parse(time.RFC3339, c.Query("end_time"))
	if err != nil {
		c.String(http.StatusBadRequest, "%s ,%s", "The time format must be RFC3339", err.Error())
		return
	}

	req := struct {
		Targets []string `json:"targets" form:"targets"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	req.Targets = utils.RemoveDuplicateElement(req.Targets)
	metas, err := s.store.ListProfileMeta(sampleType, req.Targets, startTime, endTime)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(200, metas)
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

	p, err := profile.ParseData(data)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Writer.Header().Set("Content-Type", "application/vnd.google.protobuf+gzip")
	c.Writer.Header().Set("Content-Disposition", "attachment;filename=profile.pb.gz")
	err = p.Write(c.Writer)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}

func (s *APIServer) webPprof(c *gin.Context) {
	c.Request.URL.RawQuery = removePrefixSampleType(c.Request.URL.RawQuery)
	s.pprof.web(c.Writer, c.Request)
}
