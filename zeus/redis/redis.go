package redis

import (
	"zeus/env"
	"github.com/go-redis/redis"
	"time"
)

const Nil = redis.Nil

var client *redis.Client

func Init() (*redis.Client, error) {
	client = redis.NewClient(&redis.Options{
		Addr:     env.RedisHostName + ":" + env.RedisPort,
		Password: env.RedisPassword,
		DB:       0, // use default DB
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return client, nil
}

func GetClient() *redis.Client {
	return client
}

func NewUserSession(uuid string) (string, error) {
	const n = 256

	key, err := generateRandomString(n)
	if err != nil {
		return "", err
	}

	return key, client.Set(key, uuid, 60*time.Second).Err()
}
