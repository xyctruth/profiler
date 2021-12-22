package apiserver

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
	"github.com/xyctruth/profiler/pkg/storage/badger"
)

func TestPprofServer(t *testing.T) {
	dir, err := ioutil.TempDir("./", "temp-*")
	require.Equal(t, nil, err)
	defer os.RemoveAll(dir)

	store := badger.NewStore(dir)

	invalidId, err := store.SaveProfile([]byte{}, time.Second*10)
	require.Equal(t, nil, err)
	require.Equal(t, uint64(0), invalidId)

	invalidId2, err := store.SaveProfile([]byte("haha"), time.Second*10)
	require.Equal(t, nil, err)
	require.Equal(t, uint64(1), invalidId2)

	profileBytes, err := ioutil.ReadFile("./profile.pb.gz")
	require.Equal(t, nil, err)
	id, err := store.SaveProfile(profileBytes, time.Second*10)
	require.Equal(t, nil, err)
	require.Equal(t, uint64(2), id)

	pprofServer := newPprofServer("/api/pprof/ui", store)

	httpServer := httptest.NewServer(pprofServer.mux)
	defer httpServer.Close()

	// create httpexpect instance
	e := httpexpect.New(t, httpServer.URL)

	e.GET("/badUrl").
		Expect().
		Status(http.StatusBadRequest).Text().Equal("Invalid parameter\n")

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
