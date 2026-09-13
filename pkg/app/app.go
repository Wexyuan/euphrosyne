package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Wexyuan/euphrosyne/pkg/server"
)

// App manages the lifecycle of application.
type App struct {
	name    string
	servers []server.Server
}

// NewApp creates the application with the given options.
func NewApp(opts ...Option) *App {
	a := &App{}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// Option configures an application.
type Option func(*App)

// WithName sets the application name.
func WithName(name string) Option {
	return func(a *App) {
		a.name = name
	}
}

// WithServers sets the application servers.
func WithServers(servers ...server.Server) Option {
	return func(a *App) {
		a.servers = servers
	}
}

// Run starts the application and blocks until it is stopped.
func (a *App) Run(ctx context.Context) error {
	if len(a.servers) == 0 {
		return fmt.Errorf("[app] no server is configured")
	}

	// Turn termination signals into context cancellation.
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var mu sync.Mutex
	var firstErr error

	// shutdown stops every server exactly once.
	var once sync.Once
	shutdown := func() {
		once.Do(func() {
			// Cancel first so the servers wind down.
			cancel()
			if err := a.stopServers(); err != nil {
				// Keep the earliest failure as the root cause.
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
			}
		})
	}

	// Trigger the shutdown on any cancellation, including internal ones.
	go func() {
		<-runCtx.Done()
		shutdown()
	}()

	var wg sync.WaitGroup
	for _, srv := range a.servers {
		wg.Add(1)
		go func(srv server.Server) {
			defer wg.Done()
			if err := srv.Start(runCtx); err != nil {
				// Keep the earliest failure as the root cause.
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				shutdown()
			}
		}(srv)
	}

	// Wait for every server to exit.
	wg.Wait()
	return firstErr
}

// stopServers stops application servers.
func (a *App) stopServers() error {
	// Use a separate context so a cancelled caller cannot shorten the shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var mu sync.Mutex
	var firstErr error

	var wg sync.WaitGroup
	for _, srv := range a.servers {
		wg.Add(1)
		go func(srv server.Server) {
			defer wg.Done()
			if err := srv.Stop(ctx); err != nil {
				// Keep the earliest failure.
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
			}
		}(srv)
	}

	wg.Wait()
	return firstErr
}
