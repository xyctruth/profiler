package collector

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xyctruth/profiler/pkg/utils"
)

func defaultProfileConfigs() map[string]ProfileConfig {
	return map[string]ProfileConfig{
		"profile": {
			Path:   "/debug/pprof/profile?seconds=10",
			Enable: utils.BoolPtr(true),
		},
		"fgprof": {
			Path:   "/debug/fgprof?seconds=10",
			Enable: utils.BoolPtr(true),
		},
		"mutex": {
			Path:   "/debug/pprof/mutex",
			Enable: utils.BoolPtr(true),
		},
		"heap": {
			Path:   "/debug/pprof/heap",
			Enable: utils.BoolPtr(true),
		},
		"goroutine": {
			Path:   "/debug/pprof/goroutine",
			Enable: utils.BoolPtr(true),
		},
		"allocs": {
			Path:   "/debug/pprof/allocs",
			Enable: utils.BoolPtr(true),
		},
		"block": {
			Path:   "/debug/pprof/block",
			Enable: utils.BoolPtr(true),
		},
		"threadcreate": {
			Path:   "/debug/pprof/threadcreate",
			Enable: utils.BoolPtr(true),
		},
	}

}

func LoadConfig(configPath string, fn func(configmap CollectorConfig)) error {
	conf := viper.New()
	conf.SetConfigFile(configPath)
	conf.SetConfigType("yaml")

	err := conf.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Fatal error config file: %w", err)
	}

	var config CollectorConfig
	err = conf.UnmarshalKey("collector", &config)
	if err != nil {
		return fmt.Errorf("Fatal error config CollectorConfig: %w", err)
	}

	conf.OnConfigChange(func(in fsnotify.Event) {
		var newConfig CollectorConfig
		err = conf.UnmarshalKey("collector", &newConfig)
		if err != nil {
			log.Info("Fatal error config CollectorConfig: %w")
		}
		fn(newConfig)
	})
	conf.WatchConfig()
	fn(config)

	return nil
}

type Config struct {
	Collector CollectorConfig `yaml:"collector"`
}

type CollectorConfig struct {
	TargetConfigs map[string]TargetConfig `yaml:"targetConfigs"`
}

type TargetConfig struct {
	ProfileConfigs map[string]ProfileConfig `yaml:"profileConfigs"`
	Interval       time.Duration            `yaml:"interval"`
	Expiration     int64                    `yaml:"expiration"` // unit day
	Host           string                   `yaml:"host"`
}

type ProfileConfig struct {
	Path   string `yaml:"path"`
	Enable *bool  `yaml:"enable"`
}

func buildProfileConfigs(profileConfig map[string]ProfileConfig) map[string]ProfileConfig {
	defaultConfigs := defaultProfileConfigs()
	if profileConfig == nil {
		return defaultConfigs
	}

	profiles := make(map[string]ProfileConfig, len(defaultConfigs))

	for key, defaultConfig := range defaultConfigs {
		if config, ok := profileConfig[key]; ok {
			if config.Path == "" {
				config.Path = defaultConfig.Path
			}
			if config.Enable == nil {
				config.Enable = defaultConfig.Enable
			}

			profiles[key] = config
			continue
		}

		profiles[key] = defaultConfig
	}
	return profiles
}
