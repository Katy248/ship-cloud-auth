package keyval

import (
	"github.com/redis/go-redis/v9"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/config"
)

var RDB *redis.Client

func Setup() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.Config.GetString("redis-address"),
		Password: config.Config.GetString("redis-password"),
		DB:       0,
	})
}
