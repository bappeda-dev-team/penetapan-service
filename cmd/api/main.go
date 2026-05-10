package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/bappeda-dev-team/penetapan-service/internal/api"
	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()

	var cfg api.Config

	flag.IntVar(
		&cfg.Port,
		"port",
		4000,
		"API server port",
	)

	flag.StringVar(
		&cfg.Env,
		"env",
		"development",
		"Environment (production|staging|development)",
	)

	flag.Parse()

	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			nil,
		),
	)

	cfg.DB.Dsn = os.Getenv("DB_DSN")

	if cfg.DB.Dsn == "" {

		logger.Error(
			"DB_DSN environment variable is required",
		)

		os.Exit(1)
	}

	// database connection
	db, err := openDB(cfg.DB.Dsn)

	if err != nil {

		logger.Error(err.Error())

		os.Exit(1)
	}

	defer db.Close()

	// repository
	penetapanOpdRepo := &api.PenetapanOpdRepository{
		DB: db,
	}

	// service
	penetapanOpdService := &api.PenetapanOpdService{
		Repo: penetapanOpdRepo,
	}

	// application
	app := &api.Application{
		Config: cfg,
		Logger: logger,

		PenetapanOpdService: penetapanOpdService,
	}

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
