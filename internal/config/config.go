package config

import "time"

// Config holds the application configuration.
type Config struct {
	App      App      `mapstructure:"app"`      // application identity settings
	Server   Server   `mapstructure:"server"`   // HTTP server settings
	Logger   Logger   `mapstructure:"logger"`   // logger settings
	Database Database `mapstructure:"database"` // database settings
	Cache    Cache    `mapstructure:"cache"`    // cache settings
	JWT      JWT      `mapstructure:"jwt"`      // JSON web token settings
}

// App holds the application identity settings.
type App struct {
	Name          string `mapstructure:"name"`           // application name
	Env           string `mapstructure:"env"`            // runtime environment (debug/release/test)
	SnowflakeNode int64  `mapstructure:"snowflake_node"` // snowflake worker node (0-1023)
}

// Server holds the HTTP server settings.
type Server struct {
	Addr            string        `mapstructure:"addr"`             // listen address (default to :8080)
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`     // max request read time (default to 10s)
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`    // max response write time (default to 10s)
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"` // max shutdown time (default to 5s)
}

// Logger holds the logger settings.
type Logger struct {
	Level  string `mapstructure:"level"`  // minimum log level (debug/info/warn/error)
	Format string `mapstructure:"format"` // log encoding format (console/json)
	Output string `mapstructure:"output"` // output stream (stdout/stderr)
}

// Database holds the database settings.
type Database struct {
	Driver          string `mapstructure:"driver"`            // database driver (mysql/postgres/sqlite)
	DSN             string `mapstructure:"dsn"`               // data source name
	MaxOpenConns    int    `mapstructure:"max_open_conns"`    // max open connections (default to 10)
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`    // max idle connections (default to 2)
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"` // max connection lifetime in minutes (default to 30)
}

// Cache holds the cache settings.
type Cache struct {
	Driver   string `mapstructure:"driver"`   // cache driver (memory/redis)
	TTL      int    `mapstructure:"ttl"`      // expiration in seconds (default to 300)
	Addr     string `mapstructure:"addr"`     // redis address (host:port)
	Password string `mapstructure:"password"` // redis password
	DB       int    `mapstructure:"db"`       // redis database (default to 0)
}

// JWT holds the access and refresh token settings.
type JWT struct {
	Secret     string        `mapstructure:"secret"`      // HS256 signing key
	Issuer     string        `mapstructure:"issuer"`      // issuer claim
	AccessTTL  time.Duration `mapstructure:"access_ttl"`  // access token lifetime (default to 15m)
	RefreshTTL time.Duration `mapstructure:"refresh_ttl"` // refresh token lifetime (default to 168h)
}
