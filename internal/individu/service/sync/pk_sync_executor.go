package sync

import (
	"context"
	"database/sql"
	"log/slog"
	"strconv"

	"github.com/bappeda-dev-team/penetapan-service/internal/individu/client"
	"github.com/bappeda-dev-team/penetapan-service/internal/individu/domain"
	"github.com/bappeda-dev-team/penetapan-service/internal/individu/repository"
	"github.com/bappeda-dev-team/penetapan-service/internal/individu/web"
	"github.com/bappeda-dev-team/penetapan-service/internal/kode"
)

type PkSyncExecutor struct {
	Repo   *repository.PenetapanIndividuRepository
	Client *client.Client
	Logger *slog.Logger
}

func NewPkSyncExecutor(
	repo *repository.PenetapanIndividuRepository,
	client *client.Client,
	logger *slog.Logger,
) *PkSyncExecutor {
	return &PkSyncExecutor{
		Repo:   repo,
		Client: client,
		Logger: logger,
	}
}

func (ex *PkSyncExecutor) Sync(
	ctx context.Context,
	syncId int64,
	req *web.SyncPenetapanRequest,
	currentUser string,

) (web.SyncPenetapanSummary, error) {
	request := client.SyncRequest{
		PegawaiId: req.PegawaiId,
		KodeOpd:   req.KodeOpd,
		Tahun:     req.Tahun,
	}
	// outsource
	pkPegawais, err := ex.Client.SyncPkPenetapan(ctx, request)
	if err != nil {
		return web.SyncPenetapanSummary{}, err
	}

	// db transaction
	tx, err := ex.Repo.DB.BeginTx(ctx, nil)
	if err != nil {
		return web.SyncPenetapanSummary{}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	snapshot := domain.SnapshotPenetapan{
		JenisSnapshot: domain.JenisPenetapanPk,
		PegawaiId:     req.PegawaiId,
		KodeOpd:       req.KodeOpd,
		Tahun:         req.Tahun,
	}

	snapshotId, err := ex.createActiveSnapshot(ctx, tx, snapshot, currentUser)
	if err != nil {
		return web.SyncPenetapanSummary{}, err
	}

	// SAVE PK
	summary := SummaryCounter{}
	for _, pk := range pkPegawais {
		sasaranOpdIdStr := strconv.Itoa(int(pk.SasaranOpdId))
		kodeSasaranOpd := kode.KodeSasaranOpd(sasaranOpdIdStr)
		pkPenetapan := domain.PkPenetapan{
			PegawaiId:           snapshot.PegawaiId,
			KodeOpd:             snapshot.KodeOpd,
			TahunAktif:          snapshot.Tahun,
			LevelPk:             pk.LevelPk,
			KodeSasaranOpd:      kodeSasaranOpd,
			KodePk:              pk.IdRekinPemilikPk,
			NamaPk:              pk.RekinPemilikPk,
			KeteranganPk:        pk.Keterangan,
			NamaPemilikPk:       pk.NamaPemilikPk,
			AnggaranPk:          pk.AnggaranPk,
			PenetapanIndividuId: snapshotId,
			CreatedBy:           &currentUser,
		}
		pkId, err := ex.Repo.SavePkPenetapan(ctx, tx, pkPenetapan)
		if err != nil {
			return web.SyncPenetapanSummary{}, err
		}
		summary.AddRekin(1)
		for _, ind := range pk.Indikators {
			indikatorPk := domain.IndikatorPk{
				IdPk:            pkId,
				KodeOpd:         snapshot.KodeOpd,
				TahunAktif:      snapshot.Tahun,
				KodeIndikatorPk: ind.IdIndikator,
				NamaIndikatorPk: ind.Indikator,
				CreatedBy:       &currentUser,
			}
			indPkId, err := ex.Repo.SaveIndikatorPkPenetapan(ctx, tx, indikatorPk)
			if err != nil {
				return web.SyncPenetapanSummary{}, err
			}
			summary.AddIndikator(1)
			for _, tgt := range ind.Targets {
				targetFloat, errConv := strconv.ParseFloat(tgt.Target, 64)
				if errConv != nil {
					return web.SyncPenetapanSummary{}, errConv
				}
				targetPk := domain.TargetPk{
					IdIndikatorPk: indPkId,
					KodeTargetPk:  tgt.IdTarget,
					Target:        targetFloat,
					Satuan:        tgt.Satuan,
					Tahun:         snapshot.Tahun,
					CreatedBy:     &currentUser,
				}
				_, err := ex.Repo.SaveTargetPkPenetapan(ctx, tx, targetPk)
				if err != nil {
					return web.SyncPenetapanSummary{}, err
				}
				summary.AddTarget(1)
			}
		}
		for _, ren := range pk.Renaksis {
			renaksiIndividu := domain.RenaksiIndividu{
				IdPk:            pkId,
				KodeOpd:         snapshot.KodeOpd,
				TahunAktif:      snapshot.Tahun,
				KodeRenaksi:     ren.Id,
				Urutan:          ren.Urutan,
				NamaRencanaAksi: ren.NamaRencanaAksi,
				Anggaran:        ren.Anggaran,
				CreatedBy:       &currentUser,
			}
			renId, err := ex.Repo.SaveRenaksiIndividuPkPenetapan(ctx, tx, renaksiIndividu)
			if err != nil {
				return web.SyncPenetapanSummary{}, err
			}
			summary.AddRenaksi(1)
			for _, pl := range ren.Pelaksanaan {
				pelRenaksi := domain.PelaksanaanRenaksi{
					IdRenaksiIndividu: renId,
					KodePelaksanaan:   pl.Id,
					Bulan:             pl.Bulan,
					Bobot:             pl.Bobot,
					CreatedBy:         &currentUser,
				}
				_, err := ex.Repo.SavePelaksananRenaksiIndividuPkPenetapan(ctx, tx, pelRenaksi)
				if err != nil {
					return web.SyncPenetapanSummary{}, err
				}
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return web.SyncPenetapanSummary{}, err
	}

	return summary.Response(), nil
}

func (ex *PkSyncExecutor) createActiveSnapshot(
	ctx context.Context,
	tx *sql.Tx,
	snapshot domain.SnapshotPenetapan,
	currentUser string,
) (int64, error) {
	versi, err := ex.Repo.GetPenetapanNextVersion(ctx, tx, snapshot)
	if err != nil {
		return 0, err
	}

	if versi > 1 {
		if err := ex.Repo.DeactivateOldSnapshot(ctx, tx, snapshot); err != nil {
			return 0, err
		}
	}

	snapshot.Versi = versi
	snapshot.SnapshotStatus = domain.SnapshotStatusActive
	snapshot.GeneratedBy = &currentUser
	snapshot.IsActive = true

	snapshotID, err := ex.Repo.SaveSnapshot(ctx, tx, snapshot)
	if err != nil {
		ex.Logger.Error("save penetapan individu error", "snapshot", snapshot)
		return 0, err
	}

	return snapshotID, nil
}
