package main

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/nexus-minimal/")
	viper.AddConfigPath("$HOME/.nexus-minimal")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		log.WithError(err).Warn("Reading config file")
	}
}

// HttpConfig keeps the configuration for HTTP server
type HttpConfig struct {
	Address     string
	HTTPS       bool
	Certificate string
	Key         string
	Username    string
	Password    string
}

// StorageConfig keeps configuration for cloud storage backend
type StorageConfig struct {
	Enabled   bool
	BucketURL string
	Timeout   time.Duration
}

// CacheConfig keeps the configuration for local file system cache
type CacheConfig struct {
	BaseDir string
}

// RepositoryConfig keeps the configuration for remote artifacts repository
type RepositoryConfig struct {
	URL     string
	Timeout time.Duration
}

// Config keeps the entire configuration
type Config struct {
	HTTP       HttpConfig
	Storage    StorageConfig
	Cache      CacheConfig
	Repository RepositoryConfig
}

// NewConfig parse the configuration from file and returns a configuration object
func NewConfig() Config {
	initConfig()

	config := Config{}
	config.HTTP.Address = viper.GetString("http.addr")
	config.HTTP.HTTPS = viper.GetBool("http.https")
	config.HTTP.Certificate = viper.GetString("http.crt")
	config.HTTP.Key = viper.GetString("http.key")
	config.HTTP.Username = viper.GetString("http.username")
	config.HTTP.Password = viper.GetString("http.password")

	config.Storage.Enabled = viper.GetBool("storage.enabled")
	config.Storage.BucketURL = viper.GetString("storage.bucket_url")
	config.Storage.Timeout = viper.GetDuration("storage.timeout")
	if config.Storage.Timeout == 0 {
		config.Storage.Timeout = 5 * time.Minute
	}
	config.Cache.BaseDir = viper.GetString("cache.base_dir")
	if config.Cache.BaseDir == "" {
		config.Cache.BaseDir = "./.nexus"
	}

	config.Repository.URL = viper.GetString("repository.url")
	if config.Repository.URL == "" {
		config.Repository.URL = "https://repo1.maven.org/maven2"
	}
	config.Repository.Timeout = viper.GetDuration("repository.timeout")
	if config.Repository.Timeout == 0 {
		config.Repository.Timeout = 1 * time.Minute
	}

	return config
}
