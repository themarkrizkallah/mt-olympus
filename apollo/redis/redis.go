package redis

import (
	"encoding/json"

	"github.com/go-redis/cache/v7"
	"github.com/go-redis/redis/v7"

	"apollo/env"
)

const (
	Nil         = redis.Nil
	sessionLife = 300
)

var (
	client *redis.Client
	Codec *cache.Codec
)

func Init() (*redis.Client, error) {
	client = redis.NewClient(&redis.Options{
		Addr:     env.RedisHost + ":" + env.RedisPort,
		Password: env.RedisPassword,
		DB:       0, // use default DB
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	Codec = &cache.Codec{
		Redis:     client,
		Marshal: func(i interface{}) (bytes []byte, err error) {
			return json.Marshal(i)
		},
		Unmarshal: func(bytes []byte, i interface{}) error {
			return json.Unmarshal(bytes, &i)
		},
	}
		
	return client, nil
}

func GetClient() *redis.Client {
	return client
}
