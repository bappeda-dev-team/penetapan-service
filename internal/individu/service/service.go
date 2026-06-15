package service

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/bappeda-dev-team/penetapan-service/internal/individu/client"
	"github.com/bappeda-dev-team/penetapan-service/internal/individu/domain"
	"github.com/bappeda-dev-team/penetapan-service/internal/individu/repository"
	"github.com/bappeda-dev-team/penetapan-service/internal/individu/service/sync"
	"github.com/bappeda-dev-team/penetapan-service/internal/individu/web"
)

type PenetapanIndividuService struct {
	Repo         *repository.PenetapanIndividuRepository
	Client       *client.Client
	SyncExecutor *sync.Registry
	Logger       *slog.Logger
}

func NewPenetapanIndividuService(
	repo *repository.PenetapanIndividuRepository,
	client *client.Client,
	syncExecutor *sync.Registry,
	// logger taruh paling bawah
	logger *slog.Logger,
) *PenetapanIndividuService {
	return &PenetapanIndividuService{
		Repo:         repo,
		Client:       client,
		SyncExecutor: syncExecutor,
		Logger:       logger,
	}
}

func (s *PenetapanIndividuService) FindRekinsIndividu(
	ctx context.Context,
	req web.SyncPenetapanRequest,
) (web.RekinPenetapanIndividuResponse, error) {
	return web.RekinPenetapanIndividuResponse{}, nil
}

func (s *PenetapanIndividuService) SyncPenetapanPkIndividu(
	ctx context.Context,
	req *web.SyncPenetapanRequest,
) (web.SyncPenetapanResponse, error) {
	now := time.Now()
	currentUser := "super_admin" // ambil dari ctx nanti
	jenisPenetapan := domain.JenisPenetapanPk

	executor, err := s.SyncExecutor.Get(jenisPenetapan)
	if err != nil {
		return web.SyncPenetapanResponse{}, err
	}

	// insert metadata pertama
	metadata := domain.SyncPenetapanMetadata{
		PegawaiId:      req.PegawaiId,
		KodeOpd:        req.KodeOpd,
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
		"pegawai_id", metadata.PegawaiId,
		"kode_opd", metadata.KodeOpd,
		"tahun", metadata.Tahun,
		"jenis_sync", metadata.JenisPenetapan,
		"sync_by", metadata.SyncBy,
		"started_at", metadata.StartedAt,
	)

	syncId, err := s.Repo.InsertMetadata(ctx, metadata)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return web.SyncPenetapanResponse{}, err
	}

	// updateMetadata jadi InProgress
	err = s.markSyncAsInProgress(ctx, syncId)
	if err != nil {
		return web.SyncPenetapanResponse{}, err
	}

	summary, err := executor.Sync(ctx, syncId, req, currentUser)
	if err != nil {
		return web.SyncPenetapanResponse{}, s.failSync(ctx, syncId, err)
	}

	// jika berhasil updateMetadata jadi Success
	err = s.markSyncAsSuccess(ctx, syncId)
	if err != nil {
		return web.SyncPenetapanResponse{}, err
	}

	return web.SyncPenetapanResponse{
		SyncId:           syncId,
		Status:           domain.SyncStatusSuccess,
		PegawaiId:        req.PegawaiId,
		KodeOpd:          req.KodeOpd,
		Tahun:            req.Tahun,
		JenisPenetapan:   jenisPenetapan,
		ProcessedAt:      time.Now(),
		ProcessedSummary: summary,
	}, nil
}

func (s *PenetapanIndividuService) markSyncAsFailed(ctx context.Context, syncId int64, msg string) error {
	statusSync := domain.SyncStatusFailed
	errorMessage := &msg
	return s.Repo.UpdateMetadataStatus(
		ctx,
		syncId,
		statusSync,
		errorMessage,
	)
}

func (s *PenetapanIndividuService) markSyncAsInProgress(ctx context.Context, syncId int64) error {
	statusSync := domain.SyncStatusInProgress
	return s.Repo.UpdateMetadataStatus(
		ctx,
		syncId,
		statusSync,
		nil,
	)
}

func (s *PenetapanIndividuService) markSyncAsSuccess(ctx context.Context, syncId int64) error {
	statusSync := domain.SyncStatusSuccess
	return s.Repo.UpdateMetadataStatus(
		ctx,
		syncId,
		statusSync,
		nil,
	)
}

func (s *PenetapanIndividuService) failSync(ctx context.Context, syncId int64, cause error) error {
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
