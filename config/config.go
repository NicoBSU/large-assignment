package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort                   string `mapstructure:"app-port"`
	AppHost                   string `mapstructure:"app-host"`
	DockerVersion             string `mapstructure:"docker-version"`
	MinioContainerNamePattern string `mapstructure:"minio-container-name-pattern"`
	DockerNetworkName         string `mapstructure:"docker-network-name"`
	UserEnvVarName            string `mapstructure:"user-env-var-name"`
	PasswordEnvVarName        string `mapstructure:"password-env-var-name"`
	MinioPort                 string `mapstructure:"minio-port"`
	MinioSecureConnection     bool   `mapstructure:"minio-secure-connection"`
	BucketName                string `mapstructure:"bucket-name"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config/")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}
