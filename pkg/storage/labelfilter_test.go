package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLabelFilter(t *testing.T) {
	filters := []LabelFilter{
		{
			Label:     Label{},
			condition: "AND",
		}, {
			Label:     Label{},
			condition: "AND",
		}, {
			Label:     Label{},
			condition: "OR",
		}}

	s1 := []string{"a1", "a2"}

	s2 := []string{"a1", "a2", "a3"}

	s3 := []string{"a4"}

	s4 := filters[1].Policy(filters[0].Policy(s1, s2), s3)
	require.Equal(t, []string{}, s4)

	s5 := filters[2].Policy(filters[0].Policy(s1, s2), s3)
	require.Equal(t, []string{"a1", "a2", "a4"}, s5)
}

func TestLabelFilterStrategy(t *testing.T) {
	s1 := []string{"a1", "a2"}

	s2 := []string{"a1", "a2", "a3"}

	s3 := []string{"a4"}

	s4 := Intersect(s1, s2)
	require.Equal(t, []string{"a1", "a2"}, s4)

	s4 = Intersect(s4, s3)
	require.Equal(t, []string{}, s4)

	s5 := Union(Intersect(s1, s2), s3)
	require.Equal(t, []string{"a1", "a2", "a4"}, s5)
}
