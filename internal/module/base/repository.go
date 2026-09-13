package base

import (
	"context"

	"gorm.io/gorm"

	"github.com/Wexyuan/euphrosyne/pkg/cache"
	"github.com/Wexyuan/euphrosyne/pkg/database"
	"github.com/Wexyuan/euphrosyne/pkg/id"
)

// ctxKey defines context key type.
type ctxKey string

// txKey is the transaction context key.
const txKey ctxKey = "tx"

// Repository provides the shared data access for business repositories.
type Repository struct {
	db *database.Database
	sf *id.Snowflake
	c  cache.Cache
}

// NewRepository creates the base repository.
func NewRepository(db *database.Database, sf *id.Snowflake, c cache.Cache) *Repository {
	return &Repository{
		db: db,
		sf: sf,
		c:  c,
	}
}

// GenerateID returns a new unique ID.
func (r *Repository) GenerateID() int64 {
	return r.sf.NextID()
}

// DB returns the database handle bound to the context.
func (r *Repository) DB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok {
		return tx
	}
	return r.db.WithContext(ctx)
}

// Transaction runs the callback within a transaction.
func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, txKey, tx))
	})
}

// Cache returns the shared cache store.
func (r *Repository) Cache() cache.Cache {
	return r.c
}
