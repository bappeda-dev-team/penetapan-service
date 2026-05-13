package main

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/bappeda-dev-team/penetapan-service/internal/api"
	"github.com/bappeda-dev-team/penetapan-service/internal/client/perencanaan"

	"net/http"
)

// buildApplication adalah fungsi untuk injeksi repo dan service
// tambahkan service dan repository baru kesini
func buildApplication(
	cfg api.Config,
	logger *slog.Logger,
	db *sql.DB,
) *api.Application {

	// repository
	penetapanOpdRepo := &api.PenetapanOpdRepository{
		DB: db,
	}

	// external service
	// perencanaan
	perencanaan := perencanaan.NewPerencanaanClient(
		cfg.Services.Perencanaan.BaseURL,
		cfg.Services.Perencanaan.ApiPath,
		&http.Client{Timeout: 60 * time.Second},
	)

	// service
	penetapanOpdService := &api.PenetapanOpdService{
		Repo:        penetapanOpdRepo,
		Perencanaan: perencanaan,
		Logger:      logger,
	}

	// application
	app := &api.Application{
		Config: cfg,
		Logger: logger,

		PenetapanOpdService: penetapanOpdService,
	}

	return app
}
