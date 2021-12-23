package apiserver

import (
	"regexp"
	"strings"
)

func extractProfileID(path string) string {
	reg, _ := regexp.Compile(`/([\d]+)(/|$)`)
	return strings.ReplaceAll(reg.FindString(path), "/", "")
}

func removePrefixSampleType(rawQuery string) string {
	reg, _ := regexp.Compile(`si=(profile|heap|allocs|black|mutex|fgprof)_`)
	return reg.ReplaceAllString(rawQuery, "si=")
}
