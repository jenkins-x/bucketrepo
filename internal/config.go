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

type HttpConfig struct {
	Address     string
	HTTPS       bool
	Certificate string
	Key         string
	Username    string
	Password    string
}

type StorageConfig struct {
	Type string

	Bucket    string
	AccessKey string
	SecretKey string

	BaseDir string
}

type RepositoryConfig struct {
	URL     string
	Timeout time.Duration
}

type Config struct {
	HTTP       HttpConfig
	Storage    StorageConfig
	Repository RepositoryConfig
}

func NewConfig() Config {
	initConfig()

	config := Config{}
	config.HTTP.Address = viper.GetString("http.addr")
	config.HTTP.HTTPS = viper.GetBool("http.https")
	config.HTTP.Certificate = viper.GetString("http.crt")
	config.HTTP.Key = viper.GetString("http.key")
	config.HTTP.Username = viper.GetString("http.username")
	config.HTTP.Password = viper.GetString("http.password")

	config.Storage.Type = viper.GetString("storage.type")
	config.Storage.Bucket = viper.GetString("storage.bucket")
	config.Storage.AccessKey = viper.GetString("storage.access_key")
	config.Storage.SecretKey = viper.GetString("storage.secret_key")

	config.Storage.BaseDir = viper.GetString("storage.base_dir")
	if config.Storage.BaseDir == "" {
		config.Storage.BaseDir = "./.nexus"
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
