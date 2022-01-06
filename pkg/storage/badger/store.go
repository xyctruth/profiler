package badger

import (
	"errors"
	"strconv"
	"time"

	"github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
	"github.com/xyctruth/profiler/pkg/storage"
)

type store struct {
	db         *badger.DB
	opt        Options
	profileSeq *badger.Sequence
	metaSeq    *badger.Sequence
}

func NewStore(opt Options) storage.Store {
	db, err := badger.Open(
		badger.DefaultOptions(opt.Path).
			WithLoggingLevel(3).
			WithBypassLockGuard(true).
			WithValueThreshold(1 << 10))

	if err != nil {
		panic(err)
	}

	s := &store{
		db:  db,
		opt: opt,
	}
	s.profileSeq, err = s.db.GetSequence(ProfileSequence, 1000)
	if err != nil {
		panic(err)
	}

	s.metaSeq, err = s.db.GetSequence(MetaSequence, 1000)
	if err != nil {
		panic(err)
	}

	go s.GC()

	return s
}

func (s *store) GC() {
	s.gc()

	ticker := time.NewTicker(s.opt.GCInternal)
	defer ticker.Stop()
	for range ticker.C {
		s.gc()
	}
}

func (s *store) gc() {
	log.Info("store gc start")
	if err := s.db.RunValueLogGC(0.7); err != nil {
		log.WithError(err).Info("stop store gc")
		return
	}
	s.gc()
}

func (s *store) GetProfile(id string) ([]byte, error) {
	var date []byte
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(buildProfileKey(id))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			date = val
			return nil
		})
	})

	if errors.Is(err, badger.ErrKeyNotFound) {
		return nil, storage.ErrProfileNotFound
	}
	return date, err
}

func (s *store) SaveProfile(profileData []byte, ttl time.Duration) (string, error) {
	id, err := s.profileSeq.Next()
	if err != nil {
		return "", err
	}
	idStr := strconv.FormatUint(id, 10)
	err = s.db.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(newProfileEntry(idStr, profileData, ttl))
	})

	return idStr, err
}

func (s *store) SaveProfileMeta(metas []*storage.ProfileMeta, ttl time.Duration) error {
	err := s.db.Update(func(txn *badger.Txn) error {

		for _, meta := range metas {

			id, err := s.metaSeq.Next()
			if err != nil {
				return err
			}
			idStr := strconv.FormatUint(id, 10)

			var profileMetaEntry *badger.Entry
			if profileMetaEntry, err = newProfileMetaEntry(idStr, meta, ttl); err != nil {
				return err
			}
			if err = txn.SetEntry(profileMetaEntry); err != nil {
				return err
			}

			if err = txn.SetEntry(newSampleTypeEntry(meta.SampleType, meta.ProfileType, ttl)); err != nil {
				return err
			}

			if err = txn.SetEntry(newTargetEntry(meta.TargetName, ttl)); err != nil {
				return err
			}

			labelsEntry := newLabelsEntry(meta.Labels, id, ttl)
			for _, entry := range labelsEntry {
				if err = txn.SetEntry(entry); err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}

func (s *store) ListProfileMeta(sampleType string, targetFilter []string, startTime, endTime time.Time) ([]*storage.ProfileMetaByTarget, error) {
	targets := make([]*storage.ProfileMetaByTarget, 0)
	var err error
	if len(targetFilter) == 0 {
		if targetFilter, err = s.ListTarget(); err != nil {
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

				err = item.Value(func(v []byte) error {
					meta := &storage.ProfileMeta{}
					if err = meta.Decode(v); err != nil {
						return err
					}
					target.ProfileMetas = append(target.ProfileMetas, meta)
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
	return targets, err
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
			sampleTypes = append(sampleTypes, deletePrefixKey(k))
		}

		return nil
	})

	return sampleTypes, err
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
				sampleTypes[string(v)] = append(sampleTypes[string(v)], deletePrefixKey(k))
				return nil
			})

			if err != nil {
				return err
			}
		}
		return nil
	})

	return sampleTypes, err
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
			targets = append(targets, deletePrefixKey(k))
		}
		return nil
	})

	return targets, err
}

func (s *store) ListLabel() ([]string, error) {
	targets := make([]string, 0)
	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 100
		opts.Prefix = PrefixLabel
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(PrefixLabel); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			targets = append(targets, deletePrefixKey(k))
		}
		return nil
	})

	return targets, err
}

func (s *store) Release() {
	if err := s.profileSeq.Release(); err != nil {
		log.WithError(err).Error("store release")
		return
	}

	if err := s.metaSeq.Release(); err != nil {
		log.WithError(err).Error("store release")
		return
	}

	if err := s.db.Close(); err != nil {
		log.WithError(err).Error("store close")
		return
	}

	log.Info("store release")
}
