package utils

import (
	"regexp"
	"strings"
)

func ExtractProfileID(path string) string {
	reg, _ := regexp.Compile(`/([\d]+)(/|$)`)
	return strings.ReplaceAll(reg.FindString(path), "/", "")
}

func RemovePrefixSampleType(rawQuery string) string {
	reg, _ := regexp.Compile(`si=(profile|heap|allocs|black|mutex|fgprof)_`)
	return reg.ReplaceAllString(rawQuery, "si=")
}
