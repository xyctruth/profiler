package utils

import (
	"regexp"
)

func RemoveSampleTypePrefix(sampleType string) string {
	if sampleType == "" {
		return ""
	}
	reg, _ := regexp.Compile(`(profile|heap|allocs|black|mutex)_`)
	return reg.ReplaceAllString(sampleType, "")
}
