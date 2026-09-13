package base

import (
	"context"

	"github.com/Wexyuan/euphrosyne/pkg/logger"
)

// Service provides the shared helpers for business services.
type Service struct {
	log  *logger.Logger
	repo *Repository
}

// NewService creates the base service.
func NewService(log *logger.Logger, repo *Repository) *Service {
	return &Service{
		log:  log,
		repo: repo,
	}
}

// Logger returns the shared logger.
func (s *Service) Logger() *logger.Logger {
	return s.log
}

// Transaction runs the callback within a transaction.
func (s *Service) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return s.repo.Transaction(ctx, fn)
}
