package redis

import (
	"time"

	"github.com/go-redis/cache/v7"
)

func NewUserSession(uuid string) (string, error) {
	const n = 256

	key, err := generateRandomString(n)
	if err != nil {
		return "", err
	}

	session := cache.Item{
		Key:        key,
		Object:     uuid,
		Func: func() (i interface{}, err error) {
			return uuid, nil
		},
		Expiration: sessionLife * time.Second,
	}

	return key, Codec.Once(&session)
}
