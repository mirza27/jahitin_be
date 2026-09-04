package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	AppName        string `mapstructure:"APP_NAME"`
	AppVersion     string `mapstructure:"APP_VERSION"`
	AppPort        int    `mapstructure:"APP_PORT"`
	Debug          bool   `mapstructure:"DEBUG"`
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         int    `mapstructure:"DB_PORT"`
	DBUser         string `mapstructure:"DB_USER"`
	DBPassword     string `mapstructure:"DB_PASSWORD"`
	DBName         string `mapstructure:"DB_NAME"`
	TokenSecretKey string `mapstructure:"TOKEN_SECRET_KEY"`
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	envFilePath := filepath.Join(path, ".env")

	if _, err := os.Stat(envFilePath); err == nil {
		v.SetConfigFile(envFilePath)
		v.SetConfigType("env")
		if err := v.ReadInConfig(); err != nil {
			return nil, err
		}
	} else {
		v.AutomaticEnv()
		keys := []string{
			"APP_NAME",
			"APP_VERSION",
			"APP_PORT",
			"DEBUG",
			"DB_HOST",
			"DB_PORT",
			"DB_USER",
			"DB_PASSWORD",
			"DB_NAME",
			"TOKEN_SECRET_KEY",
		}
		for _, key := range keys {
			_ = v.BindEnv(key)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
