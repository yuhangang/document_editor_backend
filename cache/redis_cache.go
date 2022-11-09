package cache

import "github.com/go-redis/redis/v9"

func initRedis() (*redis.Client, error) {
	opt, err := redis.ParseURL("redis://<user>:<pass>@localhost:6379/<db>")
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opt)
	return rdb, nil
}
