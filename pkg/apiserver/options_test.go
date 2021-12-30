package apiserver

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOptions(t *testing.T) {
	opt := DefaultOptions(nil)
	require.Equal(t, nil, opt.Store)
	require.Equal(t, ":8080", opt.Addr)
	require.Equal(t, 2*time.Minute, opt.GCInternal)

	opt = opt.WithGCInternal(3 * time.Minute)
	require.Equal(t, 3*time.Minute, opt.GCInternal)

	opt = opt.WithAddr(":8081")
	require.Equal(t, 3*time.Minute, opt.GCInternal)
}
