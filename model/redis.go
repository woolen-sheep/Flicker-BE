package model

import (
	"github.com/go-redis/redis/v7"
	"github.com/woolen-sheep/Flicker-BE/config"
)

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(
		&redis.Options{
			Addr:     config.C.Redis.Addr,
			Password: config.C.Redis.Password,
			DB:       config.C.Redis.DB,
		},
	)

	err := redisClient.Ping().Err()
	if err != nil {
		panic(err)
	}
}
