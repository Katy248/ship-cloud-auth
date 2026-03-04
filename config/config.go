package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var Config *viper.Viper

func Setup() {
	godotenv.Load()

	Config = viper.New()
	Config.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	Config.AutomaticEnv()
}
