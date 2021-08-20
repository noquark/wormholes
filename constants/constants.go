package constants

const (
	DefaultPort  = 3000
	DefaultConf  = "config.yml"
	DefaultDir   = "wormholes"
	DirPerm      = 0o775
	FilePerm     = 0o600
	EnvPrefix    = "WH"
	DotDir       = ".wormholes"
	EmptyString  = ""
	CacheControl = "private, max-age=90"
	CookieName   = "_wh"
	BloomDB      = "bloom.db"
	MaxLimit     = 1e7
	ErrorRate    = 1e-3
	MaxTry       = 10
	IDSize       = 7
	CookieSize   = 21
	TokenSize    = 43
	Streams      = 8
	BatchSize    = 1e4
	CityDB       = "GeoLite2-City.mmdb"
	EN           = "en"
)
