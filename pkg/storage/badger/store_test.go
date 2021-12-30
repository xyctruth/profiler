package badger

import (
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/require"
	"github.com/xyctruth/profiler/pkg/storage"
)

var (
	profileMeta = &storage.ProfileMeta{
		ProfileID:      1,
		Timestamp:      time.Now().UnixNano() / time.Millisecond.Nanoseconds(),
		Duration:       time.Now().UnixNano(),
		SampleTypeUnit: "count",
		SampleType:     "heap_alloc_objects",
		ProfileType:    "heap",
		TargetName:     "profiler-server",
		Value:          1,
	}

	traceMeta = &storage.ProfileMeta{
		ProfileID:   1,
		Timestamp:   time.Now().UnixNano() / time.Millisecond.Nanoseconds(),
		Duration:    time.Now().UnixNano(),
		SampleType:  "trace",
		ProfileType: "trace",
		TargetName:  "profiler-server",
	}
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

func TestNewStore(t *testing.T) {
	dir, err := ioutil.TempDir("./", "temp-*")
	defer os.RemoveAll(dir)
	require.Equal(t, nil, err)
	s := NewStore(DefaultOptions(dir))
	defer s.Release()
	require.NotEqual(t, nil, s)
}

func TestGC(t *testing.T) {
	dir, err := ioutil.TempDir("./", "temp-*")
	defer os.RemoveAll(dir)
	require.Equal(t, nil, err)
	s := NewStore(DefaultOptions(dir))
	defer s.Release()
	require.NotEqual(t, nil, s)
	s.(*store).gc()
}

func TestProfile(t *testing.T) {
	dir, err := ioutil.TempDir("./", "temp-*")
	defer os.RemoveAll(dir)
	require.Equal(t, nil, err)
	s := NewStore(DefaultOptions(dir))
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
	dir, err := ioutil.TempDir("./", "temp-*")
	defer os.RemoveAll(dir)
	require.Equal(t, nil, err)
	s := NewStore(DefaultOptions(dir))
	defer s.Release()
	require.NotEqual(t, nil, s)

	err = s.SaveProfileMeta([]*storage.ProfileMeta{profileMeta}, time.Second*2)
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
	time.Sleep(2 * time.Second)

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
	dir, err := ioutil.TempDir("./", "temp-*")
	defer os.RemoveAll(dir)
	require.Equal(t, nil, err)
	s := NewStore(DefaultOptions(dir).WithGCInternal(time.Second))
	defer s.Release()
	require.NotEqual(t, nil, s)

	err = s.SaveProfileMeta(profileMetas, time.Second*3)
	require.Equal(t, nil, err)

	min := time.Now().Add(-1 * time.Hour)
	max := time.Now()

	targets, err := s.ListTarget()
	require.Equal(t, nil, err)
	require.Equal(t, 3, len(targets))

	groupTargets, err := s.ListGroupSampleType()
	require.Equal(t, nil, err)
	require.Equal(t, 1, len(groupTargets))
	for _, tg := range groupTargets {
		require.Equal(t, 4, len(tg))
	}

	sampleTypes, err := s.ListSampleType()
	require.Equal(t, nil, err)
	require.Equal(t, 4, len(sampleTypes))

	{
		profileMetas, err := s.ListProfileMeta("heap_inuse_space", targets, min, max)
		require.Equal(t, nil, err)
		require.Equal(t, 2, len(profileMetas))

		profileMetas, err = s.ListProfileMeta("heap_inuse_space", nil, min, max)
		require.Equal(t, nil, err)
		require.Equal(t, 2, len(profileMetas))

		profileMetas, err = s.ListProfileMeta("heap_inuse_space", []string{"server2"}, min, max)
		require.Equal(t, nil, err)
		require.Equal(t, 1, len(profileMetas))

		profileMetas, err = s.ListProfileMeta("heap_inuse_objects", nil, min, max)
		require.Equal(t, nil, err)
		require.Equal(t, 1, len(profileMetas))

		profileMetas, err = s.ListProfileMeta("heap_inuse_objects1", nil, min, max)
		require.Equal(t, nil, err)
		require.Equal(t, 0, len(profileMetas))
	}

	// Waiting for the overdue
	time.Sleep(3 * time.Second)

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
	s.Release()
}

func BenchmarkBadger1(b *testing.B) {
	dir, err := ioutil.TempDir("./", "temp-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	db, err := badger.Open(
		badger.DefaultOptions(dir).
			WithLoggingLevel(3).
			WithBypassLockGuard(true))

	if err != nil {
		panic(err)
	}

	s := &store{
		db:  db,
		opt: DefaultOptions(dir),
	}

	s.seq, err = s.db.GetSequence(Sequence, 1000)
	if err != nil {
		panic(err)
	}

	defer s.Release()
	res, err := os.ReadFile("./trace_119091.gz")
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		_, err = s.SaveProfile(res, time.Hour*24*7)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkBadger2(b *testing.B) {
	dir, err := ioutil.TempDir("./", "temp-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	db, err := badger.Open(
		badger.DefaultOptions(dir).
			WithLoggingLevel(3).
			WithBypassLockGuard(true).
			//WithNumMemtables(1).
			//WithNumLevelZeroTables(1).
			//WithNumLevelZeroTablesStall(2).
			WithValueLogFileSize(64 << 20))

	if err != nil {
		panic(err)
	}

	s := &store{
		db:  db,
		opt: DefaultOptions(dir),
	}

	s.seq, err = s.db.GetSequence(Sequence, 1000)
	if err != nil {
		panic(err)
	}

	defer s.Release()
	res, err := os.ReadFile("./trace_119091.gz")
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		_, err = s.SaveProfile(res, time.Hour*24*7)
		if err != nil {
			panic(err)
		}
	}
}
