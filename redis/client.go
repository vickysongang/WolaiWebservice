package redis

import (
	"gopkg.in/redis.v3"

	"WolaiWebservice/config"
)

var (
	redisClient *redis.Client

	RedisFailErr error
)

func Initialize() error {
	var err error

	redisClient = redis.NewClient(
		&redis.Options{
			Addr:     config.Env.Redis.Host + config.Env.Redis.Port,
			Password: config.Env.Redis.Password,
			DB:       config.Env.Redis.Db,
			PoolSize: config.Env.Redis.PoolSize,
		})
	_, err = redisClient.Ping().Result()

	RedisFailErr = err

	return err
}
