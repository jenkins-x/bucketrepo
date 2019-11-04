package main

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func initConfig(configPath string) {
	viper.SetConfigName("config")
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		log.WithError(err).Warn("Reading config file")
	}
}

// HTTPConfig keeps the configuration for HTTP server
type HTTPConfig struct {
	Address     string `mapstructure:"addr"`
	HTTPS       bool   `mapstructure:"https"`
	Certificate string `mapstructure:"crt"`
	Key         string `mapstructure:"key"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	ChartPath   string `mapstructure:"chartPath"`
}

// StorageConfig keeps configuration for cloud storage backend
type StorageConfig struct {
	Enabled   bool          `mapstructure:"enabled"`
	BucketURL string        `mapstructure:"bucket_url"`
	Timeout   time.Duration `mapstructure:"timeout"`
	Prefix    string        `mapstructure:"prefix"`
}

// CacheConfig keeps the configuration for local file system cache
type CacheConfig struct {
	BaseDir string `mapstructure:"base_dir"`
}

// RepositoryConfig keeps the configuration for remote artifacts repository
type RepositoryConfig struct {
	URL     string        `mapstructure:"url"`
	Timeout time.Duration `mapstructure:"timeout"`
	Header  http.Header   `mapstructure:"header"`
}

// Config keeps the entire configuration
type Config struct {
	HTTP         HTTPConfig         `mapstructure:"http"`
	Storage      StorageConfig      `mapstructure:"storage"`
	Cache        CacheConfig        `mapstructure:"cache"`
	Repositories []RepositoryConfig `mapstructure:"repositories"`
}

// NewConfig parse the configuration from file and returns a configuration object
func NewConfig(configPath string) Config {
	initConfig(configPath)

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable determine configuration, %v", err)
	}

	if config.Storage.Timeout == 0 {
		config.Storage.Timeout = 5 * time.Minute
	}

	if config.Cache.BaseDir == "" {
		config.Cache.BaseDir = "./.bucketrepo"
	}

	if len(config.Repositories) == 0 {
		config.Repositories = []RepositoryConfig{RepositoryConfig{"https://repo1.maven.org/maven2", 1 * time.Minute, nil}}
	}
	for i := range config.Repositories {
		if config.Repositories[i].Timeout == 0 {
			config.Repositories[i].Timeout = 1 * time.Minute
		}
	}

	if log.IsLevelEnabled(log.InfoLevel) {
		b, err := yaml.Marshal(config)
		if err == nil {
			log.Infof("Configuration: %s", string(b))
		}
	}

	return config
}
