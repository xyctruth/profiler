package trace

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
	"github.com/xyctruth/profiler/pkg/storage"
	"github.com/xyctruth/profiler/pkg/storage/badger"
)

func initTraceData(s storage.Store, t *testing.T) (uint64, uint64, uint64) {
	invalidId, err := s.SaveProfile([]byte{}, time.Second*10)
	require.Equal(t, nil, err)
	require.Equal(t, uint64(0), invalidId)

	invalidId2, err := s.SaveProfile([]byte("haha"), time.Second*10)
	require.Equal(t, nil, err)
	require.Equal(t, uint64(1), invalidId2)

	traceBytes, err := ioutil.ReadFile("../testdata/trace.gz")
	require.Equal(t, nil, err)
	id, err := s.SaveProfile(traceBytes, time.Second*10)
	require.Equal(t, nil, err)
	require.Equal(t, uint64(2), id)
	return invalidId, invalidId2, id
}

func TestServer(t *testing.T) {
	dir, err := ioutil.TempDir("./", "temp-*")
	require.Equal(t, nil, err)
	defer os.RemoveAll(dir)

	store := badger.NewStore(badger.DefaultOptions(dir))

	pprofServer := NewServer("/api/trace/ui", store, 1*time.Second)
	defer pprofServer.Exit()

	httpServer := httptest.NewServer(pprofServer.mux)
	defer httpServer.Close()

	// create httpexpect instance
	e := httpexpect.New(t, httpServer.URL)

	e.GET("/badUrl").
		Expect().
		Status(http.StatusBadRequest).Text().Equal("Invalid parameter\n")

	testUI(e, store, t, pprofServer)
}

func testUI(e *httpexpect.Expect, store storage.Store, t *testing.T, server *Server) {
	invalidId, invalidId2, id := initTraceData(store, t)

	e.GET("/api/trace/ui/1999").
		Expect().
		Status(http.StatusNotFound).Text().Equal("Profile not found\n")

	e.GET(fmt.Sprintf("/api/trace/ui/%d", invalidId)).
		Expect().
		Status(http.StatusInternalServerError).Text().Equal("EOF\n")

	e.GET(fmt.Sprintf("/api/trace/ui/%d", invalidId2)).
		Expect().
		Status(http.StatusInternalServerError).Text().Equal("unexpected EOF\n")

	e.GET(fmt.Sprintf("/api/trace/ui/%d", id)).
		Expect().
		Status(http.StatusOK).Header("Content-Type").Equal("text/html; charset=utf-8")

	server.gc()

	e.GET(fmt.Sprintf("/api/trace/ui/%d", id)).
		Expect().
		Status(http.StatusOK).Header("Content-Type").Equal("text/html; charset=utf-8")

}
