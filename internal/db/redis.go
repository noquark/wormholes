package db

import "github.com/gomodule/redigo/redis"

// Config for Redis.
type Redis struct {
	URI string `env:"REDIS_URI" envDefault:"localhost:6379"`
}

func (r *Redis) Connect() *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) { return redis.Dial("tcp", r.URI) },
	}
}
