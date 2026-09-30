package database

import (
	"fmt"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// Database provides the data access methods.
type Database struct {
	*gorm.DB
}

// Options holds the database settings.
type Options struct {
	DSN             string // data source name
	MaxOpenConns    int    // max open connections (default to 10)
	MaxIdleConns    int    // max idle connections (default to 2)
	ConnMaxLifetime int    // max connection lifetime in minutes (default to 30)
}

func New(opts Options) (*Database, error) {
	if opts.DSN == "" {
		return nil, fmt.Errorf("[database] dsn is required")
	}

	db, err := gorm.Open(sqlite.Open(opts.DSN), &gorm.Config{
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
