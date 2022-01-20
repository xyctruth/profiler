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
	PrefixIndex       = []byte{0x86}
)

// 内置label
const TargetLabel = "_target"

func deletePrefixKey(key []byte) string {
	return string(key[1:])
}

func buildProfileKey(id string) []byte {
	var buf bytes.Buffer
	buf.Grow(len(PrefixProfiles) + len(id))
	buf.Write(PrefixProfiles)
	buf.WriteString(id)
	return buf.Bytes()
}

func buildProfileMetaKey(id string) []byte {
	var buf bytes.Buffer
	buf.Grow(len(PrefixProfileMeta) + len(id))
	buf.Write(PrefixProfileMeta)
	buf.WriteString(id)
	return buf.Bytes()
}

func buildSampleTypeKey(sampleType string) []byte {
	var buf bytes.Buffer
	buf.Grow(len(PrefixSampleType) + len(sampleType))
	buf.Write(PrefixSampleType)
	buf.WriteString(sampleType)
	return buf.Bytes()
}

func buildTargetKey(target string) []byte {
	var buf bytes.Buffer
	buf.Grow(len(PrefixTarget) + len(target))
	buf.Write(PrefixTarget)
	buf.WriteString(target)
	return buf.Bytes()
}

func buildLabelKey(key, val string) []byte {
	var buf bytes.Buffer
	buf.Grow(len(PrefixLabel) + len(key) + len("=") + len(val))
	buf.Write(PrefixLabel)
	buf.WriteString(key)
	buf.WriteString("=")
	buf.WriteString(val)
	return buf.Bytes()
}

func buildIndexKey(sampleType, key, val string, createAt *time.Time, id *string) []byte {
	var createAtBytes, idBytes []byte
	if createAt != nil {
		createAtBytes = storage.BuildTimeKey(*createAt)
	}
	if id != nil {
		idBytes = []byte(*id)
	}

	var buf bytes.Buffer
	buf.Grow(len(PrefixIndex) + len(sampleType) + len(key) + len("=") + len(val) + len(createAtBytes) + len(idBytes))

	buf.Write(PrefixIndex)
	buf.WriteString(sampleType)
	buf.WriteString(key)
	buf.WriteString("=")
	buf.WriteString(val)
	buf.Write(createAtBytes)
	buf.Write(idBytes)

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

func newSampleTypeEntry(sampleType string, ttl time.Duration) *badger.Entry {
	entry := badger.NewEntry(buildSampleTypeKey(sampleType), nil)
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

func newLabelEntry(labels []storage.Label, ttl time.Duration) []*badger.Entry {
	entries := make([]*badger.Entry, 0, len(labels))
	if len(labels) == 0 {
		return entries
	}
	for _, l := range labels {
		entry := badger.NewEntry(buildLabelKey(l.Key, l.Value), nil)
		if ttl > 0 {
			entry = entry.WithTTL(ttl)
		}
		entries = append(entries, entry)
	}
	return entries
}

func newIndexEntry(sampleType string, labels []storage.Label, id string, createAt time.Time, ttl time.Duration) []*badger.Entry {
	entries := make([]*badger.Entry, 0, len(labels))
	if len(labels) == 0 {
		return entries
	}
	for _, l := range labels {
		entry := badger.NewEntry(buildIndexKey(sampleType, l.Key, l.Value, &createAt, &id), nil)
		if ttl > 0 {
			entry = entry.WithTTL(ttl)
		}
		entries = append(entries, entry)
	}
	return entries
}
