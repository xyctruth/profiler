package collector

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/google/pprof/profile"
	"github.com/sirupsen/logrus"
	"github.com/xyctruth/profiler/pkg/storage"
)

type Collector struct {
	TargetName string
	*TargetConfig
	exitChan        chan struct{}
	resetTickerChan chan time.Duration
	wg              *sync.WaitGroup
	httpClient      *http.Client
	mu              sync.RWMutex
	log             *logrus.Entry
	store           storage.Store
}

func newCollector(targetName string, target *TargetConfig, store storage.Store) *Collector {
	job := &Collector{
		TargetName:      targetName,
		TargetConfig:    target,
		exitChan:        make(chan struct{}),
		resetTickerChan: make(chan time.Duration),
		wg:              &sync.WaitGroup{},
		httpClient:      &http.Client{},
		log:             logrus.WithField("collector", targetName),
		store:           store,
	}
	return job
}

func (collector *Collector) run(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	wg.Add(1)
	collector.log.Info("collector run")

	gob.Register(storage.ProfileMeta{})

	ticker := time.NewTicker(collector.Interval)
	defer ticker.Stop()
	collector.scrape()

	for {
		select {
		case <-collector.exitChan:
			collector.log.Info("collector exit")
			return
		case interval := <-collector.resetTickerChan:
			ticker.Reset(interval)
		case <-ticker.C:
			collector.scrape()
		}
	}
}

func (collector *Collector) reload(target *TargetConfig) {
	collector.mu.Lock()
	defer collector.mu.Unlock()
	if collector.Interval != target.Interval {
		collector.resetTickerChan <- target.Interval
	}
	collector.TargetConfig = target
}

func (collector *Collector) exit() {
	collector.mu.Lock()
	defer collector.mu.Unlock()
	close(collector.exitChan)
}

func (collector *Collector) scrape() {
	collector.log.Info("collector scrape start")
	collector.mu.RLock()
	for profileType, profileConfig := range DefaultProfileConfigs() {
		collector.wg.Add(1)
		go collector.fetch(profileType, profileConfig)
	}
	collector.mu.RUnlock()
	collector.wg.Wait()
}

func (collector *Collector) fetch(profileType string, profileConfig *ProfileConfig) {
	defer collector.wg.Done()
	logEntry := collector.log.WithFields(logrus.Fields{"profile_type": profileType, "profile_url": profileConfig.Path})
	req, err := http.NewRequest("GET", "http://"+collector.Host+profileConfig.Path, nil)
	if err != nil {
		logEntry.WithError(err).Error("invoke task error")
		return
	}
	req.Header.Set("User-Agent", "")

	resp, err := collector.httpClient.Do(req)
	if err != nil {
		logEntry.WithError(err).Error("http request error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		logEntry.WithError(err).Error("http resp status code is", resp.StatusCode)
		return
	}

	profileBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logEntry.WithError(err).Error("read resp error")
		return
	}

	err = collector.analysis(profileType, profileBytes)
	if err != nil {
		logEntry.WithError(err).Error("analysis result error")
		return
	}
}

func (collector *Collector) analysis(profileType string, profileBytes []byte) error {
	p, err := profile.ParseData(profileBytes)
	if err != nil {
		return err
	}
	if p.SampleType == nil || len(p.SampleType) == 0 {
		return errors.New("sample type is nil")
	}

	profileID, err := collector.store.SaveProfile(profileBytes)
	if err != nil {
		collector.log.WithError(err).Error("save profile error")
		return err
	}

	metas := make([]*storage.ProfileMeta, 0, len(p.SampleType))
	for i, _ := range p.SampleType {
		meta := &storage.ProfileMeta{}
		meta.Timestamp = time.Now().UnixNano() / time.Millisecond.Nanoseconds()
		meta.ProfileID = profileID
		meta.Duration = p.DurationNanos
		meta.SampleType = p.SampleType[i].Type
		meta.SampleTypeUnit = p.SampleType[i].Unit
		meta.ProfileType = profileType
		meta.TargetName = collector.TargetName
		for _, s := range p.Sample {
			meta.Value += s.Value[i]
		}
		if len(p.SampleType) > 1 {
			meta.SampleType = fmt.Sprintf("%s_%s", profileType, p.SampleType[i].Type)
		} else {
			meta.SampleType = fmt.Sprintf("%s", profileType)
		}
		metas = append(metas, meta)
	}

	err = collector.store.SaveProfileMeta(metas)
	if err != nil {
		collector.log.WithError(err).Error("save profile meta error")
		return err
	}
	return nil
}
