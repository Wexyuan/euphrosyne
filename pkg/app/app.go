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

// defaultShutdownTimeout bounds the shutdown timeout of the whole application.
const defaultShutdownTimeout = 15 * time.Second

// App manages the application lifecycle.
type App struct {
	name    string          // application name
	servers []server.Server // application servers
}

// Option configures an application.
type Option func(*App)

func NewApp(opts ...Option) *App {
	a := &App{}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func WithName(name string) Option {
	return func(a *App) {
		a.name = name
	}
}

func WithServers(servers ...server.Server) Option {
	return func(a *App) {
		a.servers = servers
	}
}

// Run runs the application.
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

	// Stop every server exactly once.
	shutdown := sync.OnceFunc(func() {
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

	// Trigger the shutdown on any cancellation, including internal ones.
	go func() {
		<-runCtx.Done()
		shutdown()
	}()

	var wg sync.WaitGroup
	for _, srv := range a.servers {
		wg.Go(func() {
			if err := srv.Start(runCtx); err != nil {
				// Keep the earliest failure as the root cause.
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				shutdown()
			}
		})
	}

	// Wait for every server to exit.
	wg.Wait()

	// Wait for the shutdown as well so Run only returns after every server has
	// stopped and the earliest stop error has been observed.
	shutdown()

	return firstErr
}

// stopServers stops all servers.
func (a *App) stopServers() error {
	// Use a separate context so a cancelled caller cannot shorten the shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	var mu sync.Mutex
	var firstErr error

	var wg sync.WaitGroup
	for _, srv := range a.servers {
		wg.Go(func() {
			if err := srv.Stop(ctx); err != nil {
				// Keep the earliest failure.
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
			}
		})
	}

	wg.Wait()
	return firstErr
}
