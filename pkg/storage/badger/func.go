package badger

import (
	"bytes"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/dgraph-io/badger/v3"

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

func buildTargetKey(target string) []byte {
	buf := bytes.NewBuffer(PrefixTarget)
	buf.Write([]byte(target))
	return buf.Bytes()
}

func newProfileEntry(id uint64, val []byte, ttl time.Duration) *badger.Entry {
	entry := badger.NewEntry(buildProfileKey(strconv.FormatUint(id, 10)), val)

	if ttl > 0 {
		entry = entry.WithTTL(ttl)
	}
	return entry
}

func newProfileMetaEntry(meta *storage.ProfileMeta, ttl time.Duration) *badger.Entry {
	metaBytes, err := meta.Encode()
	if err != nil {
		logrus.WithError(err).Error("ProfileMeta encode error")
		return nil
	}
	entry := badger.NewEntry(buildProfileMetaKey(meta.SampleType, meta.TargetName, time.Now()), metaBytes)
	if ttl > 0 {
		entry = entry.WithTTL(ttl)
	}
	return entry
}

func newSampleTypeEntry(sampleType string, profileType string, ttl time.Duration) *badger.Entry {
	entry := badger.NewEntry(buildSampleTypeKey(sampleType), []byte(profileType))
	if ttl > 0 {
		entry = entry.WithTTL(ttl)
	}
	return entry
}

func newTargetKeyEntry(target string, ttl time.Duration) *badger.Entry {
	entry := badger.NewEntry(buildTargetKey(target), []byte{})
	if ttl > 0 {
		entry = entry.WithTTL(ttl)
	}
	return entry
}

func deleteSampleTypeKey(sampleType []byte) string {
	return strings.Replace(string(sampleType), string(PrefixSampleType), "", 1)
}

func deleteTargetKey(target []byte) string {
	return strings.Replace(string(target), string(PrefixTarget), "", 1)
}
