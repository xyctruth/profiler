package storage

import (
	"errors"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

type Store interface {
	// GetProfile Get profile binaries by profile id, return profile binaries
	GetProfile(id string) (string, []byte, error)

	// SaveProfile Save profile，return profile id
	// binary data
	// ttl profile expiration time
	SaveProfile(name string, data []byte, ttl time.Duration) (string, error)

	// SaveProfileMeta Save profile meta data
	SaveProfileMeta(metas []*ProfileMeta, ttl time.Duration) error

	// ListProfileMeta Get profile mete data list
	ListProfileMeta(sampleType string, startTime, endTime time.Time, filters ...LabelFilter) ([]*ProfileMetaByTarget, error)

	// ListSampleType Get collected sample types list (heap_alloc_objects ,heap_alloc_space ,heap_inuse_objects ,heap_inuse_space...)
	ListSampleType() ([]string, error)

	// ListTarget  Get collection target list
	ListTarget() ([]string, error)

	// ListLabel  Get collection target labels list
	ListLabel() ([]Label, error)

	// Release 释放 Store
	Release()
}

type ProfileMeta struct {
	ProfileID          string
	ProfileType        string
	SampleType         string
	TargetName         string // target + "/" + instance
	OriginalTargetName string
	SampleTypeUnit     string
	Value              int64
	Timestamp          int64
	Duration           int64
	Labels             []Label
}

func (meta *ProfileMeta) Encode() ([]byte, error) {
	b, err := msgpack.Marshal(meta)
	if len(b) > (1 << 10) {
		return nil, errors.New("meta size > (1 << 10) , badger WithValueThreshold is 1kb")
	}
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (meta *ProfileMeta) Decode(v []byte) error {
	return msgpack.Unmarshal(v, meta)
}

type Label struct {
	Key   string
	Value string
}

type ProfileMetaByTarget struct {
	TargetName   string
	ProfileMetas []*ProfileMeta
}
