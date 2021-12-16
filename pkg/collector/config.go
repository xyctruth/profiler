package collector

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	defaultProfileConfigs = map[string]*ProfileConfig{
		"profile": {
			Path:   "/debug/pprof/profile?seconds=10",
			Enable: true,
		},
		"fgprof": {
			Path:   "/debug/fgprof?seconds=10",
			Enable: true,
		},
		"mutex": {
			Path:   "/debug/pprof/mutex",
			Enable: true,
		},
		"heap": {
			Path:   "/debug/pprof/heap",
			Enable: true,
		},
		"goroutine": {
			Path:   "/debug/pprof/goroutine",
			Enable: true,
		},
		"allocs": {
			Path:   "/debug/pprof/allocs",
			Enable: true,
		},
		"block": {
			Path:   "/debug/pprof/block",
			Enable: true,
		},
		"threadcreate": {
			Path:   "/debug/pprof/threadcreate",
			Enable: true,
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
	Host           string
}

type ProfileConfig struct {
	Path   string
	Enable bool
}

func buildProfileConfigs(profileConfig map[string]*ProfileConfig) map[string]*ProfileConfig {
	if profileConfig == nil {
		return defaultProfileConfigs
	}

	for profileName, defaultConfig := range defaultProfileConfigs {
		if config, ok := profileConfig[profileName]; ok {
			if config.Path == "" {
				config.Path = defaultConfig.Path
			}
		} else {
			profileConfig[profileName] = defaultConfig
		}
	}
	return profileConfig
}
