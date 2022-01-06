package badger

import (
	"bytes"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/xyctruth/profiler/pkg/storage"
)

var (
	ProfileSequence = []byte{0x80}
	MetaSequence    = []byte{0x86}

	PrefixProfiles    = []byte{0x81}
	PrefixProfileMeta = []byte{0x82}
	PrefixSampleType  = []byte{0x83}
	PrefixTarget      = []byte{0x84}
	PrefixLabel       = []byte{0x85}
)

func deletePrefixKey(key []byte) string {
	return string(key[1:])
}

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

func buildProfileMetaKey(id string) []byte {
	buf := bytes.NewBuffer(PrefixProfileMeta)
	buf.Write([]byte(id))
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

func buildLabelKey(key, val, id string) []byte {
	buf := bytes.NewBuffer(PrefixLabel)
	buf.Write([]byte(key))
	buf.Write([]byte("="))
	buf.Write([]byte(val))
	buf.Write([]byte(id))
	return buf.Bytes()
}

func newProfileEntry(id string, val []byte, ttl time.Duration) *badger.Entry {
	entry := badger.NewEntry(buildProfileKey(id), val)
	if ttl > 0 {
		entry = entry.WithTTL(ttl)
	}
	return entry
}

func newProfileMetaEntry(id string, meta *storage.ProfileMeta, ttl time.Duration) (*badger.Entry, error) {
	metaBytes, err := meta.Encode()
	if err != nil {
		return nil, err
	}
	entry := badger.NewEntry(buildProfileMetaKey(id), metaBytes)
	if ttl > 0 {
		entry = entry.WithTTL(ttl)
	}
	return entry, nil
}

func newSampleTypeEntry(sampleType string, profileType string, ttl time.Duration) *badger.Entry {
	entry := badger.NewEntry(buildSampleTypeKey(sampleType), []byte(profileType))
	if ttl > 0 {
		entry = entry.WithTTL(ttl)
	}
	return entry
}

func newTargetEntry(target string, ttl time.Duration) *badger.Entry {
	entry := badger.NewEntry(buildTargetKey(target), nil)
	if ttl > 0 {
		entry = entry.WithTTL(ttl)
	}
	return entry
}

func newLabelsEntry(labels map[string]string, id string, ttl time.Duration) []*badger.Entry {
	entries := make([]*badger.Entry, 0, len(labels))
	if len(labels) == 0 {
		return entries
	}
	for key, val := range labels {
		entry := badger.NewEntry(buildLabelKey(key, val, id), nil)
		if ttl > 0 {
			entry = entry.WithTTL(ttl)
		}
		entries = append(entries, entry)
	}
	return entries
}
