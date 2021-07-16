package constants

const (
	// default
	DEFAULT_PORT = 3000
	DEFAULT_CONF = "config.yml"
	DEFAULT_DIR  = "wormholes"
	DIR_PERM     = 0775
	ENV_PREFIX   = "WH"
	// common
	DOT_DIR      = ".wormholes"
	EMPTY_STRING = ""
	// factory
	BLOOM_DB   = "bloom.db"
	MAX_LIMIT  = 1e7
	ERROR_RATE = 1e-3
	MAX_TRY    = 10
	ID_SIZE    = 7
	// pipe
	STREAMS    = 8
	BATCH_SIZE = 1e4
	CITY_DB    = "GeoLite2-City.mmdb"
	EN         = "en"
)
