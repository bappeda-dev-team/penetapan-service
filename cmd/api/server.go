package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bappeda-dev-team/penetapan-service/internal/api"
)

func serve(app *api.Application, cfg api.Config, logger *slog.Logger) error {

	srv := &http.Server{
		Addr: fmt.Sprintf(
			":%d",
			cfg.Port,
		),

		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog: slog.NewLogLogger(
			logger.Handler(),
			slog.LevelError,
		),
	}

	shutdownError := make(chan error)

	// Start a background goroutine.
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		logger.Info("shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.Env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdownError
	if err != nil {
		return err
	}

	logger.Info("stopped server", "addr", srv.Addr)

	return nil
}
