package apiserver

import (
	"regexp"
	"strings"
)

func extractProfileID(path string) string {
	reg, _ := regexp.Compile(`/([\d]+)(/|$)`)
	s := reg.FindString(path)
	return strings.ReplaceAll(s, "/", "")
}

func removePrefixSampleType(rawQuery string) string {
	if rawQuery == "" {
		return ""
	}
	reg, _ := regexp.Compile(`si=(profile|heap|allocs|black|mutex|fgprof)_`)
	return reg.ReplaceAllString(rawQuery, "si=")
}
