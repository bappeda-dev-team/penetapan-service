package sync

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/bappeda-dev-team/penetapan-service/internal/client/perencanaan"
	"github.com/bappeda-dev-team/penetapan-service/internal/model/domain"
	"github.com/bappeda-dev-team/penetapan-service/internal/model/web"
	"github.com/bappeda-dev-team/penetapan-service/internal/repository"
)

type SasaranSyncExecutor struct {
	Repo              *repository.PenetapanOpdRepository
	PerencanaanClient *perencanaan.PerencanaanClient
	Logger            *slog.Logger
}

func NewSasaranSyncExecutor(
	repo *repository.PenetapanOpdRepository,
	perencanaanClient *perencanaan.PerencanaanClient,
	logger *slog.Logger,
) *SasaranSyncExecutor {
	return &SasaranSyncExecutor{
		Repo:              repo,
		PerencanaanClient: perencanaanClient,
		Logger:            logger,
	}
}

func (ex *SasaranSyncExecutor) Sync(
	ctx context.Context,
	syncId int64,
	req *web.SyncPenetapanOpdRequest,
	currentUser string,

) (web.SyncPenetapanOpdSummary, error) {
	perencanaanRequest := perencanaan.PerencanaanRequest{
		KodeOpd: req.KodeOpd,
		Tahun:   req.Tahun,
	}

	sasaranPerencanaans, err := ex.PerencanaanClient.GetPenetapanSasaranOpd(ctx, perencanaanRequest)
	if err != nil {
		return web.SyncPenetapanOpdSummary{}, err
	}

	// mulai butuh tx
	tx, err := ex.Repo.DB.BeginTx(ctx, nil)
	if err != nil {
		return web.SyncPenetapanOpdSummary{}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	versiPenetapan, err := ex.Repo.GetPenetapanNextVersion(
		ctx, tx, req.KodeOpd, domain.JenisPenetapanSasaran, req.Tahun)
	if err != nil {
		return web.SyncPenetapanOpdSummary{}, err
	}
	if versiPenetapan > 1 {
		errDeact := ex.Repo.DeactivateOldSnapshot(
			ctx, tx, req.KodeOpd, domain.JenisPenetapanSasaran, req.Tahun)
		if errDeact != nil {
			return web.SyncPenetapanOpdSummary{}, err
		}
	}

	// simpan snapshot penetapan
	snapshot := domain.PenetapanOpd{
		KodeOpd:        req.KodeOpd,
		Tahun:          req.Tahun,
		JenisPenetapan: domain.JenisPenetapanSasaran,
		Versi:          versiPenetapan,
		SnapshotStatus: domain.SnapshotStatusActive,
		GeneratedBy:    &currentUser,
		IsActive:       true,
	}

	penetapanId, err := ex.Repo.SavePenetapanOpd(ctx, tx, snapshot)
	if err != nil {
		ex.Logger.Error("save penetapan opd error",
			"snapshot", snapshot)
		return web.SyncPenetapanOpdSummary{}, err
	}

	// convert response perencanaan to snapshot
	snapshotSasaran, err := ex.toSasaranSnapshot(sasaranPerencanaans, currentUser, penetapanId, req.Tahun)
	if err != nil {
		return web.SyncPenetapanOpdSummary{}, err
	}

	// simpan snapshot
	var jumlahSasaran int
	var jumlahIndikator int
	var jumlahTarget int
	for _, sasaran := range snapshotSasaran {
		sasaranId, err := ex.Repo.SaveSasaranPenetapanOpd(ctx, tx, sasaran)
		if err != nil {
			return web.SyncPenetapanOpdSummary{}, err
		}
		jumlahSasaran += 1

		// indikator
		for _, indikator := range sasaran.Indikator {

			indId, err := ex.Repo.SaveIndikatorSasaranPenetapanOpd(ctx, tx, indikator, sasaranId)
			if err != nil {
				return web.SyncPenetapanOpdSummary{}, err
			}
			jumlahIndikator += 1

			// target
			targets := indikator.Target
			tgtInserted, err := ex.Repo.SaveTargetIndikatorSasaranBatch(ctx, tx, indId, targets)
			if err != nil {
				return web.SyncPenetapanOpdSummary{}, err
			}
			jumlahTarget += tgtInserted
		}
	}

	// commit transaction penetapan tujuan
	err = tx.Commit()
	if err != nil {
		return web.SyncPenetapanOpdSummary{}, err
	}

	return web.SyncPenetapanOpdSummary{
		Sasaran:   &jumlahSasaran,
		Indikator: jumlahIndikator,
		Target:    jumlahTarget,
	}, nil
}

func (ex *SasaranSyncExecutor) toSasaranSnapshot(sasaranPerencanaans []perencanaan.PerencanaanSasaranOpdResponse, currentUser string, penetapanId int64, tahunAktif int) ([]domain.SasaranPenetapanOpd, error) {

	createdBy := &currentUser
	var penetapanSasaranOpds = []domain.SasaranPenetapanOpd{}
	for _, per := range sasaranPerencanaans {
		for _, sasaran := range per.SasaranOpd {

			kodeSasaran := fmt.Sprintf("SAS-OPD-%s", sasaran.Id)
			periodeTujuan := fmt.Sprintf("%s-%s",
				sasaran.TahunAwal,
				sasaran.TahunAkhir)

			indikators := make(

				[]domain.IndikatorSasaranPenetapanOpd,
				0,
				len(sasaran.Indikator),
			)

			for _, ind := range sasaran.Indikator {
				indSnapshot, err := ex.toIndikatorSasaranSnapshot(
					ind, per.KodeOpd,
					createdBy,
					tahunAktif,
					penetapanId)
				if err != nil {
					return nil, err
				}
				indikators = append(
					indikators,
					indSnapshot,
				)

			}
			penetapanSasaranOpds = append(penetapanSasaranOpds,
				domain.SasaranPenetapanOpd{
					KodeOpd:        per.KodeOpd,
					KodeSasaranOpd: &kodeSasaran,
					SasaranOpd:     sasaran.NamaSasaranOpd,
					Periode:        periodeTujuan,
					TahunAktif:     tahunAktif,
					CreatedBy:      createdBy,
					PenetapanId:    penetapanId,
					Indikator:      indikators,
				})
		}
	}

	return penetapanSasaranOpds, nil
}

func (ex *SasaranSyncExecutor) toIndikatorSasaranSnapshot(ind perencanaan.IndikatorSasaranResponse, kodeOpd string, createdBy *string, tahunAktif int, penetapanId int64) (domain.IndikatorSasaranPenetapanOpd, error) {
	targets, err := ex.toTargetSnapshots(ind.Target, ind.NamaIndikator, createdBy, penetapanId)
	if err != nil {
		return domain.IndikatorSasaranPenetapanOpd{}, err
	}
	return domain.IndikatorSasaranPenetapanOpd{
		KodeIndikator:       ind.KodeIndikator,
		KodeOpd:             kodeOpd,
		Indikator:           ind.NamaIndikator,
		RumusPerhitungan:    &ind.RumusPerhitungan,
		SumberData:          &ind.SumberData,
		DefinisiOperasional: &ind.DefinisiOperasional,
		TahunAktif:          tahunAktif,
		CreatedBy:           createdBy,
		PenetapanId:         penetapanId,
		Target:              targets,
	}, nil
}

func (ex *SasaranSyncExecutor) toTargetSnapshots(targets []perencanaan.TargetResponse, namaIndikator string, createdBy *string, penetapanId int64) ([]domain.TargetIndikatorSasaranPenetapanOpd, error) {
	result := make([]domain.TargetIndikatorSasaranPenetapanOpd, 0, len(targets))
	for _, tgt := range targets {
		tahunTarget, err := strconv.Atoi(tgt.Tahun)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid tahun indikator %q value %q",
				namaIndikator,
				tgt.TargetIndikator,
			)
		}
		target, err := strconv.ParseFloat(tgt.TargetIndikator, 64)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid target indikator %q value %q",
				namaIndikator,
				tgt.TargetIndikator,
			)
		}
		tgtSnapshot := ex.toTargetIndikatorSasaranSnapshot(tgt, tahunTarget, target, createdBy, penetapanId)
		result = append(result, tgtSnapshot)
	}
	return result, nil
}

func (ex *SasaranSyncExecutor) toTargetIndikatorSasaranSnapshot(tgt perencanaan.TargetResponse, tahunTarget int, target float64, createdBy *string, penetapanId int64) domain.TargetIndikatorSasaranPenetapanOpd {
	return domain.TargetIndikatorSasaranPenetapanOpd{
		Tahun:       tahunTarget,
		Target:      target,
		Satuan:      tgt.SatuanIndikator,
		CreatedBy:   createdBy,
		PenetapanId: penetapanId,
	}
}
