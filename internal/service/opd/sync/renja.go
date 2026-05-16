package sync

import (
	"context"
	"log/slog"

	"github.com/bappeda-dev-team/penetapan-service/internal/client/perencanaan"
	"github.com/bappeda-dev-team/penetapan-service/internal/model/web"
	"github.com/bappeda-dev-team/penetapan-service/internal/repository"
)

type RenjaSyncExecutor struct {
	Repo              *repository.PenetapanOpdRepository
	PerencanaanClient *perencanaan.PerencanaanClient
	Logger            *slog.Logger
}

func NewRenjaSyncExecutor(
	repo *repository.PenetapanOpdRepository,
	perencanaanClient *perencanaan.PerencanaanClient,
	logger *slog.Logger,
) *RenjaSyncExecutor {
	return &RenjaSyncExecutor{
		Repo:              repo,
		PerencanaanClient: perencanaanClient,
		Logger:            logger,
	}
}

func (ex *RenjaSyncExecutor) Sync(
	ctx context.Context,
	syncId int64,
	req *web.SyncPenetapanOpdRequest,
	currentUser string,
) (web.SyncPenetapanOpdSummary, error) {

	// simpan snapshot
	var jumlahRenja int
	var jumlahIndikator int
	var jumlahTarget int

	return web.SyncPenetapanOpdSummary{
		Renja:     &jumlahRenja,
		Indikator: jumlahIndikator,
		Target:    jumlahTarget,
	}, nil
}
