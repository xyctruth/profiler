package collector

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xyctruth/profiler/pkg/storage"
	"github.com/xyctruth/profiler/pkg/utils"
)

// LoadConfig watch configPath change, callback fn
func LoadConfig(configPath string, fn func(CollectorConfig)) error {
	var err error
	var config CollectorConfig

	conf := viper.New()
	conf.SetConfigFile(configPath)
	conf.SetConfigType("yaml")

	if err = conf.ReadInConfig(); err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}

	if err = conf.UnmarshalKey("collector", &config); err != nil {
		return fmt.Errorf("fatal error config CollectorConfig: %w", err)
	}

	conf.OnConfigChange(func(in fsnotify.Event) {
		var newConfig CollectorConfig
		if err = conf.UnmarshalKey("collector", &newConfig); err != nil {
			log.Error("Fatal error config CollectorConfig: %w")
			return
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
	//key TargetName
	TargetConfigs map[string]TargetConfig `yaml:"targetConfigs"`
}

type TargetConfig struct {
	//key is profile name (profile, fgprof, mutex, heap, goroutine, allocs, block, threadcreate)
	ProfileConfigs map[string]ProfileConfig `yaml:"profileConfigs"`
	Interval       time.Duration            `yaml:"interval"`
	Expiration     time.Duration            `yaml:"expiration"`
	Instances      []string                 `yaml:"instances"`
	Labels         LabelConfig              `yaml:"labels"`
}

type LabelConfig map[string]string

func (t LabelConfig) ToArray() []storage.Label {
	labels := make([]storage.Label, 0, len(t))
	for k, v := range t {
		labels = append(labels, storage.Label{
			Key:   k,
			Value: v,
		})
	}
	return labels
}

type ProfileConfig struct {
	Path   string `yaml:"path"`
	Enable *bool  `yaml:"enable"`
}

// defaultProfileConfigs The default fetching profile config
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
		"trace": {
			Path:   "/debug/pprof/trace?seconds=10",
			Enable: utils.BoolPtr(false),
		},
	}
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
