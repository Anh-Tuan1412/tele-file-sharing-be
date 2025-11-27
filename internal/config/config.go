package config

import (
	"os"

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

	// Try to read config file, but don't fail if it doesn't exist
	// (in production on Fly.io, env vars are set directly)
	if err = viper.ReadInConfig(); err != nil {
		// Only fail if we're not in production (no FLY_APP_NAME env var)
		if os.Getenv("FLY_APP_NAME") == "" {
			return nil, err
		}
		// In production, continue without the config file
	}

	if err = viper.Unmarshal(&config); err != nil {
		return
	}

	return
}