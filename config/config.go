package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config представляет конфигурацию приложения.
type Config struct {
	Server  Server  `yaml:"server"`
	MongoDB MongoDB `yaml:"mongodb"`
}

type MongoDB struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func New(path string) (Config, error) {
	viper.SetConfigFile(path)

	var config Config

	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("failed to read config: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}
