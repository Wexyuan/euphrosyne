package wire

import (
	"fmt"

	"github.com/Wexyuan/kairos/internal/config"
	"github.com/Wexyuan/kairos/internal/module/auth"
	"github.com/Wexyuan/kairos/pkg/cache"
	"github.com/Wexyuan/kairos/pkg/database"
	"github.com/Wexyuan/kairos/pkg/id"
	"github.com/Wexyuan/kairos/pkg/logger"
)

func provideDatabase(cfg *config.Config, log *logger.Logger) (*database.Database, func(), error) {
	db, err := database.New(database.Options{
		Driver:          cfg.Database.Driver,
		DSN:             cfg.Database.DSN,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("[wire] init database error: driver=%s: %w", cfg.Database.Driver, err)
	}
	return db, func() {
		if err := db.Close(); err != nil {
			log.Errorf("close database error: %v", err)
		}
	}, nil
}

func provideSnowflake(cfg *config.Config) (*id.Snowflake, error) {
	return id.New(cfg.App.SnowflakeNode)
}

func provideCache(cfg *config.Config, log *logger.Logger) (cache.Cache, func(), error) {
	store, err := cache.New(cache.Options{
		Driver:   cfg.Cache.Driver,
		Addr:     cfg.Cache.Addr,
		Password: cfg.Cache.Password,
		DB:       cfg.Cache.DB,
		TTL:      cfg.Cache.TTL,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("[wire] init cache error: driver=%s: %w", cfg.Cache.Driver, err)
	}
	return store, func() {
		if err := store.Close(); err != nil {
			log.Errorf("close cache error: %v", err)
		}
	}, nil
}

func provideTokenManager(cfg *config.Config) (*auth.TokenManager, error) {
	return auth.NewTokenManager(
		cfg.JWT.Secret,
		cfg.JWT.Issuer,
		cfg.JWT.AccessTTL,
		cfg.JWT.RefreshTTL,
	)
}
