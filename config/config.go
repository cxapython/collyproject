package config

import (
	"github.com/spf13/viper"
)

var Config *viper.Viper

func init() {
	Config = viper.New()
	Config.SetConfigName("config")
	Config.AddConfigPath(".")
	Config.SetConfigType("toml")
	if err := Config.ReadInConfig(); err != nil {
		panic(err)
	}
}