package collector

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/xyctruth/profiler/pkg/storage"
)

// Manger Manage multiple collectors to scraping
type Manger struct {
	collectors map[string]*Collector
	store      storage.Store
	wg         *sync.WaitGroup
	mu         sync.Mutex
}

// NewManger new Manger instance
func NewManger(store storage.Store) *Manger {
	c := &Manger{
		collectors: make(map[string]*Collector),
		store:      store,
		wg:         &sync.WaitGroup{},
	}
	return c
}

// NewManger stop Manger instance
func (manger *Manger) Stop() {
	manger.mu.Lock()
	defer manger.mu.Unlock()
	for _, c := range manger.collectors {
		c.exit()
	}
	manger.wg.Wait()
	log.Info("collector manger exit ")
}

// NewManger Loading collector configuration
// It can be called multiple times, and the collector updates the configuration
func (manger *Manger) Load(config CollectorConfig) {
	manger.mu.Lock()
	defer manger.mu.Unlock()
	// delete old collector
	for k, collector := range manger.collectors {
		if _, ok := config.TargetConfigs[k]; !ok {
			log.Info("delete collector ", k)
			collector.exit()
			delete(manger.collectors, k)

		}
	}

	for k, target := range config.TargetConfigs {
		collector, ok := manger.collectors[k]
		if !ok {
			// add collector
			log.Info("add collector ", k)
			collector := newCollector(k, target, manger.store, manger.wg)
			manger.collectors[k] = collector
			collector.run()
			continue
		}

		// update collector
		collector.reload(target)
	}
}
