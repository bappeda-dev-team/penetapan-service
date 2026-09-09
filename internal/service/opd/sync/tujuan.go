package sync

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/bappeda-dev-team/penetapan-service/internal/client/perencanaan"
	"github.com/bappeda-dev-team/penetapan-service/internal/common"
	"github.com/bappeda-dev-team/penetapan-service/internal/model/domain"
	"github.com/bappeda-dev-team/penetapan-service/internal/model/web"
	"github.com/bappeda-dev-team/penetapan-service/internal/repository"
	"github.com/bappeda-dev-team/penetapan-service/internal/service/opd/helper"
)

type TujuanSyncExecutor struct {
	Repo              *repository.PenetapanOpdRepository
	PerencanaanClient *perencanaan.PerencanaanClient
	Logger            *slog.Logger
}

func NewTujuanSyncExecutor(
	repo *repository.PenetapanOpdRepository,
	perencanaanClient *perencanaan.PerencanaanClient,
	logger *slog.Logger,
) *TujuanSyncExecutor {
	return &TujuanSyncExecutor{
		Repo:              repo,
		PerencanaanClient: perencanaanClient,
		Logger:            logger,
	}
}

func (ex *TujuanSyncExecutor) Sync(
	ctx context.Context,
	syncId int64,
	req *web.SyncPenetapanOpdRequest,
	currentUser string,

) (web.SyncPenetapanOpdSummary, error) {
	perencanaanRequest := perencanaan.PerencanaanRequest{
		KodeOpd: req.KodeOpd,
		Tahun:   req.Tahun,
	}

	// outsource
	tujuanPerencanaan, err := ex.PerencanaanClient.GetPenetapanTujuanOpd(ctx, perencanaanRequest)
	if err != nil {
		return web.SyncPenetapanOpdSummary{}, err
	}

	if !hasValidTujuanOpd(tujuanPerencanaan) {
		return web.SyncPenetapanOpdSummary{}, common.NewValidation(
			"tidak ada data penetapan tujuan OPD yang siap disinkronkan",
		)
	}

	// mulai butuh tx
	tx, err := ex.Repo.DB.BeginTx(ctx, nil)
	if err != nil {
		return web.SyncPenetapanOpdSummary{}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	versiPenetapan, err := ex.Repo.GetPenetapanNextVersion(ctx, tx, req.KodeOpd, domain.JenisPenetapanTujuan, req.Tahun)
	if err != nil {
		return web.SyncPenetapanOpdSummary{}, err
	}
	if versiPenetapan > 1 {
		errDeact := ex.Repo.DeactivateOldSnapshot(ctx, tx, req.KodeOpd, domain.JenisPenetapanTujuan, req.Tahun)
		if errDeact != nil {
			return web.SyncPenetapanOpdSummary{}, errDeact
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

	penetapanId, err := ex.Repo.SavePenetapanOpd(ctx, tx, penetapanTujuan)
	if err != nil {
		ex.Logger.Error("save penetapan opd error",
			"penetapanTujuan", penetapanTujuan)
		return web.SyncPenetapanOpdSummary{}, err
	}

	// convert response perencanaan to snapshot
	snapshotTujuan, err := ex.toTujuanSnapshots(tujuanPerencanaan, currentUser, penetapanId, req.Tahun)
	if err != nil {
		return web.SyncPenetapanOpdSummary{}, err
	}

	// simpan tujuan
	var jumlahTujuan int
	var jumlahIndikator int
	var jumlahTarget int
	for _, tujuan := range snapshotTujuan {
		tujuanId, err := ex.Repo.SaveTujuanPenetapanOpd(ctx, tx, tujuan)
		if err != nil {
			return web.SyncPenetapanOpdSummary{}, err
		}
		jumlahTujuan += 1

		// indikator
		for _, indikator := range tujuan.Indikator {

			indId, err := ex.Repo.SaveIndikatorTujuanPenetapanOpd(ctx, tx, indikator, tujuanId)
			if err != nil {
				return web.SyncPenetapanOpdSummary{}, err
			}
			jumlahIndikator += 1

			// target
			targets := indikator.Target
			tgtInserted, err := ex.Repo.SaveTargetIndikatorTujuanBatch(ctx, tx, indId, targets)
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
		Tujuan:    &jumlahTujuan,
		Indikator: jumlahIndikator,
		Target:    jumlahTarget,
	}, nil
}

func (ex *TujuanSyncExecutor) toTujuanSnapshots(tujuanPerencanaan []perencanaan.PerencanaanTujuanOpdResponse, currentUser string, penetapanId int64, tahunAktif int) ([]domain.TujuanPenetapanOpd, error) {

	createdBy := &currentUser
	var penetapanTujuanOpds = []domain.TujuanPenetapanOpd{}
	for _, per := range tujuanPerencanaan {
		for _, tujuan := range per.TujuanOpd {

			indikators := make(
				[]domain.IndikatorTujuanPenetapanOpd,
				0,
				len(tujuan.Indikator),
			)

			for _, ind := range tujuan.Indikator {
				indSnapshot, err := ex.toIndikatorTujuanSnapshot(
					ind, per.KodeOpd,
					createdBy,
					tahunAktif,
				)
				if err != nil {
					return nil, err
				}
				indikators = append(
					indikators,
					indSnapshot,
				)

			}

			tujSnapshot := ex.toTujuanSnapshot(
				tujuan,
				per.KodeOpd,
				tahunAktif,
				penetapanId,
				createdBy,
				indikators,
			)
			penetapanTujuanOpds = append(penetapanTujuanOpds,
				tujSnapshot,
			)
		}
	}

	return penetapanTujuanOpds, nil
}

func (ex *TujuanSyncExecutor) toTujuanSnapshot(
	tujuan perencanaan.TujuanOpdResponse,
	kodeOpd string,
	tahunAktif int,
	penetapanId int64,
	createdBy *string,
	indikators []domain.IndikatorTujuanPenetapanOpd,
) domain.TujuanPenetapanOpd {
	kodeTujuan := fmt.Sprintf("TUJ-OPD-%d", tujuan.Id)
	periodeTujuan := fmt.Sprintf("%s-%s",
		tujuan.TahunAwal,
		tujuan.TahunAkhir)
	return domain.TujuanPenetapanOpd{
		KodeOpd:       kodeOpd,
		KodeTujuanOpd: kodeTujuan,
		TujuanOpd:     tujuan.Tujuan,
		Periode:       periodeTujuan,
		TahunAktif:    tahunAktif,
		CreatedBy:     createdBy,
		PenetapanId:   penetapanId,
		Indikator:     indikators,
	}
}

func (ex *TujuanSyncExecutor) toIndikatorTujuanSnapshot(ind perencanaan.IndikatorResponse, kodeOpd string, createdBy *string, tahunAktif int) (domain.IndikatorTujuanPenetapanOpd, error) {
	targets, err := ex.toTargetSnapshots(ind.Target, createdBy)
	if err != nil {
		return domain.IndikatorTujuanPenetapanOpd{}, err
	}
	kodeIndikator := fmt.Sprintf("IND-%s", ind.Id)
	return domain.IndikatorTujuanPenetapanOpd{
		KodeIndikator:       kodeIndikator,
		KodeOpd:             kodeOpd,
		Indikator:           ind.NamaIndikator,
		RumusPerhitungan:    &ind.RumusPerhitungan,
		SumberData:          &ind.SumberData,
		DefinisiOperasional: &ind.DefinisiOperasional,
		TahunAktif:          tahunAktif,
		CreatedBy:           createdBy,
		Target:              targets,
	}, nil
}

func (ex *TujuanSyncExecutor) toTargetSnapshots(targets []perencanaan.TargetResponse, createdBy *string) ([]domain.TargetIndikatorTujuanPenetapanOpd, error) {
	result := make([]domain.TargetIndikatorTujuanPenetapanOpd, 0, len(targets))
	for _, tgt := range targets {
		tahunTarget, errTahun := helper.ParseTahun(tgt.Tahun)
		if errTahun != nil {
			return nil, errTahun
		}
		target, errConv := helper.ParseFloat(tgt.TargetIndikator)
		if errConv != nil {
			return nil, errConv
		}
		tgtSnapshot := ex.toTargetIndikatorTujuanSnapshot(tgt, tahunTarget, target, createdBy)
		result = append(result, tgtSnapshot)
	}
	return result, nil
}

func (ex *TujuanSyncExecutor) toTargetIndikatorTujuanSnapshot(tgt perencanaan.TargetResponse, tahunTarget int, target float64, createdBy *string) domain.TargetIndikatorTujuanPenetapanOpd {
	kodeTarget := fmt.Sprintf("TGT-%s", tgt.Id)
	return domain.TargetIndikatorTujuanPenetapanOpd{
		KodeTarget: kodeTarget,
		Tahun:      tahunTarget,
		Target:     target,
		Satuan:     tgt.SatuanIndikator,
		CreatedBy:  createdBy,
	}

}

func hasValidTujuanOpd(tujuanPerencanaans []perencanaan.PerencanaanTujuanOpdResponse) bool {
	for _, per := range tujuanPerencanaans {
		for _, tujuan := range per.TujuanOpd {
			if strings.TrimSpace(tujuan.Tujuan) == "" {
				continue
			}
			for _, indikator := range tujuan.Indikator {
				if strings.TrimSpace(indikator.NamaIndikator) == "" || len(indikator.Target) == 0 {
					continue
				}
				for _, target := range indikator.Target {
					if isValidTargetOpd(target) {
						return true
					}
				}
			}
		}
	}
	return false
}

func isValidTargetOpd(target perencanaan.TargetResponse) bool {
	val, err := helper.ParseFloat(target.TargetIndikator)
	if err != nil {
		return false
	}
	if _, errTahun := helper.ParseTahun(target.Tahun); errTahun != nil {
		return false
	}
	return val != 0 && strings.TrimSpace(target.SatuanIndikator) != ""
}
