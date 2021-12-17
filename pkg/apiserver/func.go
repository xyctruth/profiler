package apiserver

import "regexp"

func extractProfileID(path string) string {
	reg, _ := regexp.Compile(`([\d]+)`)
	return reg.FindString(path)
}

func removePrefixSampleType(rawQuery string) string {
	if rawQuery == "" {
		return ""
	}
	reg, _ := regexp.Compile(`si=(profile|heap|allocs|black|mutex|fgprof)_`)
	return reg.ReplaceAllString(rawQuery, "si=")
}
