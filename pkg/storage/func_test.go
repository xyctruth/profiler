package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCompareKey(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name string
		key1 []byte
		key2 []byte
		want bool
	}{
		{
			name: "less",
			key1: BuildKey(now),
			key2: BuildKey(now.Add(1 * time.Second)),
			want: true,
		},
		{
			name: "equal",
			key1: BuildKey(now),
			key2: BuildKey(now),
			want: true,
		},
		{
			name: "equal millisecond diff",
			key1: BuildKey(now),
			key2: BuildKey(now.Add(1 * time.Millisecond)),
			want: true,
		},
		{
			name: "greater",
			key1: BuildKey(now),
			key2: BuildKey(now.Add(-1 * time.Second)),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CompareKey(tt.key1, tt.key2)
			assert.Equal(t, tt.want, got)
		})
	}

}
