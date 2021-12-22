package badger

import (
	"io/ioutil"
	"os"
	"strconv"
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
		ProfileType:    "heap",
		TargetName:     "profiler-server",
		Value:          1,
	}
	profileMetas = []*storage.ProfileMeta{
		{
			ProfileID:      1,
			Timestamp:      time.Now().UnixNano() / time.Millisecond.Nanoseconds(),
			Duration:       time.Now().UnixNano(),
			SampleTypeUnit: "count",
			SampleType:     "alloc_objects",
			ProfileType:    "heap",
			TargetName:     "profiler-server",
			Value:          100,
		},
		{
			ProfileID:      2,
			Timestamp:      time.Now().UnixNano() / time.Millisecond.Nanoseconds(),
			Duration:       time.Now().UnixNano(),
			SampleTypeUnit: "bytes",
			SampleType:     "alloc_space",
			ProfileType:    "heap",
			TargetName:     "profiler-server",
			Value:          200,
		},
		{
			ProfileID:      3,
			Timestamp:      time.Now().UnixNano() / time.Millisecond.Nanoseconds(),
			Duration:       time.Now().UnixNano(),
			SampleTypeUnit: "count",
			SampleType:     "inuse_objects",
			ProfileType:    "heap",
			TargetName:     "server2",
			Value:          300,
		},
		{
			ProfileID:      4,
			Timestamp:      time.Now().UnixNano() / time.Millisecond.Nanoseconds(),
			Duration:       time.Now().UnixNano(),
			SampleTypeUnit: "bytes",
			SampleType:     "inuse_space",
			ProfileType:    "heap",
			TargetName:     "server2",
			Value:          400,
		},
	}
)

func TestNewStore(t *testing.T) {
	dir, err := ioutil.TempDir("./", "data-*")
	defer os.RemoveAll(dir)
	require.Equal(t, nil, err)
	s := NewStore(dir)
	defer s.Release()
	require.NotEqual(t, nil, s)
}

func TestProfile(t *testing.T) {
	dir, err := ioutil.TempDir("./", "data-*")
	defer os.RemoveAll(dir)
	require.Equal(t, nil, err)
	s := NewStore(dir)
	defer s.Release()
	require.NotEqual(t, nil, s)

	id, err := s.SaveProfile([]byte{}, 1*time.Second)
	require.Equal(t, nil, err)
	require.NotEqual(t, 0, id)

	_, err = s.GetProfile(strconv.FormatUint(id, 10))
	require.Equal(t, nil, err)

	// Waiting for the overdue
	time.Sleep(1 * time.Second)
	_, err = s.GetProfile(strconv.FormatUint(id, 10))
	require.NotEqual(t, nil, err)

}

func TestProfileMeta(t *testing.T) {
	dir, err := ioutil.TempDir("./", "data-*")
	defer os.RemoveAll(dir)
	require.Equal(t, nil, err)
	s := NewStore(dir)
	defer s.Release()
	require.NotEqual(t, nil, s)

	err = s.SaveProfileMeta([]*storage.ProfileMeta{profileMeta}, time.Second*1)
	require.Equal(t, nil, err)

	min := time.Now().Add(-1 * time.Hour)
	max := time.Now()

	targets, err := s.ListTarget()
	require.Equal(t, nil, err)
	require.Equal(t, 1, len(targets))

	groupTargets, err := s.ListGroupSampleType()
	require.Equal(t, nil, err)
	require.Equal(t, 1, len(groupTargets))
	for _, tg := range groupTargets {
		require.Equal(t, 1, len(tg))
	}

	sampleTypes, err := s.ListSampleType()
	require.Equal(t, nil, err)
	require.Equal(t, 1, len(sampleTypes))

	profileMetas, err := s.ListProfileMeta(sampleTypes[0], targets, min, max)
	require.Equal(t, nil, err)
	require.Equal(t, 1, len(profileMetas))

	// Waiting for the overdue
	time.Sleep(1 * time.Second)

	{

		ttlTargets, err := s.ListTarget()
		require.Equal(t, nil, err)
		require.Equal(t, 0, len(ttlTargets))

		ttlGroupTargets, err := s.ListGroupSampleType()
		require.Equal(t, nil, err)
		require.Equal(t, 0, len(ttlGroupTargets))

		ttlSampleTypes, err := s.ListSampleType()
		require.Equal(t, nil, err)
		require.Equal(t, 0, len(ttlSampleTypes))

		ttlProfileMetas, err := s.ListProfileMeta(sampleTypes[0], targets, min, max)
		require.Equal(t, nil, err)
		require.Equal(t, 0, len(ttlProfileMetas))
	}
}

func TestProfileMetaArray(t *testing.T) {
	dir, err := ioutil.TempDir("./", "data-*")
	defer os.RemoveAll(dir)
	require.Equal(t, nil, err)
	s := NewStore(dir)
	defer s.Release()
	require.NotEqual(t, nil, s)

	err = s.SaveProfileMeta(profileMetas, time.Second*1)
	require.Equal(t, nil, err)

	min := time.Now().Add(-1 * time.Hour)
	max := time.Now()

	targets, err := s.ListTarget()
	require.Equal(t, nil, err)
	require.Equal(t, 2, len(targets))

	groupTargets, err := s.ListGroupSampleType()
	require.Equal(t, nil, err)
	require.Equal(t, 1, len(groupTargets))
	for _, tg := range groupTargets {
		require.Equal(t, 4, len(tg))
	}

	sampleTypes, err := s.ListSampleType()
	require.Equal(t, nil, err)
	require.Equal(t, 4, len(sampleTypes))

	for _, sampleType := range sampleTypes {
		profileMetas, err := s.ListProfileMeta(sampleType, targets, min, max)
		require.Equal(t, nil, err)
		require.Equal(t, 1, len(profileMetas))
	}

	// Waiting for the overdue
	time.Sleep(1 * time.Second)

	{

		ttlTargets, err := s.ListTarget()
		require.Equal(t, nil, err)
		require.Equal(t, 0, len(ttlTargets))

		ttlGroupTargets, err := s.ListGroupSampleType()
		require.Equal(t, nil, err)
		require.Equal(t, 0, len(ttlGroupTargets))

		ttlSampleTypes, err := s.ListSampleType()
		require.Equal(t, nil, err)
		require.Equal(t, 0, len(ttlSampleTypes))

		ttlProfileMetas, err := s.ListProfileMeta(sampleTypes[0], targets, min, max)
		require.Equal(t, nil, err)
		require.Equal(t, 0, len(ttlProfileMetas))
	}
}
