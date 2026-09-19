package base

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/Wexyuan/euphrosyne/internal/constant"
	"github.com/Wexyuan/euphrosyne/pkg/cache"
	"github.com/Wexyuan/euphrosyne/pkg/database"
	"github.com/Wexyuan/euphrosyne/pkg/id"
)

// Repository provides the shared data access.
type Repository struct {
	db *database.Database // underlying database
	sf *id.Snowflake      // snowflake ID generator
	c  cache.Cache        // shared cache
}

func NewRepository(db *database.Database, sf *id.Snowflake, c cache.Cache) *Repository {
	return &Repository{
		db: db,
		sf: sf,
		c:  c,
	}
}

// GenerateID generates a unique ID.
func (r *Repository) GenerateID() int64 {
	return r.sf.NextID()
}

// DB returns the database handle.
func (r *Repository) DB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(constant.TxKey).(*gorm.DB); ok {
		return tx
	}
	return r.db.WithContext(ctx)
}

// Transaction runs the callback in a transaction.
func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, constant.TxKey, tx))
	})
	if err != nil {
		return fmt.Errorf("[base] run transaction error: %w", err)
	}
	return nil
}

// Cache returns the shared cache.
func (r *Repository) Cache() cache.Cache {
	return r.c
}
