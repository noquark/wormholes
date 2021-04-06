package redis

import (
	"time"

	"github.com/go-redis/redis/v8"
)

// Redis connection
type Redis struct {
	Address  string
	Password string
	DB       int
}

func Default() Redis {
	return Redis{
		Address:  "localhost:6379",
		Password: "",
		DB:       0,
	}
}

func (r *Redis) Connect() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         r.Address,
		Password:     r.Password,
		DB:           r.DB,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	})
}
