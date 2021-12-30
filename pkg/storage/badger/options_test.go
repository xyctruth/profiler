package badger

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOptions(t *testing.T) {
	opt := DefaultOptions("./test")
	require.Equal(t, "./test", opt.Path)
	require.Equal(t, 5*time.Minute, opt.GCInternal)

	opt = opt.WithGCInternal(3 * time.Minute)
	require.Equal(t, 3*time.Minute, opt.GCInternal)
}
