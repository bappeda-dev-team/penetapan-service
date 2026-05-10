package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {

	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			nil,
		),
	)

	cfg := loadConfig(logger)

	// database
	db, err := openDB(cfg.DB.Dsn)

	if err != nil {

		logger.Error(err.Error())

		os.Exit(1)
	}

	defer db.Close()

	// application
	app := buildApplication(
		cfg,
		logger,
		db,
	)

	srv := &http.Server{
		Addr: fmt.Sprintf(
			":%d",
			cfg.Port,
		),

		Handler: app.Routes(),

		IdleTimeout: time.Minute,

		ReadTimeout: 5 * time.Second,

		WriteTimeout: 10 * time.Second,

		ErrorLog: slog.NewLogLogger(
			logger.Handler(),
			slog.LevelError,
		),
	}

	logger.Info(
		"starting server",
		"addr",
		srv.Addr,
		"env",
		cfg.Env,
	)

	err = srv.ListenAndServe()

	if err != nil {

		logger.Error(err.Error())

		os.Exit(1)
	}
}
