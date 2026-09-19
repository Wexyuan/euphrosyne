package server

import "context"

// Server defines the lifecycle of a long-running server.
type Server interface {
	// Start starts the server.
	Start(ctx context.Context) error
	// Stop stops the server.
	Stop(ctx context.Context) error
}
