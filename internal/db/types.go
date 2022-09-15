package db

type Type string

const (
	TypePostgres  Type = "postgres"
	TypeTimescale Type = "timescale"
	TypeRedis     Type = "redis"
)
