package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemoveDuplicateElement(t *testing.T) {
	s := []string{"a", "b", "c", "d", "b", "c", "c"}
	s = RemoveDuplicateElement(s)
	require.Equal(t, 4, len(s))
	require.Equal(t, []string{"a", "b", "c", "d"}, s)
}
