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
	// Driver is the database driver (mysql/postgres/sqlite, case-insensitive).
	Driver string
	// DSN is the data source name of the target database.
	DSN string
	// MaxOpenConns is the maximum number of open connections, defaulting to 10.
	MaxOpenConns int
	// MaxIdleConns is the maximum number of idle connections, defaulting to 2.
	MaxIdleConns int
	// ConnMaxLifetime is the maximum lifetime of a connection in minutes, defaulting to 30.
	ConnMaxLifetime int
}

// New creates the database from the given options.
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

	if err := opts.applyPool(db); err != nil {
		return nil, err
	}

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

// applyPool applies the connection pool settings with sane defaults.
func (o Options) applyPool(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("[database] get underlying sql database error: %w", err)
	}

	maxOpen := o.MaxOpenConns
	if maxOpen <= 0 {
		// Cap open connections to avoid exhausting the server under load.
		maxOpen = 10
	}
	maxIdle := o.MaxIdleConns
	if maxIdle <= 0 {
		maxIdle = 2
	}
	lifetime := o.ConnMaxLifetime
	if lifetime <= 0 {
		// Recycle connections periodically so stale ones get replaced.
		lifetime = 30
	}

	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(time.Duration(lifetime) * time.Minute)
	return nil
}
