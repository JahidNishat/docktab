package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	DockerHost string
}

func LoadEnv() error {
	return godotenv.Load()
}

func New() *Config {
	viper.SetDefault("docker.host", "")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("DOCKTAB")

	return &Config{
		DockerHost: viper.GetString("docker.host"),
	}
}
