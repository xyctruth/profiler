package collector

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xyctruth/profiler/pkg/storage/badger"
	"github.com/xyctruth/profiler/pkg/utils"
	yaml "gopkg.in/yaml.v2"
)

func TestNewCollector(t *testing.T) {
	config := &CollectorConfig{}
	yaml.Unmarshal([]byte(configStr), config)
	collector := newCollector("profiler-server", *config.TargetConfigs["profiler-server"], nil)
	require.NotEqual(t, nil, collector)
	require.Equal(t, collector.Interval, 15*time.Second)
	require.Equal(t, collector.Expiration, int64(0))
	require.Equal(t, collector.Host, "localhost:9000")
	require.Equal(t, len(collector.ProfileConfigs), 8)
}

func TestCollectorReload(t *testing.T) {
	config := &CollectorConfig{}
	yaml.Unmarshal([]byte(configStr), config)
	targetConfig := config.TargetConfigs["profiler-server"]
	collector := newCollector("profiler-server", *targetConfig, nil)
	require.NotEqual(t, nil, collector)
	require.Equal(t, 15*time.Second, collector.Interval)
	require.Equal(t, int64(0), collector.Expiration)
	require.Equal(t, "localhost:9000", collector.Host)
	require.Equal(t, 8, len(collector.ProfileConfigs))

	targetConfig.Interval = 20 * time.Second
	go func() {
		<-collector.resetTickerChan
	}()
	collector.reload(*targetConfig)
	require.Equal(t, collector.Interval, 20*time.Second)

	targetConfig.Expiration = 200
	collector.reload(*targetConfig)
	require.Equal(t, int64(200), collector.Expiration)

	targetConfig.Host = "localhost:9001"
	targetConfig.ProfileConfigs = make(map[string]*ProfileConfig)
	targetConfig.ProfileConfigs["fgprof"] = &ProfileConfig{
		Enable: utils.Bool(false),
		Path:   "/test/path?s=123",
	}
	collector.reload(*targetConfig)
	require.Equal(t, "localhost:9001", collector.Host)
	require.Equal(t, false, *targetConfig.ProfileConfigs["fgprof"].Enable)
	require.Equal(t, "/test/path?s=123", targetConfig.ProfileConfigs["fgprof"].Path)
}

func TestCollectorRun(t *testing.T) {
	err := os.RemoveAll("./data")
	require.Equal(t, nil, err)
	store := badger.NewStore("./data")

	config := &CollectorConfig{}
	yaml.Unmarshal([]byte(configStr), config)
	targetConfig := config.TargetConfigs["profiler-server"]

	collector := newCollector("profiler-server1", *targetConfig, store)

	wg := &sync.WaitGroup{}
	go collector.run(wg)

	time.Sleep(1 * time.Second)
	collector.exit()
	wg.Wait()

	targets, err := store.ListTarget()
	require.Equal(t, nil, err)
	require.Equal(t, 1, len(targets))

	sampleTypes, err := store.ListSampleType()
	require.Equal(t, 18, len(sampleTypes))
}
