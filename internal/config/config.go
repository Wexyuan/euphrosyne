package config

import "time"

// Config holds the whole application configuration loaded from the yaml file.
type Config struct {
	// App is the application identity settings.
	App App `mapstructure:"app"`
	// Server is the HTTP server settings.
	Server Server `mapstructure:"server"`
	// Logger is the logger settings.
	Logger Logger `mapstructure:"logger"`
	// Database is the database settings.
	Database Database `mapstructure:"database"`
	// Cache is the cache settings.
	Cache Cache `mapstructure:"cache"`
	// JWT is the JSON web token settings.
	JWT JWT `mapstructure:"jwt"`
}

// App holds the application identity settings.
type App struct {
	// Name is the application name used in logs and identification.
	Name string `mapstructure:"name"`
	// Env is the runtime environment (debug/release/test, case-insensitive).
	Env string `mapstructure:"env"`
	// SnowflakeNode is the snowflake worker node number (0-1023).
	SnowflakeNode int64 `mapstructure:"snowflake_node"`
}

// Server holds the HTTP server settings.
type Server struct {
	// Addr is the listen address, ":8080" by default.
	Addr string `mapstructure:"addr"`
	// ReadTimeout bounds request reads, defaulting to 10s.
	ReadTimeout time.Duration `mapstructure:"read_timeout"`
	// WriteTimeout bounds response writes, defaulting to 10s.
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	// ShutdownTimeout bounds the shutdown, defaulting to 5s.
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

// Logger holds the logger settings.
type Logger struct {
	// Level is the minimum log level (debug/info/warn/error, case-insensitive).
	Level string `mapstructure:"level"`
	// Format is the log encoding format (console/json, case-insensitive).
	Format string `mapstructure:"format"`
	// Output is the output stream (stdout/stderr, case-insensitive).
	Output string `mapstructure:"output"`
}

// Database holds the database settings.
type Database struct {
	// Driver is the database driver (mysql/postgres/sqlite, case-insensitive).
	Driver string `mapstructure:"driver"`
	// DSN is the data source name of the target database.
	DSN string `mapstructure:"dsn"`
	// MaxOpenConns is the maximum number of open connections, defaulting to 10.
	MaxOpenConns int `mapstructure:"max_open_conns"`
	// MaxIdleConns is the maximum number of idle connections, defaulting to 2.
	MaxIdleConns int `mapstructure:"max_idle_conns"`
	// ConnMaxLifetime is the maximum lifetime of a connection in minutes, defaulting to 30.
	ConnMaxLifetime int `mapstructure:"conn_max_lifetime"`
}

// Cache holds the cache settings.
type Cache struct {
	// Driver is the cache driver (memory/redis, case-insensitive).
	Driver string `mapstructure:"driver"`
	// TTL is the default expiration in seconds, defaulting to 300.
	TTL int `mapstructure:"ttl"`
	// Addr is the redis server address in host:port form.
	Addr string `mapstructure:"addr"`
	// Password is the redis auth password, empty if auth is disabled.
	Password string `mapstructure:"password"`
	// DB is the redis logical database number, defaulting to 0.
	DB int `mapstructure:"db"`
}

// JWT holds the access and refresh token settings.
type JWT struct {
	// Secret is the HS256 signing key used to sign tokens.
	Secret string `mapstructure:"secret"`
	// Issuer is the issuer claim written into every token.
	Issuer string `mapstructure:"issuer"`
	// AccessTTL is the access token lifetime, defaulting to 15m.
	AccessTTL time.Duration `mapstructure:"access_ttl"`
	// RefreshTTL is the refresh token lifetime, defaulting to 168h.
	RefreshTTL time.Duration `mapstructure:"refresh_ttl"`
}
