package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// HTTPServer serves HTTP requests.
type HTTPServer struct {
	addr            string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration
	handler         http.Handler

	mu  sync.Mutex
	srv *http.Server
}

// HTTPOptions holds the HTTP server settings.
type HTTPOptions struct {
	// Addr is the listen address, ":8080" by default.
	Addr string
	// ReadTimeout bounds request reads, defaulting to 10s.
	ReadTimeout time.Duration
	// WriteTimeout bounds response writes, defaulting to 10s.
	WriteTimeout time.Duration
	// ShutdownTimeout bounds the shutdown, defaulting to 5s.
	ShutdownTimeout time.Duration
	// Handler is the root handler.
	Handler http.Handler
}

// NewHTTP creates an HTTP server from the given options.
func NewHTTP(opts HTTPOptions) *HTTPServer {
	if opts.Addr == "" {
		opts.Addr = ":8080"
	}
	if opts.ReadTimeout <= 0 {
		opts.ReadTimeout = 10 * time.Second
	}
	if opts.WriteTimeout <= 0 {
		opts.WriteTimeout = 10 * time.Second
	}
	if opts.ShutdownTimeout <= 0 {
		opts.ShutdownTimeout = 5 * time.Second
	}

	return &HTTPServer{
		addr:            opts.Addr,
		readTimeout:     opts.ReadTimeout,
		writeTimeout:    opts.WriteTimeout,
		shutdownTimeout: opts.ShutdownTimeout,
		handler:         opts.Handler,
	}
}

// Start runs the HTTP server and blocks until it stops.
func (s *HTTPServer) Start(ctx context.Context) error {
	handler := s.handler
	if handler == nil {
		handler = http.NotFoundHandler()
	}

	srv := &http.Server{
		Addr:         s.addr,
		Handler:      handler,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
	}

	s.mu.Lock()
	s.srv = srv
	s.mu.Unlock()

	errCh := make(chan error, 1)
	go func() {
		// ErrServerClosed signals a graceful shutdown, not a failure.
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("[server] listen server %s error: %w", s.addr, err)
			return
		}
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		// A serve failure stops the server as well.
		return err
	case <-ctx.Done():
		// Drain in-flight requests before returning.
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("[server] shutdown %s error: %w", s.addr, err)
		}
		return <-errCh
	}
}

// Stop shuts the HTTP server down gracefully.
func (s *HTTPServer) Stop(ctx context.Context) error {
	s.mu.Lock()
	srv := s.srv
	s.mu.Unlock()
	if srv == nil {
		return nil
	}
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("[server] shutdown %s error: %w", s.addr, err)
	}
	return nil
}
