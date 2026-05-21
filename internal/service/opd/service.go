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
	jenisPenetapan domain.JenisPenetapan,
) (web.SyncPenetapanOpdResponse, error) {
	now := time.Now()
	currentUser := "super_admin" // ambil dari ctx nanti

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
) (web.TujuanPenetapanOpdResponse, error) {
	s.Logger.Info("FindTujuan")

	jenisPenetapan := domain.JenisPenetapanTujuan
	snapshot, err := s.getActiveSnapshot(ctx, req.KodeOpd, jenisPenetapan, req.Tahun)
	if err != nil {
		return web.TujuanPenetapanOpdResponse{}, err
	}
	req.SnapshotId = &snapshot.Id

	tujuanOpd, err := s.Repo.FindTujuanBySnapshot(ctx, req)
	if err != nil {
		s.Logger.Error("FindTujuanBySnapshot")
		return web.TujuanPenetapanOpdResponse{}, err
	}
	tujuanOpdIds := make([]int64, 0, len(tujuanOpd))
	for _, tj := range tujuanOpd {
		tujuanOpdIds = append(tujuanOpdIds, tj.Id)
	}
	indikators, err := s.Repo.FindIndikatorTujuanByTujuanIds(ctx, tujuanOpdIds)
	if err != nil {
		s.Logger.Error("FindIndikatorTujuanByTujuanIds")
		return web.TujuanPenetapanOpdResponse{}, err
	}
	indikatorIds := make([]int64, 0, len(indikators))
	for _, ind := range indikators {
		indikatorIds = append(indikatorIds, ind.Id)
	}
	targets, err := s.Repo.FindTargetIndikatorTujuanByIndikatorIds(ctx, indikatorIds)
	if err != nil {
		s.Logger.Error("FindTargetIndikatorTujuanByIndikatorIds")
		return web.TujuanPenetapanOpdResponse{}, err
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

	indikatorMap := make(map[int64][]web.IndikatorTujuanPenetapanResponse)
	for _, indikator := range indikators {
		indikatorResp := ToIndikatorTujuanOpdResponse(indikator)
		indikatorResp.Target = targetMap[indikator.Id]
		indikatorMap[indikator.IdTujuanOpd] = append(
			indikatorMap[indikator.IdTujuanOpd],
			indikatorResp,
		)
	}

	tujuanResponse := make([]web.TujuanOpdResponse, 0, len(tujuanOpd))
	for _, tujuan := range tujuanOpd {
		tujResp := ToTujuanOpdResponse(tujuan)
		tujResp.Indikator = indikatorMap[tujuan.Id]
		tujuanResponse = append(tujuanResponse, tujResp)
	}

	return web.TujuanPenetapanOpdResponse{
		KodeOpd:    req.KodeOpd,
		TahunAktif: req.Tahun,
		Versi:      snapshot.Versi,
		IsLocked:   true,
		Tujuans:    tujuanResponse,
	}, nil
}

func (s *PenetapanOpdService) FindSasaran(
	ctx context.Context,
	req domain.PenetapanOpdRequest,
) (web.SasaranPenetapanOpdResponse, error) {
	s.Logger.Info("FindSasaran")

	jenisPenetapan := domain.JenisPenetapanSasaran
	snapshot, err := s.getActiveSnapshot(ctx, req.KodeOpd, jenisPenetapan, req.Tahun)
	if err != nil {
		return web.SasaranPenetapanOpdResponse{}, err
	}
	req.SnapshotId = &snapshot.Id
	sasaranOpds, err := s.Repo.FindSasaranBySnapshot(ctx, req)
	if err != nil {
		s.Logger.Error("FindSasaranBySnapshot")
		return web.SasaranPenetapanOpdResponse{}, err
	}
	sasaranOpdIds := make([]int64, 0, len(sasaranOpds))
	for _, sas := range sasaranOpds {
		sasaranOpdIds = append(sasaranOpdIds, sas.Id)
	}
	indikators, err := s.Repo.FindIndikatorSasaranBySasarands(ctx, sasaranOpdIds)
	if err != nil {
		s.Logger.Error("FindIndikatorSasaranBySasarands")
		return web.SasaranPenetapanOpdResponse{}, err
	}
	indikatorIds := make([]int64, 0, len(indikators))
	for _, ind := range indikators {
		indikatorIds = append(indikatorIds, ind.Id)
	}

	targets, err := s.Repo.FindTargetIndikatorSasaranByIndikatorIds(ctx, indikatorIds)
	if err != nil {
		s.Logger.Error("FindTargetIndikatorSasaranByIndikatorIds")
		return web.SasaranPenetapanOpdResponse{}, err
	}
	targetMap := make(
		map[int64][]web.TargetIndikatorResponse,
	)
	for _, target := range targets {
		targetResp := ToTargetIndikatorSasaranOpdResponse(target)
		targetMap[target.IndikatorSasaranId] = append(
			targetMap[target.IndikatorSasaranId],
			targetResp,
		)
	}

	indikatorMap := make(map[int64][]web.IndikatorSasaranPenetapanResponse)
	for _, indikator := range indikators {
		indikatorResp := ToIndikatorSasaranOpdResponse(indikator)
		indikatorResp.Target = targetMap[indikator.Id]
		indikatorMap[indikator.IdSasaranOpd] = append(
			indikatorMap[indikator.IdSasaranOpd],
			indikatorResp,
		)
	}

	sasaranResponse := make([]web.SasaranOpdResponse, 0, len(sasaranOpds))
	for _, sasaran := range sasaranOpds {
		sasaranResp := ToSasaranOpdResponse(sasaran)
		sasaranResp.Indikator = indikatorMap[sasaran.Id]
		sasaranResponse = append(sasaranResponse, sasaranResp)
	}

	return web.SasaranPenetapanOpdResponse{
		KodeOpd:    req.KodeOpd,
		TahunAktif: req.Tahun,
		Versi:      snapshot.Versi,
		IsLocked:   true,
		Sasarans:   sasaranResponse,
	}, nil
}

func (s *PenetapanOpdService) FindRenja(
	ctx context.Context,
	req domain.PenetapanOpdRequest,
) (web.RenjaPenetapanOpdResponse, error) {
	s.Logger.Info("FindRenja")

	jenisPenetapan := domain.JenisPenetapanRenja
	snapshot, err := s.getActiveSnapshot(ctx, req.KodeOpd, jenisPenetapan, req.Tahun)
	if err != nil {
		return web.RenjaPenetapanOpdResponse{}, err
	}

	req.SnapshotId = &snapshot.Id

	urusans, err := s.Repo.FindRenjaUrusanBySnapshot(ctx, req)
	if err != nil {
		s.Logger.Error("FindRenjaUrusanBySnapshot")
		return web.RenjaPenetapanOpdResponse{}, err
	}
	err = s.findPaguRenjaUrusan(ctx, &urusans)
	if err != nil {
		s.Logger.Error("findPaguRenjaUrusan")
		return web.RenjaPenetapanOpdResponse{}, err
	}

	bidangUrusans, err := s.Repo.FindRenjaBidangUrusanBySnapshot(ctx, req)
	if err != nil {
		s.Logger.Error("FindRenjaBidangUrusanBySnapshot")
		return web.RenjaPenetapanOpdResponse{}, err
	}
	err = s.findPaguRenjaBidangUrusan(ctx, &bidangUrusans)
	if err != nil {
		s.Logger.Error("findPaguRenjaBidangUrusan")
		return web.RenjaPenetapanOpdResponse{}, err
	}

	programs, err := s.Repo.FindRenjaProgramBySnapshot(ctx, req)
	if err != nil {
		s.Logger.Error("FindRenjaProgramBySnapshot")
		return web.RenjaPenetapanOpdResponse{}, err
	}
	err = s.findIndikatorRenjaProgram(ctx, &programs)
	if err != nil {
		s.Logger.Error("findIndikatorRenjaProgram")
		return web.RenjaPenetapanOpdResponse{}, err
	}
	err = s.findPaguRenjaProgram(ctx, &programs)
	if err != nil {
		s.Logger.Error("findPaguRenjaProgram")
		return web.RenjaPenetapanOpdResponse{}, err
	}

	kegiatans, err := s.Repo.FindRenjaKegiatanBySnapshot(ctx, req)
	if err != nil {
		s.Logger.Error("FindRenjaKegiatanBySnapshot")
		return web.RenjaPenetapanOpdResponse{}, err
	}
	err = s.findIndikatorRenjaKegiatan(ctx, &kegiatans)
	if err != nil {
		s.Logger.Error("findIndikatorRenjaKegiatan")
		return web.RenjaPenetapanOpdResponse{}, err
	}
	err = s.findPaguRenjaKegiatan(ctx, &kegiatans)
	if err != nil {
		s.Logger.Error("findPaguRenjaKegiatan")
		return web.RenjaPenetapanOpdResponse{}, err
	}

	subkegiatans, err := s.Repo.FindRenjaSubkegiatanBySnapshot(ctx, req)
	if err != nil {
		s.Logger.Error("FindRenjaSubkegiatanBySnapshot")
		return web.RenjaPenetapanOpdResponse{}, err
	}
	err = s.findIndikatorRenjaSubkegiatan(ctx, &subkegiatans)
	if err != nil {
		s.Logger.Error("findIndikatorRenjaSubkegiatan")
		return web.RenjaPenetapanOpdResponse{}, err
	}
	err = s.findPaguRenjaSubkegiatan(ctx, &subkegiatans)
	if err != nil {
		s.Logger.Error("findPaguRenjaSubkegiatan")
		return web.RenjaPenetapanOpdResponse{}, err
	}

	// responses

	var urusanResponses []web.RenjaUrusanResponse
	for _, urusan := range urusans {
		// pagu urusan
		var anggarans []web.PaguAnggaranRenjaResponse
		for _, pagu := range urusan.PaguAnggaran {
			anggarans = append(anggarans,
				web.PaguAnggaranRenjaResponse{
					Id:        pagu.Id,
					KodePagu:  pagu.KodePagu,
					Pagu:      pagu.PaguAnggaran,
					JenisPagu: pagu.JenisPagu,
				},
			)
		}
		urusanResponses = append(
			urusanResponses,
			web.RenjaUrusanResponse{
				Id:           urusan.Id,
				KodeUrusan:   urusan.KodeUrusan,
				Urusan:       urusan.Urusan,
				IsLocked:     true,
				PaguAnggaran: anggarans,
			},
		)
	}

	var bidangUrusanResponses []web.RenjaBidangUrusanResponse
	for _, bidang := range bidangUrusans {
		// pagu urusan
		var anggarans []web.PaguAnggaranRenjaResponse
		for _, pagu := range bidang.PaguAnggaran {
			anggarans = append(anggarans,
				web.PaguAnggaranRenjaResponse{
					Id:        pagu.Id,
					KodePagu:  pagu.KodePagu,
					Pagu:      pagu.PaguAnggaran,
					JenisPagu: pagu.JenisPagu,
				},
			)
		}
		bidangUrusanResponses = append(
			bidangUrusanResponses,
			web.RenjaBidangUrusanResponse{
				Id:               bidang.Id,
				KodeBidangUrusan: bidang.KodeBidangUrusan,
				BidangUrusan:     bidang.BidangUrusan,
				IsLocked:         true,
				PaguAnggaran:     anggarans,
			},
		)
	}

	var programResponses []web.RenjaProgramResponse
	for _, program := range programs {
		var indResp []web.IndikatorRenjaResponse
		for _, ind := range program.Indikators {
			var targetResp []web.TargetIndikatorResponse
			for _, tar := range ind.Targets {
				targetResp = append(targetResp,
					web.TargetIndikatorResponse{
						Id:         tar.Id,
						KodeTarget: tar.KodeTarget,
						Tahun:      tar.Tahun,
						Target:     tar.Target,
						Satuan:     tar.Satuan,
					})
			}
			indResp = append(indResp,
				web.IndikatorRenjaResponse{
					Id:            ind.Id,
					KodeIndikator: ind.KodeIndikator,
					Indikator:     ind.Indikator,
					Targets:       targetResp,
				})
		}

		// pagu urusan
		var anggarans []web.PaguAnggaranRenjaResponse
		for _, pagu := range program.PaguAnggaran {
			anggarans = append(anggarans,
				web.PaguAnggaranRenjaResponse{
					Id:        pagu.Id,
					KodePagu:  pagu.KodePagu,
					Pagu:      pagu.PaguAnggaran,
					JenisPagu: pagu.JenisPagu,
				},
			)
		}
		programResponses = append(
			programResponses,
			web.RenjaProgramResponse{
				Id:           program.Id,
				KodeProgram:  program.KodeProgram,
				Program:      program.Program,
				IsLocked:     true,
				Indikators:   indResp,
				PaguAnggaran: anggarans,
			},
		)
	}

	var kegiatanResponses []web.RenjaKegiatanResponse
	for _, kegiatan := range kegiatans {
		var indResp []web.IndikatorRenjaResponse
		for _, ind := range kegiatan.Indikators {
			var targetResp []web.TargetIndikatorResponse
			for _, tar := range ind.Targets {
				targetResp = append(targetResp,
					web.TargetIndikatorResponse{
						Id:         tar.Id,
						KodeTarget: tar.KodeTarget,
						Tahun:      tar.Tahun,
						Target:     tar.Target,
						Satuan:     tar.Satuan,
					})
			}
			indResp = append(indResp,
				web.IndikatorRenjaResponse{
					Id:            ind.Id,
					KodeIndikator: ind.KodeIndikator,
					Indikator:     ind.Indikator,
					Targets:       targetResp,
				})
		}

		// pagu urusan
		var anggarans []web.PaguAnggaranRenjaResponse
		for _, pagu := range kegiatan.PaguAnggaran {
			anggarans = append(anggarans,
				web.PaguAnggaranRenjaResponse{
					Id:        pagu.Id,
					KodePagu:  pagu.KodePagu,
					Pagu:      pagu.PaguAnggaran,
					JenisPagu: pagu.JenisPagu,
				},
			)
		}
		kegiatanResponses = append(
			kegiatanResponses,
			web.RenjaKegiatanResponse{
				Id:           kegiatan.Id,
				KodeKegiatan: kegiatan.KodeKegiatan,
				Kegiatan:     kegiatan.Kegiatan,
				IsLocked:     true,
				Indikators:   indResp,
				PaguAnggaran: anggarans,
			},
		)
	}

	var subkegiatanResponses []web.RenjaSubkegiatanResponse
	for _, sub := range subkegiatans {
		var indResp []web.IndikatorRenjaResponse
		for _, ind := range sub.Indikators {
			var targetResp []web.TargetIndikatorResponse
			for _, tar := range ind.Targets {
				targetResp = append(targetResp,
					web.TargetIndikatorResponse{
						Id:         tar.Id,
						KodeTarget: tar.KodeTarget,
						Tahun:      tar.Tahun,
						Target:     tar.Target,
						Satuan:     tar.Satuan,
					})
			}
			indResp = append(indResp,
				web.IndikatorRenjaResponse{
					Id:            ind.Id,
					KodeIndikator: ind.KodeIndikator,
					Indikator:     ind.Indikator,
					Targets:       targetResp,
				})
		}

		// pagu urusan
		var anggarans []web.PaguAnggaranRenjaResponse
		for _, pagu := range sub.PaguAnggaran {
			anggarans = append(anggarans,
				web.PaguAnggaranRenjaResponse{
					Id:        pagu.Id,
					KodePagu:  pagu.KodePagu,
					Pagu:      pagu.PaguAnggaran,
					JenisPagu: pagu.JenisPagu,
				},
			)
		}
		subkegiatanResponses = append(
			subkegiatanResponses,
			web.RenjaSubkegiatanResponse{
				Id:              sub.Id,
				KodeSubkegiatan: sub.KodeSubkegiatan,
				Subkegiatan:     sub.Subkegiatan,
				IsLocked:        true,
				Indikators:      indResp,
				PaguAnggaran:    anggarans,
			},
		)
	}

	return web.RenjaPenetapanOpdResponse{
		KodeOpd:       req.KodeOpd,
		TahunAktif:    req.Tahun,
		Versi:         12,
		IsLocked:      true,
		Urusans:       urusanResponses,
		BidangUrusans: bidangUrusanResponses,
		Programs:      programResponses,
		Kegiatans:     kegiatanResponses,
		Subkegiatans:  subkegiatanResponses,
	}, nil
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

func (s *PenetapanOpdService) getActiveSnapshot(
	ctx context.Context,
	kodeOpd string,
	jenisPenetapan string,
	tahun int,
) (*domain.ActiveSnapshot, error) {

	snapshot, err := s.Repo.GetActiveSnapshot(ctx, kodeOpd, jenisPenetapan, tahun)
	if err != nil {
		s.Logger.Error("GetActiveSnapshot")
		return nil, err
	}
	if snapshot == nil {
		s.Logger.ErrorContext(
			ctx,
			"get active snapshot failed",
			"kode_opd", kodeOpd,
			"jenis", jenisPenetapan,
			"tahun", tahun,
			"err", err,
		)
		return nil, nil
	}

	return snapshot, nil
}
