package opd

import (
	"context"
	"log/slog"
	"time"

	"github.com/bappeda-dev-team/penetapan-service/internal/client/perencanaan"
	"github.com/bappeda-dev-team/penetapan-service/internal/model/domain"
	"github.com/bappeda-dev-team/penetapan-service/internal/model/web"
	"github.com/bappeda-dev-team/penetapan-service/internal/repository"
	"github.com/bappeda-dev-team/penetapan-service/internal/service/opd/sync"
)

type PenetapanOpdService struct {
	Repo         *repository.PenetapanOpdRepository
	Perencanaan  *perencanaan.PerencanaanClient
	SyncExecutor *sync.Registry
	// logger taruh paling bawah
	Logger *slog.Logger
}

func NewPenetapanOpdService(
	repo *repository.PenetapanOpdRepository,
	perencanaanClient *perencanaan.PerencanaanClient,
	syncExecutor *sync.Registry,
	// logger taruh paling bawah
	logger *slog.Logger,
) *PenetapanOpdService {
	return &PenetapanOpdService{
		Repo:         repo,
		Perencanaan:  perencanaanClient,
		SyncExecutor: syncExecutor,
		Logger:       logger,
	}
}

func (s *PenetapanOpdService) SyncPenetapanOpd(
	ctx context.Context,
	req *web.SyncPenetapanOpdRequest,
) (web.SyncPenetapanOpdResponse, error) {
	now := time.Now()
	currentUser := "super_admin" // ambil dari ctx nanti

	jenisPenetapan := req.JenisPenetapan

	executor, err := s.SyncExecutor.Get(jenisPenetapan)
	if err != nil {
		return web.SyncPenetapanOpdResponse{}, err
	}

	// insert metadata pertama (pending)
	metadata := domain.SyncPenetapanMetadataOpd{
		KodeOpd:        req.KodeOpd,
		Tahun:          req.Tahun,
		Status:         domain.SyncStatusPending,
		JenisPenetapan: jenisPenetapan,
		StartedAt:      now,
		SyncBy:         &currentUser,
	}

	// log start
	s.Logger.InfoContext(
		ctx,
		"sync penetapan started",
		"kode_opd", metadata.KodeOpd,
		"tahun", metadata.Tahun,
		"jenis_sync", metadata.JenisPenetapan,
		"sync_by", metadata.SyncBy,
		"started_at", metadata.StartedAt,
	)

	syncId, err := s.Repo.InsertMetadata(ctx, metadata)
	if err != nil {
		return web.SyncPenetapanOpdResponse{}, err
	}

	// updateMetadata jadi InProgress
	err = s.markSyncAsInProgress(ctx, syncId)
	if err != nil {
		return web.SyncPenetapanOpdResponse{}, err
	}

	summary, err := executor.Sync(ctx, syncId, req, currentUser)
	if err != nil {
		return web.SyncPenetapanOpdResponse{}, s.failSync(ctx, syncId, err)
	}

	// jika berhasil updateMetadata jadi Success
	err = s.markSyncAsSuccess(ctx, syncId)
	if err != nil {
		return web.SyncPenetapanOpdResponse{}, err
	}

	return web.SyncPenetapanOpdResponse{
		SyncId:           syncId,
		Status:           domain.SyncStatusSuccess,
		KodeOpd:          req.KodeOpd,
		Tahun:            req.Tahun,
		JenisPenetapan:   jenisPenetapan,
		ProcessedAt:      time.Now(),
		ProcessedSummary: summary,
	}, nil

}

