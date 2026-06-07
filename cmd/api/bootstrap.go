package main

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/bappeda-dev-team/penetapan-service/internal/api"
	"github.com/bappeda-dev-team/penetapan-service/internal/client/perencanaan"
	"github.com/bappeda-dev-team/penetapan-service/internal/repository"
	"github.com/bappeda-dev-team/penetapan-service/internal/service/individu"
	"github.com/bappeda-dev-team/penetapan-service/internal/service/opd"
	"github.com/bappeda-dev-team/penetapan-service/internal/service/opd/sync"

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
	penetapanOpdRepo := repository.NewPenetapanOpdRepository(db)
	penetapanIndiRepo := repository.NewPenetapanIndividuRepository(db)

	// external service
	// perencanaan
	perencanaanClient := perencanaan.NewPerencanaanClient(
		cfg.Services.Perencanaan.BaseURL,
		cfg.Services.Perencanaan.ApiPath,
		&http.Client{Timeout: 60 * time.Second},
	)

	// executor
	//// tujuan
	tujuanSyncExecutor := sync.NewTujuanSyncExecutor(
		penetapanOpdRepo,
		perencanaanClient,
		logger,
	)

	//// sasaran
	sasaranSyncExecutor := sync.NewSasaranSyncExecutor(
		penetapanOpdRepo,
		perencanaanClient,
		logger,
	)

	//// renja
	renjaSyncExecutor := sync.NewRenjaSyncExecutor(
		penetapanOpdRepo,
		perencanaanClient,
		logger,
	)

	// register executor
	syncRegistry := &sync.Registry{
		TujuanSyncExecutor:  tujuanSyncExecutor,
		SasaranSyncExecutor: sasaranSyncExecutor,
		RenjaSyncExecutor:   renjaSyncExecutor,
	}

	// service
	penetapanOpdService := opd.NewPenetapanOpdService(
		penetapanOpdRepo,
		perencanaanClient,
		syncRegistry,
		logger,
	)

	individuService := individu.NewPenetapanIndividuService(
		penetapanIndiRepo,
		perencanaanClient,
		syncRegistry,
		logger,
	)

	// application
	app := &api.Application{
		Config: cfg,
		Logger: logger,

		PenetapanOpdService: penetapanOpdService,
		IndividuService:     individuService,
	}

	return app
}
