package viper

import (
	"github.com/mehdy/twixter/pkg/entities"
	"github.com/spf13/viper"
)

// Ensure Config implements ConfigGetter interface.
var _ entities.ConfigGetter = &Config{}

type Config struct {
	*viper.Viper
}

func NewConfig() *Config {
	v := viper.New()
	v.AutomaticEnv()

	return &Config{Viper: v}
}
