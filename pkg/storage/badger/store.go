package badger

import (
	"strconv"
	"time"

	"github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
	"github.com/xyctruth/profiler/pkg/storage"
)

type store struct {
	db   *badger.DB
	path string
	seq  *badger.Sequence
}

func NewStore(path string) storage.Store {
	db, err := badger.Open(badger.DefaultOptions(path).WithLoggingLevel(3).WithBypassLockGuard(true))
	if err != nil {
		panic(err)
	}

	s := &store{
		db:   db,
		path: path,
	}
	s.seq, err = s.db.GetSequence(Sequence, 1000)
	if err != nil {
		panic(err)
	}

	go s.GC()

	return s
}

func (s *store) GC() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
	again:
		log.Info("store gc start")
		err := s.db.RunValueLogGC(0.7)
		if err == nil {
			goto again
		} else {
			log.WithError(err).Info("store gc error")
		}
	}
}

func (s *store) GetProfile(id string) ([]byte, error) {
	var date []byte
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(buildProfileKey(id))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			date = val
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	return date, err
}

func (s *store) SaveProfile(profileData []byte) (uint64, error) {
	id, err := s.seq.Next()
	if err != nil {
		return 0, err
	}

	err = s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(buildProfileKey(strconv.FormatUint(id, 10)), profileData)

	})

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *store) SaveProfileMeta(metas []*storage.ProfileMeta) error {
	return s.db.Update(func(txn *badger.Txn) error {
		for _, meta := range metas {
			err := txn.Set(buildSampleTypeKey(meta.SampleType), []byte(meta.ProfileType))
			if err != nil {
				return err
			}

			err = txn.Set(buildTargetKey(meta.TargetName), []byte(meta.TargetName))
			if err != nil {
				return err
			}

			metaBytes, err := meta.Encode()
			if err != nil {
				return err
			}

			err = txn.Set(buildProfileMetaKey(meta.SampleType, meta.TargetName, time.Now()), metaBytes)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *store) ListProfileMeta(sampleType string, targetFilter []string, startTime, endTime time.Time) ([]*storage.ProfileMetaByTarget, error) {
	targets := make([]*storage.ProfileMetaByTarget, 0)
	var err error
	if len(targetFilter) == 0 {
		targetFilter, err = s.ListTarget()
		if err != nil {
			return nil, err
		}
	}
	err = s.db.View(func(txn *badger.Txn) error {
		for _, targetName := range targetFilter {
			target := &storage.ProfileMetaByTarget{TargetName: targetName, ProfileMetas: make([]*storage.ProfileMeta, 0)}
			min := buildProfileMetaKey(sampleType, targetName, startTime)
			max := buildProfileMetaKey(sampleType, targetName, endTime)

			opts := badger.DefaultIteratorOptions
			opts.PrefetchSize = 100
			opts.Prefix = buildBaseProfileMetaKey(sampleType, targetName)
			it := txn.NewIterator(opts)
			defer it.Close()
			for it.Seek(min); it.Valid(); it.Next() {
				item := it.Item()
				k := item.Key()
				if !storage.CompareKey(k, max) {
					break
				}
				err := item.Value(func(v []byte) error {
					sample := &storage.ProfileMeta{}
					err = sample.Decode(v)
					if err != nil {
						return err
					}
					target.ProfileMetas = append(target.ProfileMetas, sample)
					return nil
				})
				if err != nil {
					return err
				}
			}

			if len(target.ProfileMetas) > 0 {
				targets = append(targets, target)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return targets, nil
}

func (s *store) ListSampleType() ([]string, error) {
	sampleTypes := make([]string, 0)
	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 100
		opts.Prefix = PrefixSampleType
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(PrefixSampleType); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			sampleTypes = append(sampleTypes, deleteSampleTypeKey(k))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return sampleTypes, nil
}

func (s *store) ListGroupSampleType() (map[string][]string, error) {
	sampleTypes := make(map[string][]string)

	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 100
		opts.Prefix = PrefixSampleType
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(PrefixSampleType); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				if _, ok := sampleTypes[string(v)]; !ok {
					sampleTypes[string(v)] = make([]string, 0, 5)
				}
				sampleTypes[string(v)] = append(sampleTypes[string(v)], deleteSampleTypeKey(k))
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return sampleTypes, nil
}

func (s *store) ListTarget() ([]string, error) {
	targets := make([]string, 0)

	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 100
		opts.Prefix = PrefixTarget
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(PrefixTarget); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			targets = append(targets, deleteTargetKey(k))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return targets, nil
}

func (s *store) Clear(targetName string, days int64) error {
	sampleTypes, err := s.ListSampleType()
	if err != nil {
		return err
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	ago := today.Add(-time.Hour * 24 * time.Duration(days))

	err = s.db.Update(func(txn *badger.Txn) error {
		for _, sampleType := range sampleTypes {
			max := buildProfileMetaKey(sampleType, targetName, ago)
			opts := badger.DefaultIteratorOptions
			opts.PrefetchSize = 100
			opts.Prefix = buildBaseProfileMetaKey(sampleType, targetName)
			it := txn.NewIterator(opts)
			defer it.Close()
			for it.Rewind(); it.Valid(); it.Next() {
				item := it.Item()
				k := item.Key()
				var profileID uint64
				if !storage.CompareKey(k, max) {
					break
				}
				err = item.Value(func(v []byte) error {
					sample := &storage.ProfileMeta{}
					err = sample.Decode(v)
					if err != nil {
						return err
					}
					profileID = sample.ProfileID
					return nil
				})
				if err != nil {
					return err
				}

				err = s.delete(k, buildProfileKey(strconv.FormatUint(profileID, 10)))
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *store) delete(keys ...[]byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		for _, key := range keys {
			err := txn.Delete(key)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *store) Release() {
	err := s.seq.Release()
	if err != nil {
		log.WithError(err).Error("store release ")
		return
	}
	log.Info("store release ")
}
