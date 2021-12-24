package collector

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/xyctruth/profiler/pkg/go/v1175/trace"

	"github.com/google/pprof/profile"
	"github.com/sirupsen/logrus"
	"github.com/xyctruth/profiler/pkg/storage"
)

// Collector Collect target pprof http endpoints
type Collector struct {
	TargetName string
	TargetConfig
	exitChan        chan struct{}
	resetTickerChan chan time.Duration
	mangerWg        *sync.WaitGroup
	wg              *sync.WaitGroup
	httpClient      *http.Client
	mu              sync.RWMutex
	log             *logrus.Entry
	store           storage.Store
}

func newCollector(targetName string, target TargetConfig, store storage.Store, mangerWg *sync.WaitGroup) *Collector {
	collector := &Collector{
		TargetName:      targetName,
		TargetConfig:    target,
		exitChan:        make(chan struct{}),
		resetTickerChan: make(chan time.Duration, 1000),
		mangerWg:        mangerWg,
		wg:              &sync.WaitGroup{},
		httpClient:      &http.Client{},
		log:             logrus.WithField("collector", targetName),
		store:           store,
	}
	collector.ProfileConfigs = buildProfileConfigs(collector.ProfileConfigs)
	return collector
}

func (collector *Collector) run() {
	collector.mu.Lock()
	defer collector.mu.Unlock()

	collector.log.Info("collector run")

	collector.mangerWg.Add(1)

	go collector.scrapeLoop(collector.Interval)
}

func (collector *Collector) scrapeLoop(interval time.Duration) {
	defer collector.mangerWg.Done()
	collector.scrape()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-collector.exitChan:
			collector.log.Info("scrape loop exit")
			return
		case i := <-collector.resetTickerChan:
			ticker.Reset(i)
		case <-ticker.C:
			collector.scrape()
		}
	}
}

func (collector *Collector) reload(target TargetConfig) {
	collector.mu.Lock()
	defer collector.mu.Unlock()

	target.ProfileConfigs = buildProfileConfigs(target.ProfileConfigs)

	if reflect.DeepEqual(collector.TargetConfig, target) {
		return
	}
	collector.log.Info("reload collector ")

	if collector.Interval != target.Interval {
		collector.resetTickerChan <- target.Interval
	}

	collector.TargetConfig = target
}

func (collector *Collector) exit() {
	close(collector.exitChan)
}

func (collector *Collector) scrape() {
	collector.mu.RLock()
	defer collector.mu.RUnlock()

	collector.log.Info("collector scrape start")
	for profileType, profileConfig := range collector.ProfileConfigs {
		if *profileConfig.Enable {
			collector.wg.Add(1)
			go collector.fetch(profileType, profileConfig)
		}
	}
	collector.wg.Wait()
}

func (collector *Collector) fetch(profileType string, profileConfig ProfileConfig) {
	defer collector.wg.Done()

	logEntry := collector.log.WithFields(logrus.Fields{"profile_type": profileType, "profile_url": profileConfig.Path})
	logEntry.Info("collector start fetch")

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
		logEntry.WithError(err).Error("http resp status code is ", resp.StatusCode)
		return
	}

	profileBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logEntry.WithError(err).Error("read resp error")
		return
	}

	if profileType == "trace" {
		err = collector.analysisTrace(profileType, profileBytes)
		if err != nil {
			logEntry.WithError(err).Error("analysis result error")
			return
		}
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
	if len(p.SampleType) == 0 {
		return errors.New("sample type is nil")
	}

	if len(p.Mapping) > 0 {
		p.Mapping[0].File = collector.TargetName
	}

	b := &bytes.Buffer{}
	if err = p.Write(b); err != nil {
		return err
	}

	profileID, err := collector.store.SaveProfile(b.Bytes(), collector.Expiration)
	if err != nil {
		return err
	}

	metas := make([]*storage.ProfileMeta, 0, len(p.SampleType))
	for i := range p.SampleType {
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
			meta.SampleType = profileType
		}
		metas = append(metas, meta)
	}

	err = collector.store.SaveProfileMeta(metas, collector.Expiration)
	if err != nil {
		return err
	}
	return nil
}

func (collector *Collector) analysisTrace(profileType string, profileBytes []byte) error {
	buf := &bytes.Buffer{}
	buf.Write(profileBytes)
	_, err := trace.Parse(buf, collector.TargetName)
	if err != nil {
		return err
	}

	profileID, err := collector.store.SaveProfile(profileBytes, collector.Expiration)
	if err != nil {
		return err
	}
	metas := make([]*storage.ProfileMeta, 0, 1)
	meta := &storage.ProfileMeta{}
	meta.Timestamp = time.Now().UnixNano() / time.Millisecond.Nanoseconds()
	meta.ProfileID = profileID
	meta.SampleType = profileType
	meta.ProfileType = profileType
	meta.TargetName = collector.TargetName
	metas = append(metas, meta)

	err = collector.store.SaveProfileMeta(metas, collector.Expiration)
	if err != nil {
		return err
	}
	return nil
}
