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
	v := newViper()

	return &Config{
		DockerHost: v.GetString("docker.host"),
		Debug:      v.GetBool("debug"),
	}
}

func newViper() *viper.Viper {
	v := viper.New()

	v.SetEnvPrefix("DOCKTAB")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("docker.host", "")
	v.SetDefault("debug", false)

	return v
}
