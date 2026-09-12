package server

import "context"

// Server defines the lifecycle shared by all services.
type Server interface {
	// Start runs the service and blocks until it stops.
	Start(ctx context.Context) error
	// Stop shuts the service down gracefully.
	Stop(ctx context.Context) error
}
