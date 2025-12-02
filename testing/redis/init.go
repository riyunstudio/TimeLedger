package redis

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func Initialize() (*redis.Client, *miniredis.Miniredis, error) {
	mr, err := miniredis.Run()
	if err != nil {
		return nil, nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client, mr, nil
}
