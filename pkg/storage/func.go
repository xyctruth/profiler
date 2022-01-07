package storage

import (
	"bytes"
	"time"
)

func CompareKey(k, max []byte) bool {
	return bytes.Compare(k, max) <= 0
}

func BuildTimeKey(datetime time.Time) []byte {
	return []byte(datetime.Local().Format(time.RFC3339))
}
