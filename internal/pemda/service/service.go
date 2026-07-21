package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/bappeda-dev-team/penetapan-service/internal/pemda/client"
	"github.com/bappeda-dev-team/penetapan-service/internal/pemda/domain"
	"github.com/bappeda-dev-team/penetapan-service/internal/pemda/repository"
	"github.com/bappeda-dev-team/penetapan-service/internal/pemda/service/sync"
	"github.com/bappeda-dev-team/penetapan-service/internal/pemda/web"
)

type PenetapanPemdaService struct {
	Repo         *repository.PenetapanPemdaRepository
	Client       *client.Client
	SyncExecutor *sync.Registry
	Logger       *slog.Logger
}

func NewPenetapanPemdaService(
	repo *repository.PenetapanPemdaRepository,
	client *client.Client,
	syncExecutor *sync.Registry,
	logger *slog.Logger,
) *PenetapanPemdaService {
	return &PenetapanPemdaService{
		Repo:         repo,
		Client:       client,
		SyncExecutor: syncExecutor,
		Logger:       logger,
	}
}

func (s *PenetapanPemdaService) SyncPenetapanTujuanPemda(
	ctx context.Context,
	req *web.SyncPenetapanRequest,
) (web.SyncPenetapanResponse, error) {
	now := time.Now()
	currentUser := "super_admin" // ambil dari ctx nanti
	jenisPenetapan := domain.JenisPenetapanTujuanPemda

	executor, err := s.SyncExecutor.Get(jenisPenetapan)
	if err != nil {
		s.Logger.Error("ERROR GET EXECUTOR", "err", err)
		return web.SyncPenetapanResponse{}, err
	}

	// insert metadata pertama
	metadata := domain.SyncPenetapanMetadata{
		Tahun:          req.Tahun,
		Status:         domain.SyncStatusPending,
		StartedAt:      now,
		SyncBy:         &currentUser,
		JenisPenetapan: jenisPenetapan,
	}

	// log start
	s.Logger.InfoContext(
		ctx,
		"sync penetapan pk individu started",
		"tahun", metadata.Tahun,
		"jenis_sync", metadata.JenisPenetapan,
		"sync_by", metadata.SyncBy,
		"started_at", metadata.StartedAt,
	)

	syncId, err := s.Repo.InsertMetadata(ctx, metadata)
	if err != nil {
		s.Logger.Error("ERROR INSERT METADATA", "err", err)
		return web.SyncPenetapanResponse{}, err
	}

	// updateMetadata jadi InProgress
	err = s.markSyncAsInProgress(ctx, syncId)
	if err != nil {
		s.Logger.Error("ERROR INSERT METADATA SYNC IN PROGRESS", "err", err)
		return web.SyncPenetapanResponse{}, err
	}

	// SYNC AND SAVE
	summary, err := executor.Sync(ctx, syncId, req, currentUser)
	if err != nil {
		s.Logger.Error("ERROR SYNC DATA", "err", err)
		return web.SyncPenetapanResponse{}, s.failSync(ctx, syncId, err)
	}

	// jika berhasil updateMetadata jadi Success
	err = s.markSyncAsSuccess(ctx, syncId)
	if err != nil {
		s.Logger.Error("ERROR INSERT METADATA SYNC SUCCESS", "err", err)
		return web.SyncPenetapanResponse{}, err
	}

	return web.SyncPenetapanResponse{
		SyncId:           syncId,
		Status:           domain.SyncStatusSuccess,
		Tahun:            req.Tahun,
		JenisPenetapan:   jenisPenetapan,
		ProcessedAt:      time.Now(),
		ProcessedSummary: summary,
	}, nil
}

func (s *PenetapanPemdaService) markSyncAsFailed(ctx context.Context, syncId int64, msg string) error {
	statusSync := domain.SyncStatusFailed
	errorMessage := &msg
	return s.Repo.UpdateMetadataStatus(
		ctx,
		syncId,
		statusSync,
		errorMessage,
	)
}

func (s *PenetapanPemdaService) markSyncAsInProgress(ctx context.Context, syncId int64) error {
	statusSync := domain.SyncStatusInProgress
	return s.Repo.UpdateMetadataStatus(
		ctx,
		syncId,
		statusSync,
		nil,
	)
}

func (s *PenetapanPemdaService) markSyncAsSuccess(ctx context.Context, syncId int64) error {
	statusSync := domain.SyncStatusSuccess
	return s.Repo.UpdateMetadataStatus(
		ctx,
		syncId,
		statusSync,
		nil,
	)
}

func (s *PenetapanPemdaService) failSync(ctx context.Context, syncId int64, cause error) error {
	updateErr := s.markSyncAsFailed(
		ctx,
		syncId,
		cause.Error(),
	)
	if updateErr != nil {
		return updateErr
	}
	return cause
}
