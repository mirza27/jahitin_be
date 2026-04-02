package config

import (
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	AppName    string `mapstructure:"APP_NAME"`
	AppVersion string `mapstructure:"APP_VERSION" `
	AppPort    int    `mapstructure:"APP_PORT"`
	Debug      bool   `mapstructure:"DEBUG"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
}

func LoadConfig(path string) (config *Config, err error) {
	v := viper.New()

	v.SetConfigFile(filepath.Join(path, ".env"))
	v.SetConfigType("env")
	v.AutomaticEnv()

	err = v.ReadInConfig()
	if err != nil {
		return
	}

	var cfg Config
	err = v.Unmarshal(&cfg)
	if err != nil {
		return
	}

	config = &cfg
	return

}
