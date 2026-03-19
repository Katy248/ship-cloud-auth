package keyval

import (
	"github.com/katy248/ship-cloud-auth/config"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func Setup() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.Config.GetString("redis-address"),
		Password: config.Config.GetString("redis-password"),
		DB:       0,
	})
}
