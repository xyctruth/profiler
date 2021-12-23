package apiserver

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/xyctruth/profiler/pkg/storage"

	"github.com/stretchr/testify/require"

	"github.com/xyctruth/profiler/pkg/storage/badger"

	"github.com/gavv/httpexpect/v2"
)

var (
	profileMetas = []*storage.ProfileMeta{
		{
			ProfileID:      1,
			Timestamp:      time.Now().UnixNano() / time.Millisecond.Nanoseconds(),
			Duration:       time.Now().UnixNano(),
			SampleTypeUnit: "count",
			SampleType:     "heap_alloc_objects",
			ProfileType:    "heap",
			TargetName:     "profiler-server",
			Value:          100,
		},
		{
			ProfileID:      2,
			Timestamp:      time.Now().UnixNano() / time.Millisecond.Nanoseconds(),
			Duration:       time.Now().UnixNano(),
			SampleTypeUnit: "bytes",
			SampleType:     "heap_alloc_space",
			ProfileType:    "heap",
			TargetName:     "profiler-server",
			Value:          200,
		},
		{
			ProfileID:      3,
			Timestamp:      time.Now().UnixNano() / time.Millisecond.Nanoseconds(),
			Duration:       time.Now().UnixNano(),
			SampleTypeUnit: "count",
			SampleType:     "heap_inuse_objects",
			ProfileType:    "heap",
			TargetName:     "server2",
			Value:          300,
		},
		{
			ProfileID:      4,
			Timestamp:      time.Now().UnixNano() / time.Millisecond.Nanoseconds(),
			Duration:       time.Now().UnixNano(),
			SampleTypeUnit: "bytes",
			SampleType:     "heap_inuse_space",
			ProfileType:    "heap",
			TargetName:     "server2",
			Value:          400,
		},
		{
			ProfileID:      5,
			Timestamp:      time.Now().UnixNano() / time.Millisecond.Nanoseconds(),
			Duration:       time.Now().UnixNano(),
			SampleTypeUnit: "bytes",
			SampleType:     "heap_inuse_space",
			ProfileType:    "heap",
			TargetName:     "server3",
			Value:          400,
		},
	}
)

// Rid of debug output
func init() {
	gin.SetMode(gin.TestMode)
}

func initMateData(s storage.Store, t *testing.T) {
	err := s.SaveProfileMeta(profileMetas, time.Second*3)
	require.Equal(t, nil, err)
}

func initProfileData(s storage.Store, t *testing.T) (uint64, uint64, uint64) {
	invalidId, err := s.SaveProfile([]byte{}, time.Second*10)
	require.Equal(t, nil, err)
	require.Equal(t, uint64(0), invalidId)

	invalidId2, err := s.SaveProfile([]byte("haha"), time.Second*10)
	require.Equal(t, nil, err)
	require.Equal(t, uint64(1), invalidId2)

	profileBytes, err := ioutil.ReadFile("./profile.pb.gz")
	require.Equal(t, nil, err)
	id, err := s.SaveProfile(profileBytes, time.Second*10)
	require.Equal(t, nil, err)
	require.Equal(t, uint64(2), id)
	return invalidId, invalidId2, id
}

func TestApiServer(t *testing.T) {
	dir, err := ioutil.TempDir("./", "temp-*")
	require.Equal(t, nil, err)
	defer os.RemoveAll(dir)
	s := badger.NewStore(dir)
	apiServer := NewAPIServer(":8080", s)
	apiServer.Run()
	defer apiServer.Stop()
}

func TestBasisAPI(t *testing.T) {
	dir, err := ioutil.TempDir("./", "temp-*")
	require.Equal(t, nil, err)
	defer os.RemoveAll(dir)
	s := badger.NewStore(dir)

	initMateData(s, t)

	apiServer := NewAPIServer(":8080", s)

	e := getExpect(apiServer, t)

	e.GET("/").
		Expect().
		Status(http.StatusNotFound)

	e.GET("/api/healthz").
		Expect().
		Status(http.StatusOK)

	e.GET("/api/targets").
		Expect().
		Status(http.StatusOK).JSON().Array().Contains("profiler-server", "server2", "server3")

	e.GET("/api/sample_types").
		Expect().
		Status(http.StatusOK).JSON().Array().Contains("heap_alloc_objects", "heap_alloc_space", "heap_inuse_space", "heap_inuse_space")

	res := e.GET("/api/group_sample_types").
		Expect().
		Status(http.StatusOK).JSON().Object()
	res.Path("$").Object().Keys().Length().Equal(1)
	res.Path("$.heap").Array().Contains("heap_alloc_objects", "heap_alloc_space", "heap_inuse_space", "heap_inuse_space")

	time.Sleep(3 * time.Second)

	e.GET("/api/targets").
		Expect().
		Status(http.StatusOK).JSON().Array().Length().Equal(0)

	e.GET("/api/sample_types").
		Expect().
		Status(http.StatusOK).JSON().Array().Length().Equal(0)

	e.GET("/api/group_sample_types").
		Expect().
		Status(http.StatusOK).JSON().Object().Path("$").Object().Keys().Length().Equal(0)
}

