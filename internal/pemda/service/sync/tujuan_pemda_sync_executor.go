package sync

import (
	"context"
	"log/slog"

	"github.com/bappeda-dev-team/penetapan-service/internal/pemda/client"
	"github.com/bappeda-dev-team/penetapan-service/internal/pemda/repository"
	"github.com/bappeda-dev-team/penetapan-service/internal/pemda/web"
)

type TujuanSyncExecutor struct {
	Repo   *repository.PenetapanPemdaRepository
	Client *client.Client
	Logger *slog.Logger
}

func NewTujuanSyncExecutor(
	repo *repository.PenetapanPemdaRepository,
	client *client.Client,
	logger *slog.Logger,
) *TujuanSyncExecutor {
	return &TujuanSyncExecutor{
		Repo:   repo,
		Client: client,
		Logger: logger,
	}
}

func (ex *TujuanSyncExecutor) Sync(
	ctx context.Context,
	syncId int64,
	req *web.SyncPenetapanRequest,
	currentUser string,
) (web.SyncPenetapanSummary, error) {
	request := client.SyncRequest{
		Tahun: req.Tahun,
	}
	// outsource
	tujuanPemdas, err := ex.Client.SyncTujuanPemdaPenetapan(ctx, request)
	if err != nil {
		return web.SyncPenetapanSummary{}, err
	}
	ex.Logger.Info("TUJUAN PEMDA SYNC EXECUTOR", "tujuanPemdas", tujuanPemdas)

	response := web.SyncPenetapanSummary{}
	return response, nil
}
