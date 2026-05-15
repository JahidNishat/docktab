package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	DockerHost string
	Debug      bool
}

func LoadEnv() error {
	return godotenv.Load()
}

func New() *Config {
	viper.SetEnvPrefix("DOCKTAB")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("docker.host", "")
	viper.SetDefault("debug", false)

	return &Config{
		DockerHost: viper.GetString("docker.host"),
		Debug:      viper.GetBool("debug"),
	}
}
