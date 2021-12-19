package collector

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/xyctruth/profiler/pkg/utils"
)

var (
	defaultProfileConfigs = map[string]*ProfileConfig{
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
)

func LoadConfig(configPath string, fn func(configmap Config)) {
	conf := viper.New()
	conf.SetConfigFile(configPath)
	conf.SetConfigType("yaml")

	err := conf.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w", err))
	}

	var config Config
	err = conf.UnmarshalKey("collector", &config)
	if err != nil {
		panic(fmt.Errorf("Fatal error config Config: %w", err))
	}

	conf.OnConfigChange(func(in fsnotify.Event) {
		var newConfig Config
		err = conf.UnmarshalKey("collector", &newConfig)
		if err != nil {
			panic(fmt.Errorf("Fatal error config Config: %w", err))
		}
		fn(newConfig)
	})
	conf.WatchConfig()
	fn(config)
}

type Config struct {
	Name          string
	TargetConfigs map[string]*TargetConfig
}

type TargetConfig struct {
	ProfileConfigs map[string]*ProfileConfig
	Interval       time.Duration
	Expiration     int64 // unit day
	Host           string
}

type ProfileConfig struct {
	Path   string
	Enable *bool
}

func buildProfileConfigs(profileConfig map[string]*ProfileConfig) map[string]*ProfileConfig {
	if profileConfig == nil {
		return defaultProfileConfigs
	}

	for key, defaultConfig := range defaultProfileConfigs {
		if config, ok := profileConfig[key]; ok {
			if config.Path == "" {
				config.Path = defaultConfig.Path
			}

			if config.Enable == nil {
				config.Enable = defaultConfig.Enable
			}
			continue
		}

		profileConfig[key] = defaultConfig
	}
	return profileConfig
}
