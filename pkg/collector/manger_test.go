package collector

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xyctruth/profiler/pkg/storage/badger"
	"github.com/xyctruth/profiler/pkg/utils"
	yaml "gopkg.in/yaml.v2"
)

func init() {
	utils.RegisterPProf()
}

func TestManger(t *testing.T) {
	dir, err := ioutil.TempDir("./", "temp-*")
	require.Equal(t, nil, err)
	defer os.RemoveAll(dir)
	store := badger.NewStore(badger.DefaultOptions(dir))
	defer store.Release()

	manger := NewManger(store)
	c := &Config{}
	yaml.Unmarshal([]byte(generalConfigYAML), c)
	config := c.Collector
	manger.Load(config)
	require.Equal(t, 2, len(manger.collectors))

	s2 := config.TargetConfigs["server2"]
	s2.Interval = time.Second * 1
	s2.Host = "localhost:9000"
	s2.ProfileConfigs["heap"] = ProfileConfig{
		Enable: utils.Bool(true),
		Path:   "/test/path?s=123",
	}
	config.TargetConfigs["server2"] = s2
	manger.Load(config)

	delete(config.TargetConfigs, "profiler-server")
	manger.Load(config)
	require.Equal(t, 1, len(manger.collectors))
	manger.Stop()

}

func TestErrorHostManger(t *testing.T) {
	dir, err := ioutil.TempDir("./", "temp-*")
	require.Equal(t, nil, err)
	defer os.RemoveAll(dir)
	store := badger.NewStore(badger.DefaultOptions(dir))
	defer store.Release()
	manger := NewManger(store)
	config := &CollectorConfig{}
	yaml.Unmarshal([]byte(errHostConfigYAML), config)
	manger.Load(*config)
	manger.Stop()
}
