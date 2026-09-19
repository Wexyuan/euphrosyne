package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database provides the data access methods.
type Database struct {
	*gorm.DB
}

// Options holds the database settings.
type Options struct {
	Driver          string // database driver (mysql/postgres/sqlite)
	DSN             string // data source name
	MaxOpenConns    int    // max open connections (default to 10)
	MaxIdleConns    int    // max idle connections (default to 2)
	ConnMaxLifetime int    // max connection lifetime in minutes (default to 30)
}

func New(opts Options) (*Database, error) {
	if opts.Driver == "" {
		return nil, fmt.Errorf("[database] driver is required")
	}
	if opts.DSN == "" {
		return nil, fmt.Errorf("[database] dsn is required")
	}

	dialector, err := opts.resolveDialector()
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		// Translate driver errors so unique violations become gorm.ErrDuplicatedKey.
		TranslateError: true,
	})
	if err != nil {
		return nil, fmt.Errorf("[database] open database error: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("[database] get underlying sql database error: %w", err)
	}

	maxOpen := opts.MaxOpenConns
	if maxOpen <= 0 {
		// Cap open connections to avoid exhausting the server under load.
		maxOpen = 10
	}
	maxIdle := opts.MaxIdleConns
	if maxIdle <= 0 {
		maxIdle = 2
	}
	lifetime := opts.ConnMaxLifetime
	if lifetime <= 0 {
		// Recycle connections periodically so stale ones get replaced.
		lifetime = 30
	}

	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(time.Duration(lifetime) * time.Minute)

	return &Database{DB: db}, nil
}

// resolveDialector resolves the dialector from the configured driver.
func (o Options) resolveDialector() (gorm.Dialector, error) {
	switch strings.ToLower(o.Driver) {
	case "mysql":
		return mysql.Open(o.DSN), nil
	case "postgres":
		return postgres.Open(o.DSN), nil
	case "sqlite":
		return sqlite.Open(o.DSN), nil
	default:
		return nil, fmt.Errorf("[database] invalid driver %q: must be mysql, postgres or sqlite", o.Driver)
	}
}

// Close closes the database connection.
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("[database] get underlying sql database error: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("[database] close database error: %w", err)
	}
	return nil
}
