package sync

import (
	"context"
	"log/slog"

	"github.com/bappeda-dev-team/penetapan-service/internal/client/perencanaan"
	"github.com/bappeda-dev-team/penetapan-service/internal/kode"
	"github.com/bappeda-dev-team/penetapan-service/internal/model/domain"
	"github.com/bappeda-dev-team/penetapan-service/internal/model/web"
	"github.com/bappeda-dev-team/penetapan-service/internal/repository"
)

type RenaksiSyncExecutor struct {
	Repo              *repository.PenetapanOpdRepository
	PerencanaanClient *perencanaan.PerencanaanClient
	Logger            *slog.Logger
}

func NewRenaksiSyncExecutor(
	repo *repository.PenetapanOpdRepository,
	perencanaanClient *perencanaan.PerencanaanClient,
	logger *slog.Logger,
) *RenaksiSyncExecutor {
	return &RenaksiSyncExecutor{
		Repo:              repo,
		PerencanaanClient: perencanaanClient,
		Logger:            logger,
	}
}

func (ex *RenaksiSyncExecutor) Sync(
	ctx context.Context,
	syncId int64,
	req *web.SyncPenetapanOpdRequest,
	currentUser string,
) (web.SyncPenetapanOpdSummary, error) {
	const op = "RENAKSI OPD SYNC EXECUTOR"
	perencanaanRequest := perencanaan.PerencanaanRequest{
		KodeOpd: req.KodeOpd,
		Tahun:   req.Tahun,
	}
	perencanaanResponse, err := ex.PerencanaanClient.GetPenetapanRenaksiOpd(ctx, perencanaanRequest)
	if err != nil {
		ex.Logger.Error(op, "client", "RENAKSI OPD", "err", err)
		return web.SyncPenetapanOpdSummary{}, err
	}
	ex.Logger.Info(op, "perencanaanResponse", perencanaanResponse)

	// TODO: implement guard for penetapan
	// if !hasValidRenaksi(perencanaanResponse) {
	// 	return web.SyncPenetapanOpdSummary{}, common.NewValidation(
	// 		"tidak ada data penetapan renaksi OPD yang siap disinkronkan",
	// 	)
	// }

	// mulai butuh tx
	tx, err := ex.Repo.DB.BeginTx(ctx, nil)
	if err != nil {
		return web.SyncPenetapanOpdSummary{}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// cari latest versi penetapan renja
	// versi otomatis + 1 misal 0 / belum ada -> 1
	versiPenetapan, err := ex.Repo.GetPenetapanNextVersion(ctx, tx,
		req.KodeOpd, domain.JenisPenetapanRenaksi, req.Tahun)
	if err != nil {
		ex.Logger.Error(op, "errVersiPenetapan", err)
		return web.SyncPenetapanOpdSummary{}, err
	}
	if versiPenetapan > 1 {
		ex.Logger.Info(op, "Versi Penetapan", versiPenetapan)
		// deactivate versi lama
		errDeact := ex.Repo.DeactivateOldSnapshot(ctx, tx,
			req.KodeOpd, domain.JenisPenetapanRenaksi, req.Tahun)
		if errDeact != nil {
			ex.Logger.Error(op, "errDeact", errDeact)
			return web.SyncPenetapanOpdSummary{}, err
		}
	}

	// snapshot penetapan baru
	penetapan := domain.PenetapanOpd{
		KodeOpd:        req.KodeOpd,
		Tahun:          req.Tahun,
		JenisPenetapan: domain.JenisPenetapanRenaksi,
		Versi:          versiPenetapan,
		SnapshotStatus: domain.SnapshotStatusActive,
		GeneratedBy:    &currentUser,
		IsActive:       true,
	}
	// save snapshot penetapan baru
	penetapanId, err := ex.Repo.SavePenetapanOpd(ctx, tx, penetapan)
	if err != nil {
		ex.Logger.Error(op,
			"penetapan id renaksi err", err)
		return web.SyncPenetapanOpdSummary{}, err
	}

	// convert response perencanaan to snapshot
	snapshotRenaksis := ex.toRenaksiSnapshots(
		perencanaanResponse,
		&currentUser,
		penetapanId,
		req.KodeOpd,
		req.Tahun,
	)
	ex.Logger.Info(op, "snapshotRenaksis", snapshotRenaksis)

	jumlahRenaksiTersimpan, err := ex.Repo.SaveRencanaAksiOpdBatch(ctx, tx, snapshotRenaksis)
	if err != nil {
		return web.SyncPenetapanOpdSummary{}, err
	}
	ex.Logger.Info(op, "jumlah renaksi", jumlahRenaksiTersimpan)

	// commit transaction penetapan renaksi opd
	err = tx.Commit()
	if err != nil {
		ex.Logger.Error(op, "error commit", err)
		return web.SyncPenetapanOpdSummary{}, err
	}

	return web.SyncPenetapanOpdSummary{
		Renaksi: &jumlahRenaksiTersimpan,
	}, nil
}

func (ex *RenaksiSyncExecutor) toRenaksiSnapshots(
	renaksiPerencanaans []perencanaan.RencanaAksiOpdResponse,
	createdBy *string,
	penetapanId int64,
	kodeOpd string,
	tahun int,
) []domain.RencanaAksiOpd {
	snapshots := make([]domain.RencanaAksiOpd, 0, len(renaksiPerencanaans))

	for _, ren := range renaksiPerencanaans {
		snapshots = append(snapshots, domain.RencanaAksiOpd{
			PenetapanId:        penetapanId,
			KodeRencanaAksiOpd: kode.KodeRenaksiOpd(ren.Id),
			KodeOpd:            kodeOpd,
			KodeSasaranOpd:     kode.KodeSasaranOpd(ren.SasaranId),
			KodePk:             ren.RekinId,
			NamaPk:             ren.AksiKegiatan,
			PegawaiId:          ren.NamaPemilik,
			KodeSubkegiatan:    ren.KodeSubkegiatan,
			NamaSubkegiatan:    ren.NamaSubKegiatan,
			AnggaranRenaksi:    ren.Anggaran,
			Tahun:              tahun,
			Tw1:                ren.Tw1,
			Tw2:                ren.Tw2,
			Tw3:                ren.Tw3,
			Tw4:                ren.Tw4,
			CreatedBy:          createdBy,
		})
	}

	return snapshots
}
