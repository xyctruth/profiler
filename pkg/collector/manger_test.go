package collector

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xyctruth/profiler/pkg/storage/badger"
	"github.com/xyctruth/profiler/pkg/utils"
	yaml "gopkg.in/yaml.v2"
)

func init() {
	utils.RegisterPProf()
}

func TestManger(t *testing.T) {
	err := os.RemoveAll("./data/manger")
	require.Equal(t, nil, err)
	store := badger.NewStore("./data/manger")

	manger := NewManger(store)
	config := &CollectorConfig{}
	yaml.Unmarshal([]byte(configStr), config)
	manger.Load(*config)
	require.Equal(t, 2, len(manger.collectors))

	config.TargetConfigs["profiler-server"].Host = "localhost:1111"
	manger.Load(*config)

	config.TargetConfigs["server2"].Host = "localhost:9000"
	config.TargetConfigs["server2"].ProfileConfigs["heap"] = &ProfileConfig{
		Enable: utils.Bool(true),
		Path:   "/test/path?s=123",
	}
	manger.Load(*config)

	delete(config.TargetConfigs, "profiler-server")
	manger.Load(*config)
	require.Equal(t, 1, len(manger.collectors))

	manger.Stop()

}
