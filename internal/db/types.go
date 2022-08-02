package db

type Type string

const (
	TypePostgres  Type = "postgres"
	TypeTimescale Type = "timescale"
	TypeRedis     Type = "redis"
)

const (
	ModeUnified   = "unified"
	ModeCreator   = "creator"
	ModeDirector  = "director"
	ModeGenerator = "generator"
)
