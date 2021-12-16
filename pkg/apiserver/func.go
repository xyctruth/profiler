package apiserver

import "regexp"

func extractProfileId(path string) string {
	reg, _ := regexp.Compile(`([\d]+)`)
	return reg.FindString(path)
}
