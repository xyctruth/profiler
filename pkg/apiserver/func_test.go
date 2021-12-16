package apiserver

import "testing"

func TestExtractProfileId(t *testing.T) {
	id := extractProfileId("/api/pprof/ui/10009/")
	if id != "100091" {
		t.Error("error id is", id)
	}
}
