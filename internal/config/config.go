package config

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	viper *viper.Viper
}

const ConfigFileName = "ship-cloud-auth"

func New() *Config {

	if err := godotenv.Load(); err != nil {
		log.Warn("Error occurred while loading .env file", "error", err)
	}

	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()
	v.SetConfigName(ConfigFileName)
	v.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Warn("Error occurred while reading config file:", err)
	}

	return &Config{
		viper: v,
	}
}

func (c *Config) Port() int {
	return c.viper.GetInt("port")
}

func (c *Config) DatabaseURL() string {
	return c.viper.GetString("database-connection-string")
}

func (c *Config) RootUserEmail() string {
	viper.SetDefault("root-user-email", "root@admin.com")
	return c.viper.GetString("root-user-email")
}

func (c *Config) RootUserPassword() string {
	viper.SetDefault("root-user-password", "admin123")
	return c.viper.GetString("root-user-password")
}

func (c *Config) JWTSecret() jwt.Keyfunc {

	return func(t *jwt.Token) (any, error) {
		key := c.viper.GetString("jwt-security-key")
		if key == "" {
			return nil, fmt.Errorf("security key not specified")
		}
		return []byte(key), nil
	}
}
