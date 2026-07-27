package sync

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

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

	// SAVE RENJA INDIVIDU
	summary := SummaryCounter{}
	for _, renja := range renjaIndividus {
		kodePaguProgram := fmt.Sprintf("PAGU-PRG-%s-%d-%s", renja.KodeProgram, snapshot.Tahun, "programs")
		kodePaguKegiatan := fmt.Sprintf("PAGU-KEG-%s-%d-%s", renja.KodeKegiatan, snapshot.Tahun, "kegiatans")
		kodePaguSubkegiatan := fmt.Sprintf("PAGU-SUBKEG-%s-%d-%s", renja.KodeSubkegiatan, snapshot.Tahun, "subkegiatans")
		renjaIndividu := domain.RenjaIndividu{
			PegawaiId:  snapshot.PegawaiId,
			KodeOpd:    snapshot.KodeOpd,
			TahunAktif: snapshot.Tahun,

			KodePk:        renja.IdRekinPemilikPk,
			NamaPemilikPk: renja.NamaPemilikPk,
			LevelPk:       renja.LevelPk,

			KodeProgram:     renja.KodeProgram,
			NamaProgram:     renja.NamaProgram,
			KodePaguProgram: kodePaguProgram,
			PaguProgram:     renja.PaguProgram,

			KodeKegiatan:     renja.KodeKegiatan,
			NamaKegiatan:     renja.NamaKegiatan,
			KodePaguKegiatan: kodePaguKegiatan,
			PaguKegiatan:     renja.PaguKegiatan,

			KodeSubkegiatan:     renja.KodeSubkegiatan,
			NamaSubkegiatan:     renja.NamaSubkegiatan,
			KodePaguSubkegiatan: kodePaguSubkegiatan,
			PaguSubkegiatan:     renja.PaguSubkegiatan,

			PenetapanIndividuId: snapshotId,
			CreatedBy:           &currentUser,
		}
		idRenja, err := ex.Repo.SaveRenjaIndividu(ctx, tx, renjaIndividu)
		if err != nil {
			return web.SyncPenetapanSummary{}, err
		}
		summary.AddRenjaIndividu(1)
		// INDIKATOR
		if err := ex.saveIndikatorRenja(
			ctx,
			tx,
			idRenja,
			"PROGRAM",
			renja.IndikatorPrograms,
			&summary,
		); err != nil {
			return web.SyncPenetapanSummary{}, err
		}

		if err := ex.saveIndikatorRenja(
			ctx,
			tx,
			idRenja,
			"KEGIATAN",
			renja.IndikatorKegiatans,
			&summary,
		); err != nil {
			return web.SyncPenetapanSummary{}, err
		}

		if err := ex.saveIndikatorRenja(
			ctx,
			tx,
			idRenja,
			"SUB-KEGIATAN",
			renja.IndikatorSubkegiatans,
			&summary,
		); err != nil {
			return web.SyncPenetapanSummary{}, err
		}
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

func (ex *RenjaIndividuSyncExecutor) saveIndikatorRenja(
	ctx context.Context,
	tx *sql.Tx,
	renjaID int64,
	jenis string,
	indikators []client.IndikatorRenja,
	summary *SummaryCounter,
) error {
	for _, ind := range indikators {
		kodeIndikator := fmt.Sprintf("IND-%s", ind.Id)
		indikator := domain.IndikatorRenjaIndividu{
			RenjaIndividuID: renjaID,
			JenisIndikator:  jenis,
			// WARNING KODE INDIKATOR
			KodeIndikatorRenja: kodeIndikator,
			Indikator:          ind.Indikator,
		}

		idIndikator, err := ex.Repo.SaveIndikatorRenjaIndividu(ctx, tx, indikator)
		if err != nil {
			return err
		}

		summary.AddIndikatorRenjaIndividu(1)

		for _, target := range ind.Targets {
			kodeTarget := fmt.Sprintf("TGT-%s", target.Id)
			targetFloat, err := ParseTargetFloat(target.Target)
			if err != nil {
				return fmt.Errorf(
					"invalid target indikator %q value %q",
					ind.Indikator,
					target.Target,
				)
			}

			targetRenja := domain.TargetRenjaIndividu{
				IndikatorRenjaIndividuID: idIndikator,
				JenisTarget:              jenis,
				KodeTargetRenja:          kodeTarget,
				Target:                   targetFloat,
				Satuan:                   target.Satuan,
				Tahun:                    target.Tahun,
			}

			_, err = ex.Repo.SaveTargetRenjaIndividu(ctx, tx, targetRenja)
			if err != nil {
				return err
			}

			summary.AddTargetRenjaIndividu(1)
		}
	}

	return nil
}

func ParseTargetFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)

	// Hanya ubah 10,5 -> 10.5
	if strings.Count(s, ",") == 1 &&
		!strings.Contains(s, ".") {
		s = strings.Replace(s, ",", ".", 1)
	}

	return strconv.ParseFloat(s, 64)
}
