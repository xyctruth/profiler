package collector

import (
	"testing"
	"time"

	"github.com/xyctruth/profiler/pkg/utils"

	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

var (
	configStr = `
  targetConfigs:

    profiler-server:
      interval: 15s
      expiration: 0  # no expiration time. unit day
      host: localhost:9000
      profileConfigs: # default scrape (profile, heap, allocs, black, mutex, fgprof)

    server2:
      interval: 20s
      expiration: 1
      host: localhost:9000
      profileConfigs: # rewrite default profile config
        fgprof:
          enable: false
        profile:
          path: /debug/pprof/profile?seconds=15
          enable: false
        heap:
          path: /debug/pprof/heap
`
)

func TestLoadConfig(t *testing.T) {
	LoadConfig("../../collector.yaml", func(config CollectorConfig) {
		require.NotEqual(t, config, nil)
		require.Equal(t, len(config.TargetConfigs), 2)

		serverConfig, ok := config.TargetConfigs["profiler-server"]
		require.Equal(t, ok, true)
		require.Equal(t, 15*time.Second, serverConfig.Interval)
		require.Equal(t, int64(0), serverConfig.Expiration)
		require.Equal(t, "localhost:9000", serverConfig.Host)
		require.Equal(t, 0, len(serverConfig.ProfileConfigs))

		serverConfig, ok = config.TargetConfigs["server2"]
		require.Equal(t, ok, true)
		require.Equal(t, 3, len(serverConfig.ProfileConfigs))
	})

}

func TestBuildProfileConfigs(t *testing.T) {
	config := &CollectorConfig{}
	err := yaml.Unmarshal([]byte(configStr), config)
	require.NoError(t, err)

	serverConfig, ok := config.TargetConfigs["server2"]
	require.Equal(t, true, ok)
	require.Equal(t, 3, len(serverConfig.ProfileConfigs))

	profileConfigs := buildProfileConfigs(serverConfig.ProfileConfigs)

	require.Equal(t, len(profileConfigs), 8)

	require.Equal(t, defaultProfileConfigs()["fgprof"].Path, profileConfigs["fgprof"].Path)
	require.Equal(t, utils.Bool(false), profileConfigs["fgprof"].Enable)

	require.Equal(t, profileConfigs["profile"].Path, "/debug/pprof/profile?seconds=15")
	require.Equal(t, utils.Bool(false), profileConfigs["profile"].Enable)

	require.Equal(t, defaultProfileConfigs()["heap"].Path, profileConfigs["heap"].Path)
	require.Equal(t, utils.Bool(true), profileConfigs["heap"].Enable)
}
