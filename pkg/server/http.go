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
	addr            string        // listen address
	readTimeout     time.Duration // max request read time
	writeTimeout    time.Duration // max response write time
	shutdownTimeout time.Duration // max shutdown time
	handler         http.Handler  // root HTTP handler

	mu  sync.Mutex   // guards srv
	srv *http.Server // running server
}

// HTTPOptions holds the HTTP server settings.
type HTTPOptions struct {
	Addr            string        // listen address (default to :8080)
	ReadTimeout     time.Duration // max request read time (default to 10s)
	WriteTimeout    time.Duration // max response write time (default to 10s)
	ShutdownTimeout time.Duration // max shutdown time (default to 5s)
	Handler         http.Handler  // root HTTP handler
}

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

// Start starts the HTTP server.
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
		// Treat ErrServerClosed as a graceful shutdown, not a failure.
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("[server] listen server %s error: %w", s.addr, err)
			return
		}
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		// Drain in-flight requests before returning.
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("[server] shutdown server %s error: %w", s.addr, err)
		}
		return <-errCh
	}
}

// Stop stops the HTTP server.
func (s *HTTPServer) Stop(ctx context.Context) error {
	s.mu.Lock()
	srv := s.srv
	s.mu.Unlock()
	if srv == nil {
		return nil
	}
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("[server] shutdown server %s error: %w", s.addr, err)
	}
	return nil
}
