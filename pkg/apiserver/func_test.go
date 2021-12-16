package apiserver

import "testing"

func TestExtractProfileId(t *testing.T) {
	id := extractProfileID("/api/pprof/ui/10009/")
	if id != "10009" {
		t.Error("error id is", id)
	}
}
