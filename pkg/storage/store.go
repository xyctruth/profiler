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
	// data: binary profile data
	// ttl: profile expiration time
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

	// Release Store
	Release()
}

type ProfileMeta struct {
	ProfileID      string  `json:"profile_id"`
	ProfileType    string  `json:"profile_type"`
	SampleType     string  `json:"sample_type"`
	TargetName     string  `json:"target_name"`
	Instance       string  `json:"instance"`
	SampleTypeUnit string  `json:"sample_type_unit"`
	Value          int64   `json:"value"`
	Timestamp      int64   `json:"timestamp"`
	Duration       int64   `json:"duration"`
	Labels         []Label `json:"labels"`
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
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ProfileMetaByTarget struct {
	Key          string         `json:"key"`
	ProfileMetas []*ProfileMeta `json:"profile_metas"`
}
