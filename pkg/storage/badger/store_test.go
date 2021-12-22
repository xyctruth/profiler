package badger

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/xyctruth/profiler/pkg/storage"

	"github.com/stretchr/testify/require"
)

var (
	profileMeta = &storage.ProfileMeta{
		ProfileID:      1,
		Timestamp:      time.Now().UnixNano() / time.Millisecond.Nanoseconds(),
		Duration:       time.Now().UnixNano(),
		SampleTypeUnit: "count",
		SampleType:     "alloc_objects",
		ProfileType:    "alloc",
		TargetName:     "profiler-server",
		Value:          1,
	}
)

func TestNewStore(t *testing.T) {
	dir, err := ioutil.TempDir("./", "data-*")
	defer os.RemoveAll(dir)
	require.Equal(t, nil, err)
	s := NewStore(dir)
	require.NotEqual(t, nil, s)
}

func TestSaveProfileMeta(t *testing.T) {
	dir, err := ioutil.TempDir("./", "data-*")
	defer os.RemoveAll(dir)
	require.Equal(t, nil, err)
	s := NewStore(dir)
	require.NotEqual(t, nil, s)

	err = s.SaveProfileMeta([]*storage.ProfileMeta{profileMeta}, time.Second*1)
	require.Equal(t, nil, err)

	// targets
	targets, err := s.ListTarget()
	require.Equal(t, nil, err)
	require.Equal(t, 1, len(targets))

	time.Sleep(1 * time.Second)
	// targets ttl
	targets, err = s.ListTarget()
	require.Equal(t, nil, err)
	require.Equal(t, 0, len(targets))
}
