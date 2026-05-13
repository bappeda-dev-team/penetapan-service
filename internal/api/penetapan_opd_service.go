package api

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/bappeda-dev-team/penetapan-service/internal/client/perencanaan"
	"github.com/bappeda-dev-team/penetapan-service/internal/model/domain"
	"github.com/bappeda-dev-team/penetapan-service/internal/model/web"
)

type PenetapanOpdService struct {
	Repo        *PenetapanOpdRepository
	Perencanaan *perencanaan.PerencanaanClient
	Logger      *slog.Logger
}

func (s *PenetapanOpdService) SyncPenetapanOpd(
	ctx context.Context,
	req *web.SyncPenetapanOpdRequest,
) (web.SyncPenetapanOpdResponse, error) {
	now := time.Now()
	currentUser := "super_admin" // ambil dari ctx nanti
	statusSync := domain.SyncStatusPending

	var jenisPenetapan string
	switch req.JenisPenetapan {
	case "tujuan":
		jenisPenetapan = domain.JenisPenetapanTujuan

	case "sasaran":
		jenisPenetapan = domain.JenisPenetapanSasaran
	default:
		jenisPenetapan = "unknown"
	}

	metadata := domain.SyncPenetapanMetadataOpd{
		KodeOpd:        req.KodeOpd,
		Tahun:          req.Tahun,
		Status:         statusSync,
		JenisPenetapan: jenisPenetapan,
		StartedAt:      now,
		SyncBy:         &currentUser,
	}

	// TODO: add more log
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

	// insert metadata pertama (pending)
	syncId, err := s.Repo.InsertMetadata(ctx, metadata)
	if err != nil {
		return web.SyncPenetapanOpdResponse{}, err
	}

	perencanaanRequest := perencanaan.PerencanaanRequest{
		KodeOpd: req.KodeOpd,
		Tahun:   req.Tahun,
	}
	tujuanPerencanaan, err := s.Perencanaan.GetPenetapanTujuanOpd(ctx, perencanaanRequest)
	if err != nil {
		errUpdate := s.markSyncAsFailed(ctx, syncId, err.Error())
		if errUpdate != nil {
			return web.SyncPenetapanOpdResponse{}, errUpdate
		}
		return web.SyncPenetapanOpdResponse{}, err
	}

	// updateMetadata jadi InProgress
	err = s.markSyncAsInProgress(ctx, syncId)
	if err != nil {
		return web.SyncPenetapanOpdResponse{}, err
	}

	// mulai butuh tx
	tx, err := s.Repo.DB.BeginTx(ctx, nil)
	if err != nil {
		return web.SyncPenetapanOpdResponse{}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	versiPenetapan, err := s.Repo.GetPenetapanNextVersion(ctx, tx, req.KodeOpd, domain.JenisPenetapanTujuan, req.Tahun)
	if err != nil {
		return web.SyncPenetapanOpdResponse{}, err
	}
	if versiPenetapan > 1 {
		errDeact := s.Repo.DeactivateOldSnapshot(ctx, tx, req.KodeOpd, domain.JenisPenetapanTujuan, req.Tahun)
		if errDeact != nil {
			return web.SyncPenetapanOpdResponse{}, errDeact
		}
	}

	// simpan snapshot penetapan
	penetapanTujuan := domain.PenetapanOpd{
		KodeOpd:        req.KodeOpd,
		Tahun:          req.Tahun,
		JenisPenetapan: domain.JenisPenetapanTujuan,
		Versi:          versiPenetapan,
		SnapshotStatus: domain.SnapshotStatusActive,
		GeneratedBy:    &currentUser,
		IsActive:       true,
	}

	penetapanId, err := s.Repo.SavePenetapanOpd(ctx, tx, penetapanTujuan)
	if err != nil {
		s.Logger.Error("save penetapan opd error",
			"penetapanTujuan", penetapanTujuan)
		return web.SyncPenetapanOpdResponse{}, err
	}
	snapshotTujuan, err := s.toTujuanSnapshot(tujuanPerencanaan, currentUser, penetapanId, req.Tahun)
	if err != nil {
		errUpdate := s.markSyncAsFailed(
			ctx,
			syncId,
			err.Error(),
		)
		if errUpdate != nil {
			return web.SyncPenetapanOpdResponse{}, errUpdate
		}

		return web.SyncPenetapanOpdResponse{}, errUpdate
	}
	// simpan tujuan
	var jumlahTujuan int
	var jumlahIndikator int
	var jumlahTarget int
	for _, tujuan := range snapshotTujuan {
		tujuanId, err := s.Repo.SaveTujuanPenetapanOpd(ctx, tx, tujuan)
		if err != nil {
			errUpdate := s.markSyncAsFailed(
				ctx,
				syncId,
				err.Error(),
			)
			if errUpdate != nil {
				return web.SyncPenetapanOpdResponse{}, errUpdate
			}
			return web.SyncPenetapanOpdResponse{}, err
		}
		jumlahTujuan += 1

		// indikator
		for _, indikator := range tujuan.Indikator {

			indId, err := s.Repo.SaveIndikatorTujuanPenetapanOpd(ctx, tx, indikator, tujuanId)
			if err != nil {
				errUpdate := s.markSyncAsFailed(
					ctx,
					syncId,
					err.Error(),
				)
				if errUpdate != nil {
					return web.SyncPenetapanOpdResponse{}, errUpdate
				}
				return web.SyncPenetapanOpdResponse{}, err
			}
			jumlahIndikator += 1

			// target
			targets := indikator.Target
			tgtInserted, err := s.Repo.SaveTargetIndikatorTujuanBatch(ctx, tx, indId, targets)
			if err != nil {
				errUpdate := s.markSyncAsFailed(
					ctx,
					syncId,
					err.Error(),
				)
				if errUpdate != nil {
					return web.SyncPenetapanOpdResponse{}, errUpdate
				}
				return web.SyncPenetapanOpdResponse{}, err
			}
			jumlahTarget += tgtInserted
		}
	}

	// commit transaction penetapan tujuan
	err = tx.Commit()
	if err != nil {
		return web.SyncPenetapanOpdResponse{}, err
	}

	// jika berhasil updateMetadata jadi Success
	err = s.markSyncAsSuccess(ctx, syncId)
	if err != nil {
		return web.SyncPenetapanOpdResponse{}, err
	}

	// just for response
	statusSync = domain.SyncStatusSuccess

	processedAt := time.Now()

	return web.SyncPenetapanOpdResponse{
		SyncId:         syncId,
		Status:         statusSync,
		KodeOpd:        req.KodeOpd,
		Tahun:          req.Tahun,
		JenisPenetapan: jenisPenetapan,
		ProcessedAt:    processedAt,
		ProcessedSummary: web.SyncPenetapanOpdSummary{
			Tujuan:    &jumlahTujuan,
			Indikator: 3,
			Target:    1,
		},
	}, nil
}

func (s *PenetapanOpdService) FindTujuan(
	ctx context.Context,
	req domain.PenetapanOpdRequest,
) ([]web.TujuanPenetapanOpdResponse, error) {
	jenisPenetapan := domain.JenisPenetapanTujuan
	snapshotActive, err := s.Repo.GetActiveSnapshot(ctx, req.KodeOpd, jenisPenetapan, req.Tahun)
	if err != nil {
		return nil, err
	}
	req.SnapshotId = snapshotActive

	tujuanOpd, err := s.Repo.FindTujuanBySnapshot(ctx, req)
	if err != nil {
		return nil, err
	}
	tujuanOpdIds := make([]int64, 0, len(tujuanOpd))
	for _, tj := range tujuanOpd {
		tujuanOpdIds = append(tujuanOpdIds, tj.Id)
	}
	indikators, err := s.Repo.FindIndikatorTujuanByTujuanIds(ctx, tujuanOpdIds)
	if err != nil {
		return nil, err
	}
	indikatorIds := make([]int64, 0, len(indikators))
	for _, ind := range indikators {
		indikatorIds = append(indikatorIds, ind.Id)
	}
	targets, err := s.Repo.FindTargetIndikatorTujuanByIndikatorIds(ctx, indikatorIds)
	if err != nil {
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

func (s *PenetapanOpdService) toTujuanSnapshot(tujuanPerencanaan []perencanaan.PerencanaanTujuanOpdResponse, currentUser string, penetapanId int64, tahunAktif int) ([]domain.TujuanPenetapanOpd, error) {
	penetapanTujuanOpds := make([]domain.TujuanPenetapanOpd, 0, len(tujuanPerencanaan))
	for _, per := range tujuanPerencanaan {
		for _, tujuan := range per.TujuanOpd {
			kodeTujuan := fmt.Sprintf("TUJ-OPD-%d", tujuan.Id)
			periodeTujuan := fmt.Sprintf("%s-%s",
				tujuan.TahunAwal,
				tujuan.TahunAkhir)
			indikators := make([]domain.IndikatorTujuanPenetapanOpd, 0, len(tujuan.Indikator))
			for _, ind := range tujuan.Indikator {
				targets := make([]domain.TargetIndikatorTujuanPenetapanOpd, 0, len(ind.Target))
				for _, tgt := range ind.Target {
					tahunTarget, err := strconv.Atoi(tgt.Tahun)
					if err != nil {
						return nil, fmt.Errorf(
							"invalid tahun target: %s",
							tgt.Tahun,
						)
					}
					target, err := strconv.ParseFloat(tgt.TargetIndikator, 64)
					if err != nil {
						return nil, fmt.Errorf(

							"invalid target indikator: %s",
							tgt.TargetIndikator,
						)
					}
					targets = append(targets,
						domain.TargetIndikatorTujuanPenetapanOpd{
							Tahun:       tahunTarget,
							Target:      target,
							Satuan:      tgt.SatuanIndikator,
							CreatedBy:   &currentUser,
							PenetapanId: penetapanId,
						})
				}
				indikators = append(indikators,
					domain.IndikatorTujuanPenetapanOpd{
						KodeIndikator:       ind.KodeIndikator,
						KodeOpd:             per.KodeOpd,
						Indikator:           ind.NamaIndikator,
						RumusPerhitungan:    &ind.RumusPerhitungan,
						SumberData:          &ind.SumberData,
						DefinisiOperasional: &ind.DefinisiOperasional,
						TahunAktif:          tahunAktif,
						CreatedBy:           &currentUser,
						PenetapanId:         penetapanId,
						Target:              targets,
					})
			}
			penetapanTujuanOpds = append(penetapanTujuanOpds,
				domain.TujuanPenetapanOpd{
					KodeOpd:       per.KodeOpd,
					KodeTujuanOpd: &kodeTujuan,
					TujuanOpd:     tujuan.Tujuan,
					Periode:       periodeTujuan,
					TahunAktif:    tahunAktif,
					CreatedBy:     &currentUser,
					PenetapanId:   penetapanId,
					Indikator:     indikators,
				})
		}
	}

	return penetapanTujuanOpds, nil
}
