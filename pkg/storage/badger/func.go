package badger

import (
	"bytes"
	"strings"
	"time"

	"github.com/xyctruth/profiler/pkg/storage"
)

var (
	PrefixProfiles    = []byte("Profile")
	PrefixProfileMeta = []byte("ProfileMeta")
	PrefixSampleType  = []byte("SampleType")
	PrefixTarget      = []byte("Target")
	Sequence          = []byte("Sequence")
)

func buildProfileKey(id string) []byte {
	buf := bytes.NewBuffer(PrefixProfiles)
	buf.Write([]byte(id))
	return buf.Bytes()
}

func buildBaseProfileMetaKey(sampleType string, target string) []byte {
	buf := bytes.NewBuffer(PrefixProfileMeta)
	buf.Write([]byte(sampleType))
	buf.Write([]byte(target))
	return buf.Bytes()
}

func buildProfileMetaKey(sampleType string, target string, datetime time.Time) []byte {
	buf := bytes.NewBuffer(PrefixProfileMeta)
	buf.Write([]byte(sampleType))
	buf.Write([]byte(target))
	buf.Write(storage.BuildKey(datetime))
	return buf.Bytes()
}

func buildSampleTypeKey(sampleType string) []byte {
	buf := bytes.NewBuffer(PrefixSampleType)
	buf.Write([]byte(sampleType))
	return buf.Bytes()
}

func deleteSampleTypeKey(sampleType []byte) string {
	return strings.Replace(string(sampleType), string(PrefixSampleType), "", 1)
}

func buildTargetKey(target string) []byte {
	buf := bytes.NewBuffer(PrefixTarget)
	buf.Write([]byte(target))
	return buf.Bytes()
}

func deleteTargetKey(target []byte) string {
	return strings.Replace(string(target), string(PrefixTarget), "", 1)
}
