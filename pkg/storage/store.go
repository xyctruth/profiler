package storage

import (
	"github.com/vmihailenco/msgpack/v5"
	"time"
)

type Store interface {
	GetProfile(id string) ([]byte, error)
	SaveProfile(data []byte) (uint64, error)

	SaveProfileMeta(metas []*ProfileMeta) error
	ListProfileMeta(sampleType string, targetFilter []string, startTime, endTime time.Time) ([]*ProfileMetaByTarget, error)

	ListSampleType() ([]string, error)
	ListGroupSampleType() (map[string][]string, error)

	ListTarget() ([]string, error)

	Release()
}

type ProfileMetaByTarget struct {
	TargetName   string
	ProfileMetas []*ProfileMeta
}

type ProfileMeta struct {
	ProfileID      uint64
	Value          int64
	Timestamp      int64
	Duration       int64
	SampleTypeUnit string
	ProfileType    string
	TargetName     string
	SampleType     string
}

func (meta *ProfileMeta) Encode() ([]byte, error) {
	b, err := msgpack.Marshal(meta)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (meta *ProfileMeta) Decode(v []byte) error {
	return msgpack.Unmarshal(v, meta)
}
