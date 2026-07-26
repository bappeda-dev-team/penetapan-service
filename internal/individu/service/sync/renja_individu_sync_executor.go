package sync

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/bappeda-dev-team/penetapan-service/internal/individu/client"
	"github.com/bappeda-dev-team/penetapan-service/internal/individu/domain"
	"github.com/bappeda-dev-team/penetapan-service/internal/individu/repository"
	"github.com/bappeda-dev-team/penetapan-service/internal/individu/web"
)

type RenjaIndividuSyncExecutor struct {
	Repo   *repository.PenetapanIndividuRepository
	Client *client.Client
	Logger *slog.Logger
}

func NewRenjaIndividuSyncExecutor(
	repo *repository.PenetapanIndividuRepository,
	client *client.Client,
	logger *slog.Logger,
) *RenjaIndividuSyncExecutor {
	return &RenjaIndividuSyncExecutor{
		Repo:   repo,
		Client: client,
		Logger: logger,
	}
}

func (ex *RenjaIndividuSyncExecutor) Sync(
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
	renjaIndividus, err := ex.Client.SyncRenjaIndividu(ctx, request)
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
		JenisSnapshot: domain.JenisPenetapanRenjaIndividu,
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
	for _, renja := range renjaIndividus {
		renjaIndividu := domain.RenjaIndividu{
			PegawaiId:  snapshot.PegawaiId,
			KodeOpd:    snapshot.KodeOpd,
			TahunAktif: snapshot.Tahun,

			KodePk:        renja.IdRekinPemilikPk,
			NamaPemilikPk: renja.NamaPemilikPk,
			LevelPk:       renja.LevelPk,

			KodeProgram: renja.KodeProgram,
			NamaProgram: renja.NamaProgram,
			PaguProgram: renja.PaguProgram,

			KodeKegiatan: renja.KodeKegiatan,
			NamaKegiatan: renja.NamaKegiatan,
			PaguKegiatan: renja.PaguKegiatan,

			KodeSubkegiatan: renja.KodeSubkegiatan,
			NamaSubkegiatan: renja.NamaSubkegiatan,
			PaguSubkegiatan: renja.PaguSubkegiatan,

			PenetapanIndividuId: snapshotId,
			CreatedBy:           &currentUser,
		}
		_, err := ex.Repo.SaveRenjaIndividu(ctx, tx, renjaIndividu)
		if err != nil {
			return web.SyncPenetapanSummary{}, err
		}
		summary.AddRenjaIndividu(1)
	}

	err = tx.Commit()
	if err != nil {
		return web.SyncPenetapanSummary{}, err
	}

	return summary.Response(), nil
}

func (ex *RenjaIndividuSyncExecutor) createActiveSnapshot(
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
		ex.Logger.Error("save renja individu error", "snapshot", snapshot)
		return 0, err
	}

	return snapshotID, nil
}
