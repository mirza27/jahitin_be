package util

import "github.com/spf13/viper"

type Config struct {
	AppName    string `mapstructure:"APP_NAME"`
	AppVersion string `mapstructure:"APP_VERSION" `
	Debug      bool   `mapstructure:"DEBUG"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
}

func LoadConfig(path string) (config *Config, err error) {

	viper.AddConfigPath(path)
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return

}
