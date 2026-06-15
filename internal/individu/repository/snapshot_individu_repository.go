package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bappeda-dev-team/penetapan-service/internal/individu/domain"
)

func (repo *PenetapanIndividuRepository) GetPenetapanNextVersion(
	ctx context.Context,
	tx *sql.Tx,
	req domain.SnapshotPenetapan,
) (int, error) {
	query := `SELECT COALESCE(MAX(versi), 0) + 1
		  FROM penetapan_individu
 		  WHERE jenis_penetapan = $1
		    AND pegawai_id = $2
		    AND kode_opd = $3
	            AND tahun = $4

		`
	var versi int
	err := tx.QueryRowContext(
		ctx,
		query,
		req.JenisSnapshot,
		req.PegawaiId,
		req.KodeOpd,
		req.Tahun,
	).Scan(&versi)
	if err != nil {
		return 0, err
	}

	return versi, nil
}

func (r *PenetapanIndividuRepository) DeactivateOldSnapshot(
	ctx context.Context,
	tx *sql.Tx,
	req domain.SnapshotPenetapan,
) error {
	snapshotStatus := domain.SnapshotStatusArchived
	query := `UPDATE penetapan_individu
		  SET
		    snapshot_status = $5,
		    is_active = FALSE
 		  WHERE kode_opd = $1
		    AND jenis_penetapan = $2
	            AND tahun = $3
		    AND pegawai_id = $4
		`
	_, err := tx.ExecContext(
		ctx,
		query,
		req.KodeOpd,
		req.JenisSnapshot,
		req.Tahun,
		req.PegawaiId,
		snapshotStatus,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *PenetapanIndividuRepository) SaveSnapshot(
	ctx context.Context,
	tx *sql.Tx,
	snapshot domain.SnapshotPenetapan,
) (int64, error) {
	query := `INSERT INTO penetapan_individu
		(pegawai_id, kode_opd, tahun,
		jenis_penetapan, versi,
		snapshot_status, generated_by,
 		is_active)
		VALUES ($1, $2, $3,
 			$4, $5,
 			$6, $7,
 			$8)
		RETURNING id
		`
	var id int64

	err := tx.QueryRowContext(
		ctx,
		query,
		snapshot.PegawaiId,
		snapshot.KodeOpd,
		snapshot.Tahun,
		snapshot.JenisSnapshot,
		snapshot.Versi,
		snapshot.SnapshotStatus,
		snapshot.GeneratedBy,
		snapshot.IsActive,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *PenetapanIndividuRepository) GetActiveSnapshot(
	ctx context.Context,
	pegawaiId, kodeOpd string,
	tahun int,
	jenisPenetapan domain.JenisPenetapan,
) (*domain.ActiveSnapshot, error) {
	query := `
	    SELECT id, versi
	    FROM penetapan_individu
	    WHERE
		pegawai_id = $1
	    	AND kode_opd = $2
	    	AND tahun = $3
	    	AND jenis_penetapan = $4
	    	AND is_active = TRUE
	    ORDER BY versi DESC
	    LIMIT 1`

	var result domain.ActiveSnapshot
	err := repo.DB.QueryRowContext(
		ctx,
		query,
		pegawaiId,
		kodeOpd,
		tahun,
		jenisPenetapan,
	).Scan(
		&result.Id,
		&result.Versi,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}
