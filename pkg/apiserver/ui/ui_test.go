package ui

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/xyctruth/profiler/pkg/apiserver/ui/trace"

	"github.com/xyctruth/profiler/pkg/apiserver/ui/pprof"

	"github.com/xyctruth/profiler/pkg/storage"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
	"github.com/xyctruth/profiler/pkg/storage/badger"
)

func initProfileData(s storage.Store, t *testing.T) (string, string, string) {
	invalidId, err := s.SaveProfile([]byte{}, time.Second*10)
	require.Equal(t, nil, err)
	require.Equal(t, "0", invalidId)

	invalidId2, err := s.SaveProfile([]byte("haha"), time.Second*10)
	require.Equal(t, nil, err)
	require.Equal(t, "1", invalidId2)

	profileBytes, err := ioutil.ReadFile("../testdata/profile.gz")
	require.Equal(t, nil, err)
	id, err := s.SaveProfile(profileBytes, time.Second*10)
	require.Equal(t, nil, err)
	require.Equal(t, "2", id)
	return invalidId, invalidId2, id
}

func initTraceData(s storage.Store, t *testing.T) (string, string, string) {
	invalidId, err := s.SaveProfile([]byte{}, time.Second*10)
	require.Equal(t, nil, err)
	require.Equal(t, "0", invalidId)

	invalidId2, err := s.SaveProfile([]byte("haha"), time.Second*10)
	require.Equal(t, nil, err)
	require.Equal(t, "1", invalidId2)

	traceBytes, err := ioutil.ReadFile("../testdata/trace.gz")
	require.Equal(t, nil, err)
	id, err := s.SaveProfile(traceBytes, time.Second*10)
	require.Equal(t, nil, err)
	require.Equal(t, "2", id)
	return invalidId, invalidId2, id
}

func TestPProfServer(t *testing.T) {
	dir, err := ioutil.TempDir("./", "temp-*")
	require.Equal(t, nil, err)
	defer os.RemoveAll(dir)

	store := badger.NewStore(badger.DefaultOptions(dir))

	pprofServer := NewServer("/api/pprof/ui", store, 1*time.Minute, pprof.Driver)
	defer pprofServer.Exit()

	httpServer := httptest.NewServer(pprofServer.mux)
	defer httpServer.Close()

	// create httpexpect instance
	e := httpexpect.New(t, httpServer.URL)

	e.GET("/badUrl").
		Expect().
		Status(http.StatusBadRequest).Text().Equal("Invalid parameter\n")

	testPProfUI(e, store, t, pprofServer)
}

func testPProfUI(e *httpexpect.Expect, store storage.Store, t *testing.T, pprofServer *Server) {
	invalidId, invalidId2, id := initProfileData(store, t)

	e.GET("/api/pprof/ui/1999").
		Expect().
		Status(http.StatusNotFound).Text().Equal("Profile not found\n")

	e.GET(fmt.Sprintf("/api/pprof/ui/%s", invalidId)).
		Expect().
		Status(http.StatusInternalServerError).Text().Equal("failed to fetch any source profiles\n")

	e.GET(fmt.Sprintf("/api/pprof/ui/%s", invalidId2)).
		Expect().
		Status(http.StatusInternalServerError).Text().Equal("failed to fetch any source profiles\n")

	e.GET(fmt.Sprintf("/api/pprof/ui/%s", id)).
		Expect().
		Status(http.StatusOK).Header("Content-Type").Equal("text/html")

	e.GET(fmt.Sprintf("/api/pprof/ui/%s", id)).WithQuery("si", "alloc_objects").
		Expect().
		Status(http.StatusOK).Header("Content-Type").Equal("text/html")

	e.GET(fmt.Sprintf("/api/pprof/ui/%s/", id)).WithQuery("si", "alloc_space").
		Expect().
		Status(http.StatusOK).Header("Content-Type").Equal("text/html")

	e.GET(fmt.Sprintf("/api/pprof/ui/%s/", id)).WithQuery("si", "ppp").
		Expect().
		Status(http.StatusBadRequest).Text().Equal("sample_index \"ppp\" must be one of: [alloc_objects alloc_space inuse_objects inuse_space]\n")

	e.GET(fmt.Sprintf("/api/pprof/ui/%s/top", id)).WithQuery("si", "alloc_space").
		Expect().
		Status(http.StatusOK).Header("Content-Type").Equal("text/html")

	pprofServer.gc()

	e.GET(fmt.Sprintf("/api/pprof/ui/%s/top", id)).WithQuery("si", "alloc_space").
		Expect().
		Status(http.StatusOK).Header("Content-Type").Equal("text/html")

	e.GET(fmt.Sprintf("/api/pprof/ui/%s/toperror", id)).WithQuery("si", "alloc_space").
		Expect().
		Status(http.StatusOK).Header("Content-Type").Equal("text/html")

}

func TestTraceServer(t *testing.T) {
	dir, err := ioutil.TempDir("./", "temp-*")
	require.Equal(t, nil, err)
	defer os.RemoveAll(dir)

	store := badger.NewStore(badger.DefaultOptions(dir))

	traceServer := NewServer("/api/trace/ui", store, 1*time.Minute, trace.Driver)
	defer traceServer.Exit()

	httpServer := httptest.NewServer(traceServer.mux)
	defer httpServer.Close()

	// create httpexpect instance
	e := httpexpect.New(t, httpServer.URL)

	e.GET("/badUrl").
		Expect().
		Status(http.StatusBadRequest).Text().Equal("Invalid parameter\n")

	testTraceUI(e, store, t, traceServer)
}

func testTraceUI(e *httpexpect.Expect, store storage.Store, t *testing.T, server *Server) {
	invalidId, invalidId2, id := initTraceData(store, t)

	e.GET("/api/trace/ui/1999").
		Expect().
		Status(http.StatusNotFound).Text().Equal("Profile not found\n")

	e.GET(fmt.Sprintf("/api/trace/ui/%s", invalidId)).
		Expect().
		Status(http.StatusInternalServerError).Text().Equal("EOF\n")

	e.GET(fmt.Sprintf("/api/trace/ui/%s", invalidId2)).
		Expect().
		Status(http.StatusInternalServerError).Text().Equal("unexpected EOF\n")

	e.GET(fmt.Sprintf("/api/trace/ui/%s", id)).
		Expect().
		Status(http.StatusOK).Header("Content-Type").Equal("text/html; charset=utf-8")

	e.GET(fmt.Sprintf("/api/trace/ui/%s", id)).
		Expect().
		Status(http.StatusOK).Header("Content-Type").Equal("text/html; charset=utf-8")

	server.gc()

	e.GET(fmt.Sprintf("/api/trace/ui/%s", id)).
		Expect().
		Status(http.StatusOK).Header("Content-Type").Equal("text/html; charset=utf-8")

}
