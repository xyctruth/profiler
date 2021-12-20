package collector

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

func TestLoadConfig(t *testing.T) {
	LoadConfig("../../collector.yaml", func(config CollectorConfig) {
		require.NotEqual(t, config, nil)
		require.Equal(t, len(config.TargetConfigs), 2)

		serverConfig, ok := config.TargetConfigs["profiler-server"]
		require.Equal(t, ok, true)
		require.Equal(t, serverConfig.Interval, 15*time.Second)
		require.Equal(t, serverConfig.Expiration, int64(0))
		require.Equal(t, serverConfig.Host, "localhost:9000")
		require.Equal(t, len(serverConfig.ProfileConfigs), 0)

		serverConfig, ok = config.TargetConfigs["server2"]
		require.Equal(t, ok, true)
		require.Equal(t, len(serverConfig.ProfileConfigs), 3)
	})

}

func TestBuildProfileConfigs(t *testing.T) {
	configStr := `
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

	config := &CollectorConfig{}
	err := yaml.Unmarshal([]byte(configStr), config)
	require.NoError(t, err)

	serverConfig, ok := config.TargetConfigs["server2"]
	require.Equal(t, ok, true)
	require.Equal(t, len(serverConfig.ProfileConfigs), 3)

	profileConfigs := buildProfileConfigs(serverConfig.ProfileConfigs)

	require.Equal(t, len(profileConfigs), 8)

	require.Equal(t, profileConfigs["fgprof"].Path, defaultProfileConfigs["fgprof"].Path)
	require.NotEqual(t, profileConfigs["fgprof"].Enable, false)

	require.Equal(t, profileConfigs["profile"].Path, "/debug/pprof/profile?seconds=15")
	require.NotEqual(t, profileConfigs["profile"].Enable, false)

	require.Equal(t, profileConfigs["heap"].Path, defaultProfileConfigs["heap"].Path)
	require.NotEqual(t, profileConfigs["heap"].Enable, true)
}
