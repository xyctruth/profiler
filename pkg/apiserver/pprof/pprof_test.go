package pprof

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/xyctruth/profiler/pkg/storage"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
	"github.com/xyctruth/profiler/pkg/storage/badger"
)

func initProfileData(s storage.Store, t *testing.T) (uint64, uint64, uint64) {
	invalidId, err := s.SaveProfile([]byte{}, time.Second*10)
	require.Equal(t, nil, err)
	require.Equal(t, uint64(0), invalidId)

	invalidId2, err := s.SaveProfile([]byte("haha"), time.Second*10)
	require.Equal(t, nil, err)
	require.Equal(t, uint64(1), invalidId2)

	profileBytes, err := ioutil.ReadFile("../profile.gz")
	require.Equal(t, nil, err)
	id, err := s.SaveProfile(profileBytes, time.Second*10)
	require.Equal(t, nil, err)
	require.Equal(t, uint64(2), id)
	return invalidId, invalidId2, id
}

func TestPprofServer(t *testing.T) {
	dir, err := ioutil.TempDir("./", "temp-*")
	require.Equal(t, nil, err)
	defer os.RemoveAll(dir)

	store := badger.NewStore(dir)

	pprofServer := NewPProfServer("/api/pprof/ui", store)

	httpServer := httptest.NewServer(pprofServer.mux)
	defer httpServer.Close()

	// create httpexpect instance
	e := httpexpect.New(t, httpServer.URL)

	e.GET("/badUrl").
		Expect().
		Status(http.StatusBadRequest).Text().Equal("Invalid parameter\n")

	testPprofUI(e, store, t)
}

func testPprofUI(e *httpexpect.Expect, store storage.Store, t *testing.T) {
	invalidId, invalidId2, id := initProfileData(store, t)

	e.GET("/api/pprof/ui/1999").
		Expect().
		Status(http.StatusNotFound).Text().Equal("Profile not found\n")

	e.GET(fmt.Sprintf("/api/pprof/ui/%d", invalidId)).
		Expect().
		Status(http.StatusInternalServerError).Text().Equal("failed to fetch any source profiles\n")

	e.GET(fmt.Sprintf("/api/pprof/ui/%d", invalidId2)).
		Expect().
		Status(http.StatusInternalServerError).Text().Equal("failed to fetch any source profiles\n")

	e.GET(fmt.Sprintf("/api/pprof/ui/%d", id)).
		Expect().
		Status(http.StatusOK).Header("Content-Type").Equal("text/html")

	e.GET(fmt.Sprintf("/api/pprof/ui/%d", id)).WithQuery("si", "alloc_objects").
		Expect().
		Status(http.StatusOK).Header("Content-Type").Equal("text/html")

	e.GET(fmt.Sprintf("/api/pprof/ui/%d/", id)).WithQuery("si", "alloc_space").
		Expect().
		Status(http.StatusOK).Header("Content-Type").Equal("text/html")

	e.GET(fmt.Sprintf("/api/pprof/ui/%d/", id)).WithQuery("si", "ppp").
		Expect().
		Status(http.StatusBadRequest).Text().Equal("sample_index \"ppp\" must be one of: [alloc_objects alloc_space inuse_objects inuse_space]\n")

	e.GET(fmt.Sprintf("/api/pprof/ui/%d/top", id)).WithQuery("si", "alloc_space").
		Expect().
		Status(http.StatusOK).Header("Content-Type").Equal("text/html")

	e.GET(fmt.Sprintf("/api/pprof/ui/%d/top", id)).WithQuery("si", "alloc_space").
		Expect().
		Status(http.StatusOK).Header("Content-Type").Equal("text/html")

	e.GET(fmt.Sprintf("/api/pprof/ui/%d/toperror", id)).WithQuery("si", "alloc_space").
		Expect().
		Status(http.StatusOK).Header("Content-Type").Equal("text/html")

}
