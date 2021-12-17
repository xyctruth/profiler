package apiserver

import "testing"

func TestExtractProfileId(t *testing.T) {
	id := extractProfileID("/api/pprof/ui/10009/")
	if id != "10009" {
		t.Error("error id is", id)
	}
}

func TestRemovePrefixSampleType(t *testing.T) {
	rawQuery := removePrefixSampleType("si=heap_alloc_space")
	if rawQuery != "si=alloc_space" {
		t.Error("error rawQuery is", rawQuery)
	}
}
