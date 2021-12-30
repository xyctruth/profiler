package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	profileMeta := &ProfileMeta{
		ProfileID:      1,
		Timestamp:      time.Now().UnixNano() / time.Millisecond.Nanoseconds(),
		Duration:       time.Now().UnixNano(),
		SampleTypeUnit: "count",
		SampleType:     "alloc_objects",
		ProfileType:    "heap",
		TargetName:     "profiler-server",
		Value:          1,
	}
	b, err := profileMeta.Encode()
	require.Equal(t, err, nil)

	var meta ProfileMeta
	err = meta.Decode(b)
	require.Equal(t, err, nil)
	require.Equal(t, profileMeta.ProfileID, meta.ProfileID)
	require.Equal(t, profileMeta.Timestamp, meta.Timestamp)
	require.Equal(t, profileMeta.Duration, meta.Duration)
	require.Equal(t, profileMeta.SampleTypeUnit, meta.SampleTypeUnit)
	require.Equal(t, profileMeta.SampleType, meta.SampleType)

	profileMeta.ProfileType = string(make([]byte, 1024, 1024))
	b, err = profileMeta.Encode()
	require.NotEqual(t, err, nil)
}
