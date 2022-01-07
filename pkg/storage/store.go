package storage

import (
	"errors"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

type Store interface {
	// GetProfile Get profile binaries by profile id, return profile binaries
	GetProfile(id string) ([]byte, error)
	// SaveProfile Save profile，return profile id
	// data binaries file
	// ttl profile expiration time
	SaveProfile(data []byte, ttl time.Duration) (string, error)

	// SaveProfileMeta Save profile meta data
	SaveProfileMeta(metas []*ProfileMeta, ttl time.Duration) error
	// SaveProfileMeta Get profile meta data list
	ListProfileMeta(sampleType string, targetFilter []string, startTime, endTime time.Time) ([]*ProfileMetaByTarget, error)

	// ListSampleType Get collected sample types list (heap_alloc_objects ,heap_alloc_space ,heap_inuse_objects ,heap_inuse_space...)
	ListSampleType() ([]string, error)
	// ListGroupSampleType Get collected sample types list grouped by profile types (heap,goroutine...)
	ListGroupSampleType() (map[string][]string, error)
	// ListTarget  Get collection target list
	ListTarget() ([]string, error)
	// ListLabels  Get collection target labels list
	ListLabel() ([]string, error)

	// Release
	Release()
}

type TargetLabels map[string]string

func (t TargetLabels) ToArray() []Label {
	labels := make([]Label, 0, len(t))
	for k, v := range t {
		labels = append(labels, Label{
			Key:   k,
			Value: v,
		})
	}
	return labels
}

type Label struct {
	Key   string
	Value string
}

type ProfileMetaByTarget struct {
	TargetName   string
	ProfileMetas []*ProfileMeta
}

type ProfileMeta struct {
	ProfileID      string
	ProfileType    string
	SampleType     string
	TargetName     string
	SampleTypeUnit string
	Value          int64
	Timestamp      int64
	Duration       int64
	Labels         []Label
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