func TestListProfileMeta(t *testing.T) {
	dir, err := ioutil.TempDir("./", "temp-*")
	require.Equal(t, nil, err)
	defer os.RemoveAll(dir)

	s := badger.NewStore(dir)

	initMateData(s, t)

	apiServer := NewAPIServer(":8080", s)

	e := getExpect(apiServer, t)

	e.GET("/api/profile_meta").
		Expect().
		Status(http.StatusNotFound)

	e.GET("/api/profile_meta/heap_inuse_space").
		Expect().
		Status(http.StatusBadRequest).Text().Equal("start_time or end_time is empty")

	startTimestamp := time.Now().Add(-1*time.Minute).UnixNano() / time.Millisecond.Nanoseconds()
	endTimestamp := time.Now().UnixNano() / time.Millisecond.Nanoseconds()

	e.GET("/api/profile_meta/heap_inuse_space").WithQuery("start_time", startTimestamp).
		Expect().
		Status(http.StatusBadRequest).Text().Equal("start_time or end_time is empty")

	e.GET("/api/profile_meta/heap_inuse_space").
		WithQuery("start_time", startTimestamp).WithQuery("end_time", endTimestamp).
		Expect().
		Status(http.StatusBadRequest).Text().Contains("The time format must be RFC3339")

	startTime := time.Now().Add(-1 * time.Minute).Format(time.RFC3339)
	endTime := time.Now().Format(time.RFC3339)

	e.GET("/api/profile_meta/heap_inuse_space").
		WithQuery("start_time", startTime).WithQuery("end_time", endTimestamp).
		Expect().
		Status(http.StatusBadRequest).Text().Contains("The time format must be RFC3339")

	e.GET("/api/profile_meta/heap_inuse_space").
		WithQuery("start_time", startTime).WithQuery("end_time", endTime).
		Expect().
		Status(http.StatusOK).JSON().Array().Length().Equal(2)

	e.GET("/api/profile_meta/heap_inuse_space").
		WithQuery("start_time", startTime).WithQuery("end_time", endTime).WithQuery("targets", "server2").
		Expect().
		Status(http.StatusOK).JSON().Array().Length().Equal(1)

	e.GET("/api/profile_meta/heap_inuse_space").
		WithQuery("start_time", startTime).WithQuery("end_time", endTime).WithQuery("targets", "notfound").
		Expect().
		Status(http.StatusOK).JSON().Array().Length().Equal(0)
}

func TestGetProfile(t *testing.T) {
	dir, err := ioutil.TempDir("./", "temp-*")
	require.Equal(t, nil, err)
	defer os.RemoveAll(dir)

	s := badger.NewStore(dir)
	invalidId, invalidId2, id := initProfileData(s, t)
	apiServer := NewAPIServer(":8080", s)
	e := getExpect(apiServer, t)

	e.GET("/api/profile/999").
		Expect().
		Status(http.StatusNotFound)

	e.GET(fmt.Sprintf("/api/profile/%d", invalidId)).
		Expect().
		Status(http.StatusInternalServerError).Text().Equal("parsing profile: empty input file")

	e.GET(fmt.Sprintf("/api/profile/%d", invalidId2)).
		Expect().
		Status(http.StatusInternalServerError).Text().Equal("parsing profile: unrecognized profile format")

	e.GET(fmt.Sprintf("/api/profile/%d", id)).
		Expect().
		Status(http.StatusOK).Header("Content-Type").Equal("application/vnd.google.protobuf+gzip")
}

func TestWebProfile(t *testing.T) {
	dir, err := ioutil.TempDir("./", "temp-*")
	require.Equal(t, nil, err)
	defer os.RemoveAll(dir)

	s := badger.NewStore(dir)
	apiServer := NewAPIServer(":8080", s)
	e := getExpect(apiServer, t)
	testPprofUI(e, s, t)
}

func getExpect(apiServer *APIServer, t *testing.T) *httpexpect.Expect {
	handler := apiServer.router

	return httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(handler),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
	})
}