package badger

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
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

func (s *store) GetProfile(id string) (string, []byte, error) {
	var data []byte
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(buildProfileKey(id))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			data = val
			return nil
		})
	})

	if errors.Is(err, badger.ErrKeyNotFound) {
		return "", nil, storage.ErrProfileNotFound
	}

	buf := bytes.NewBuffer(data)
	gzipReader, err := gzip.NewReader(buf)
	if err != nil {
		return "", nil, err
	}
	defer gzipReader.Close()
	b, err := ioutil.ReadAll(gzipReader)
	if err != nil && !strings.Contains(err.Error(), "unexpected EOF") {
		return "", nil, err
	}
	return gzipReader.Header.Name, b, nil
}

func (s *store) SaveProfile(name string, profileData []byte, ttl time.Duration) (string, error) {
	var compressData bytes.Buffer
	gzipWriter, _ := gzip.NewWriterLevel(&compressData, gzip.BestCompression)
	gzipWriter.Header.Name = name
	defer gzipWriter.Close()
	_, err := gzipWriter.Write(profileData)
	if err != nil {
		return "", err
	}
	err = gzipWriter.Flush()
	if err != nil {
		return "", err
	}

	id, err := s.profileSeq.Next()
	if err != nil {
		return "", err
	}
	idStr := strconv.FormatUint(id, 10)
	err = s.db.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(newProfileEntry(idStr, compressData.Bytes(), ttl))
	})

	return idStr, err
}

func (s *store) SaveProfileMeta(metas []*storage.ProfileMeta, ttl time.Duration) error {
	err := s.db.Update(func(txn *badger.Txn) error {

		now := time.Now()
		for _, meta := range metas {
			if meta.OriginalTargetName == "" {
				meta.OriginalTargetName = meta.TargetName
			}
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

			if err = txn.SetEntry(newSampleTypeEntry(meta.SampleType, ttl)); err != nil {
				return err
			}

			if err = txn.SetEntry(newTargetEntry(meta.TargetName, ttl)); err != nil {
				return err
			}

			// 添加默认target Index
			meta.Labels = append(meta.Labels, storage.Label{
				Key:   TargetLabel,
				Value: meta.OriginalTargetName,
			})

			labelEnters := newLabelEntry(meta.Labels, ttl)
			for _, entry := range labelEnters {
				if err = txn.SetEntry(entry); err != nil {
					return err
				}
			}

			indexEnters := newIndexEntry(meta.SampleType, meta.Labels, idStr, now, ttl)
			for _, entry := range indexEnters {
				if err = txn.SetEntry(entry); err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}

func (s *store) ListProfileMeta(sampleType string, startTime, endTime time.Time, filters ...storage.LabelFilter) ([]*storage.ProfileMetaByTarget, error) {
	var err error

	// Default query all target label
	if len(filters) == 0 {
		labels, err := s.ListLabel()
		if err != nil {
			return nil, err
		}
		for _, label := range labels {
			if label.Key == TargetLabel {
				filters = append(filters, storage.LabelFilter{Label: label})
			}
		}
	}

	ids, err := s.searchProfileMeta(sampleType, filters, startTime, endTime)
	if err != nil {
		return nil, err
	}

	targetMap := make(map[string][]*storage.ProfileMeta)
	err = s.db.View(func(txn *badger.Txn) error {
		for _, id := range ids {

			key := buildProfileMetaKey(id)
			item, err := txn.Get(key)
			if err != nil {
				return err
			}

			err = item.Value(func(v []byte) error {
				meta := &storage.ProfileMeta{}
				if err = meta.Decode(v); err != nil {
					return err
				}

				if metas, ok := targetMap[meta.TargetName]; ok {
					metas = append(metas, meta)
					targetMap[meta.TargetName] = metas
				} else {
					metas = make([]*storage.ProfileMeta, 0)
					metas = append(metas, meta)
					targetMap[meta.TargetName] = metas
				}
				return nil
			})
			if err != nil {
				return err
			}

		}
		return nil
	})

	res := make([]*storage.ProfileMetaByTarget, 0)
	for targetName, metas := range targetMap {
		res = append(res, &storage.ProfileMetaByTarget{TargetName: targetName, ProfileMetas: metas})
	}

	return res, err
}

func (s *store) searchProfileMeta(sampleType string, filters []storage.LabelFilter, startTime, endTime time.Time) ([]string, error) {
	ids := make([]string, 0)
	err := s.db.View(func(txn *badger.Txn) error {
		for _, filter := range filters {
			func() {

				idsByLabel := make([]string, 0)

				min := buildIndexKey(sampleType, filter.Key, filter.Value, &startTime, nil)
				max := buildIndexKey(sampleType, filter.Key, filter.Value, &endTime, nil)

				opts := badger.DefaultIteratorOptions
				opts.PrefetchSize = 1000
				opts.Prefix = buildIndexKey(sampleType, filter.Key, filter.Value, nil, nil)
				it := txn.NewIterator(opts)
				defer it.Close()

				for it.Seek(min); it.Valid(); it.Next() {
					item := it.Item()
					k := item.Key()

					if !storage.CompareKey(k, max) {
						break
					}

					id := string(k[len(min):])
					idsByLabel = append(idsByLabel, id)
				}
				if len(ids) == 0 {
					ids = idsByLabel
				} else {
					ids = filter.Policy(ids, idsByLabel)
				}
			}()
		}
		return nil
	})
	return ids, err
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

func (s *store) ListLabel() ([]storage.Label, error) {
	labels := make([]storage.Label, 0)
	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 100
		opts.Prefix = PrefixLabel
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(PrefixLabel); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			s := strings.Split(deletePrefixKey(k), "=")
			labels = append(labels, storage.Label{
				Key:   s[0],
				Value: s[1],
			})
		}
		return nil
	})
	return labels, err
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
