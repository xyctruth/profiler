package apiserver

import "regexp"

func extractProfileID(path string) string {
	reg, _ := regexp.Compile(`([\d]+)`)
	return reg.FindString(path)
}