func (s *PenetapanOpdService) FindTujuan(
	ctx context.Context,
	req domain.PenetapanOpdRequest,
) ([]web.TujuanPenetapanOpdResponse, error) {
	jenisPenetapan := domain.JenisPenetapanTujuan
	snapshotActive, err := s.Repo.GetActiveSnapshot(ctx, req.KodeOpd, jenisPenetapan, req.Tahun)
	if err != nil {
		s.Logger.Error("GetActiveSnapshot")
		return nil, err
	}
	if snapshotActive == nil {
		return []web.TujuanPenetapanOpdResponse{}, nil
	}
	req.SnapshotId = snapshotActive

	tujuanOpd, err := s.Repo.FindTujuanBySnapshot(ctx, req)
	if err != nil {
		s.Logger.Error("FindTujuanBySnapshot")
		return nil, err
	}
	tujuanOpdIds := make([]int64, 0, len(tujuanOpd))
	for _, tj := range tujuanOpd {
		tujuanOpdIds = append(tujuanOpdIds, tj.Id)
	}
	indikators, err := s.Repo.FindIndikatorTujuanByTujuanIds(ctx, tujuanOpdIds)
	if err != nil {
		s.Logger.Error("FindIndikatorTujuanByTujuanIds")
		return nil, err
	}
	indikatorIds := make([]int64, 0, len(indikators))
	for _, ind := range indikators {
		indikatorIds = append(indikatorIds, ind.Id)
	}
	targets, err := s.Repo.FindTargetIndikatorTujuanByIndikatorIds(ctx, indikatorIds)
	if err != nil {
		s.Logger.Error("FindTargetIndikatorTujuanByIndikatorIds")
		return nil, err
	}
	targetMap := make(
		map[int64][]web.TargetIndikatorResponse,
	)
	for _, target := range targets {
		targetResp := ToTargetIndikatorTujuanOpdResponse(target)
		targetMap[target.IndikatorTujuanId] = append(
			targetMap[target.IndikatorTujuanId],
			targetResp,
		)
	}

	indikatorMap := make(
		map[int64][]web.IndikatorTujuanPenetapanResponse,
	)
	for _, indikator := range indikators {
		indikatorResp := ToIndikatorTujuanOpdResponse(indikator)
		indikatorResp.Target = targetMap[indikator.Id]
		indikatorMap[indikator.IdTujuanOpd] = append(
			indikatorMap[indikator.IdTujuanOpd],
			indikatorResp,
		)
	}

	result := make(
		[]web.TujuanPenetapanOpdResponse,
		0,
		len(tujuanOpd),
	)
	for _, tujuan := range tujuanOpd {
		tujResp := ToTujuanOpdResponse(tujuan)
		tujResp.Indikator = indikatorMap[tujuan.Id]
		result = append(result, tujResp)
	}

	return result, nil
}

func (s *PenetapanOpdService) FindSasaran(
	ctx context.Context,
	req domain.PenetapanOpdRequest,
) ([]web.SasaranPenetapanOpdResponse, error) {
	sasaranOpd, err := s.Repo.FindSasaran(ctx, req)
	if err != nil {
		return nil, err
	}

	result := make(
		[]web.SasaranPenetapanOpdResponse,
		0,
		len(sasaranOpd),
	)
	for _, sasaran := range sasaranOpd {
		result = append(result, ToSasaranOpdResponse(sasaran))
	}

	return result, nil
}

func (s *PenetapanOpdService) markSyncAsFailed(ctx context.Context, syncId int64, msg string) error {
	statusSync := domain.SyncStatusFailed
	errorMessage := &msg
	return s.Repo.UpdateMetadataStatus(
		ctx,
		syncId,
		statusSync,
		errorMessage,
	)
}

func (s *PenetapanOpdService) markSyncAsInProgress(ctx context.Context, syncId int64) error {
	statusSync := domain.SyncStatusInProgress
	return s.Repo.UpdateMetadataStatus(
		ctx,
		syncId,
		statusSync,
		nil,
	)
}

func (s *PenetapanOpdService) markSyncAsSuccess(ctx context.Context, syncId int64) error {
	statusSync := domain.SyncStatusSuccess
	return s.Repo.UpdateMetadataStatus(
		ctx,
		syncId,
		statusSync,
		nil,
	)
}

func (s *PenetapanOpdService) failSync(ctx context.Context, syncId int64, cause error) error {
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
