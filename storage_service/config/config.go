package config

import (
	"flag"
	"fmt"
	"github.com/paw1a/grpc-media-converter/storage_service/pkg/minio"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "storage microservice config path")
}

type Config struct {
	GRPC  GRPC          `mapstructure:"grpc"`
	Minio *minio.Config `mapstructure:"minio"`
}

type GRPC struct {
	Port string `mapstructure:"port"`
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

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort != "" {
		cfg.GRPC.Port = grpcPort
	}

	return cfg, nil
}
