package storage

import (
	"bytes"
	"time"
)

func CompareKey(k, max []byte) bool {
	i := bytes.Compare(k, max)
	return i <= 0
}

func BuildKey(datetime time.Time) []byte {
	key := datetime.Local().Format(time.RFC3339)
	return []byte(key)
}
