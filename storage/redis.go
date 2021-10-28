package storage

import redislib "github.com/go-redis/redis"

var Redis *redislib.Client

func ConnectRedis() (err error) {
	Redis = redislib.NewClient(&redislib.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	_, err = Redis.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}
