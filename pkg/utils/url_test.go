package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractProfileId(t *testing.T) {

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "/10009/",
			input:   "/api/pprof/ui/10009/",
			want:    "10009",
			wantErr: false,
		},
		{
			name:    "/10009/top",
			input:   "/api/pprof/ui/10009/top",
			want:    "10009",
			wantErr: false,
		},
		{
			name:    "/10009",
			input:   "/api/pprof/ui/10009",
			want:    "10009",
			wantErr: false,
		},
		{
			name:    "/10009asd",
			input:   "/api/pprof/ui/10009asd",
			want:    "",
			wantErr: false,
		},
		{
			name:    "/10009asd/",
			input:   "/api/pprof/ui/10009asd/",
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractProfileID(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRemovePrefixSampleType(t *testing.T) {
	rawQuery := RemovePrefixSampleType("si=heap_alloc_space")
	assert.Equal(t, "si=alloc_space", rawQuery)

	rawQuery = RemovePrefixSampleType("")
	assert.Equal(t, "", rawQuery)
}
