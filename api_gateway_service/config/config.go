package config

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "api gateway microservice config path")
}

type Config struct {
	Http Http `mapstructure:"http"`
	Grpc Grpc `mapstructure:"grpc"`
}

type Http struct {
	Port string `mapstructure:"port"`
}

type Grpc struct {
	StorageServicePort string `mapstructure:"storageServicePort"`
}

func InitConfig() (*Config, error) {
	if configPath == "" {
		configPathFromEnv := os.Getenv("CONFIG_PATH")
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			getwd, err := os.Getwd()
			if err != nil {
				return nil, errors.Wrap(err, "os.Getwd")
			}
			configPath = fmt.Sprintf("%s/config/config.yaml", getwd)
		}
	}

	cfg := &Config{}
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort != "" {
		cfg.Http.Port = httpPort
	}

	storageServicePort := os.Getenv("STORAGE_SERVICE_PORT")
	if storageServicePort != "" {
		cfg.Grpc.StorageServicePort = storageServicePort
	}

	return cfg, nil
}
