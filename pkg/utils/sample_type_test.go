package utils

import (
	"testing"
)

func TestRemoveSampleTypePrefix(t *testing.T) {
	a1 := RemoveSampleTypePrefix("profile_haha")
	if a1 != "haha" {
		t.Error(" error, a1 is", a1)
	}

	a1 = RemoveSampleTypePrefix("heap_haha")
	if a1 != "haha" {
		t.Error(" error, a1 is", a1)
	}

	a1 = RemoveSampleTypePrefix("allocs_haha")
	if a1 != "haha" {
		t.Error(" error, a1 is", a1)
	}

	a1 = RemoveSampleTypePrefix("mutex_haha")
	if a1 != "haha" {
		t.Error(" error, a1 is", a1)
	}

	a1 = RemoveSampleTypePrefix("mutex1_haha")
	if a1 != "mutex1_haha" {
		t.Error(" error, a1 is", a1)
	}
}
