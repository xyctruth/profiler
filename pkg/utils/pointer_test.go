package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBool(t *testing.T) {
	boolPtr := Bool(true)
	b := true
	require.NotEqual(t, b, boolPtr)
	require.Equal(t, &b, boolPtr)
}

func TestBoolPtr(t *testing.T) {
	boolPtr := BoolPtr(false)
	b := false
	require.NotEqual(t, b, boolPtr)
	require.Equal(t, &b, boolPtr)
}
