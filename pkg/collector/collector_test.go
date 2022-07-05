package collector

import (
	"io/ioutil"
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
	c := &Config{}
	yaml.Unmarshal([]byte(generalConfigYAML), c)
	config := c.Collector
	collector := newCollector("profiler-server", config.TargetConfigs["profiler-server"], nil, &sync.WaitGroup{})
	require.NotEqual(t, nil, collector)
	require.Equal(t, collector.Interval, 2*time.Second)
	require.Equal(t, collector.Expiration, time.Duration(0))
	require.Equal(t, collector.Instances, []string{"localhost:9000"})
	require.Equal(t, len(collector.ProfileConfigs), 9)
}

func TestCollectorReload(t *testing.T) {
	dir, err := ioutil.TempDir("./", "temp-*")
	require.Equal(t, nil, err)
	defer os.RemoveAll(dir)

	store := badger.NewStore(badger.DefaultOptions(dir))
	defer store.Release()

	wg := &sync.WaitGroup{}
	c := &Config{}
	yaml.Unmarshal([]byte(generalConfigYAML), c)
	config := c.Collector

	targetConfig := config.TargetConfigs["server2"]
	collector := newCollector("server2", targetConfig, store, wg)

	collector.run()
	defer func() {
		collector.exit()
		wg.Wait()
	}()

	targetConfig.Interval = 1 * time.Second
	collector.reload(targetConfig)
	require.Equal(t, collector.Interval, 1*time.Second)

	collector.reload(targetConfig)
	require.Equal(t, collector.Interval, 1*time.Second)

	targetConfig.Expiration = 200 * time.Second
	collector.reload(targetConfig)
	require.Equal(t, 200*time.Second, collector.Expiration)

	targetConfig.ProfileConfigs["fgprof"] = ProfileConfig{
		Enable: utils.Bool(true),
		Path:   "/test/path?s=123",
	}

	collector.reload(targetConfig)
	require.Equal(t, utils.Bool(true), collector.ProfileConfigs["fgprof"].Enable)
	require.Equal(t, "/test/path?s=123", collector.ProfileConfigs["fgprof"].Path)
	time.Sleep(1 * time.Second)
}

func TestCollectorRun(t *testing.T) {
	dir, err := ioutil.TempDir("./", "temp-*")
	require.Equal(t, nil, err)
	defer os.RemoveAll(dir)
	store := badger.NewStore(badger.DefaultOptions(dir))
	defer store.Release()

	wg := &sync.WaitGroup{}

	c := &Config{}
	yaml.Unmarshal([]byte(generalConfigYAML), c)
	config := c.Collector

	targetConfig := config.TargetConfigs["profiler-server"]
	collector := newCollector("profiler-server", targetConfig, store, wg)

	collector.run()

	time.Sleep(1 * time.Second)

	collector.exit()
	wg.Wait()

	targets, err := store.ListTarget()
	require.Equal(t, nil, err)
	require.Equal(t, 1, len(targets))

	sampleTypes, err := store.ListSampleType()
	require.Equal(t, 19, len(sampleTypes))

	labels, err := store.ListLabel()
	require.Equal(t, 3, len(labels))
}
