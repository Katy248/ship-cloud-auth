package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var Config *viper.Viper

func Setup() {
	_ = godotenv.Load()

	Config = viper.New()
	Config.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	Config.AutomaticEnv()
}

func SecurityKey() []byte {
	return []byte(Config.GetString("jwt-security-key"))
}
