package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL       string `mapstructure:"DATABASE_URL"`
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
}

func LoadConfig(path string) (config *Config, err error) {
	viper.SetConfigName("dev")
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err = viper.Unmarshal(&config); err != nil {
		return
	}

	return
}