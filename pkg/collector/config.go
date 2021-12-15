package collector

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
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
	Path string
}

func DefaultProfileConfigs() map[string]*ProfileConfig {
	return map[string]*ProfileConfig{
		"profile": {
			Path: "/debug/pprof/profile?seconds=10",
		},
		"mutex": {
			Path: "/debug/pprof/mutex",
		},
		"heap": {
			Path: "/debug/pprof/heap",
		},
		"goroutine": {
			Path: "/debug/pprof/goroutine",
		},
		"allocs": {
			Path: "/debug/pprof/allocs",
		},
		"block": {

			Path: "/debug/pprof/block",
		},
		"threadcreate": {
			Path: "/debug/pprof/threadcreate",
		},
	}
}
