package api

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/bappeda-dev-team/penetapan-service/internal/model/domain"
)

type PenetapanOpdRepository struct {
	DB *sql.DB
}

func (r *PenetapanOpdRepository) SaveTujuanPenetapanOpd(
	ctx context.Context,
	tx *sql.Tx,
	data domain.TujuanPenetapanOpd,
) (int64, error) {

	query := `
		INSERT INTO tb_tujuan_penetapan_opd
		(
			kode_opd,
			kode_tujuan_opd,
			tujuan_opd,
			periode,
			tahun_aktif,
			created_by,
			penetapan_id
		)
		VALUES
		(
			$1, $2, $3, $4, $5, $6, $7
		)
		RETURNING id
	`

	var id int64

	err := tx.QueryRowContext(
		ctx,
		query,
		data.KodeOpd,
		data.KodeTujuanOpd,
		data.TujuanOpd,
		data.Periode,
		data.TahunAktif,
		data.CreatedBy,
		data.PenetapanId,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PenetapanOpdRepository) SaveIndikatorTujuanPenetapanOpd(
	ctx context.Context,
	tx *sql.Tx,
	data domain.IndikatorTujuanPenetapanOpd,
	tujuanId int64,
) (int64, error) {

	query := `
		INSERT INTO tb_indikator_tujuan_penetapan_opd
		(
			id_tujuan_opd,
			kode_opd,
			kode_indikator,
			indikator,
			rumus_perhitungan,
			sumber_data,
			definisi_operasional,
			tahun_aktif,
			created_by,
			penetapan_id
		)
		VALUES
		(
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
		RETURNING id
	`

	var id int64

	err := tx.QueryRowContext(
		ctx,
		query,
		tujuanId,
		data.KodeOpd,
		data.KodeIndikator,
		data.Indikator,
		data.RumusPerhitungan,
		data.SumberData,
		data.DefinisiOperasional,
		data.TahunAktif,
		data.CreatedBy,
		data.PenetapanId,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PenetapanOpdRepository) SaveTargetIndikatorTujuanBatch(
	ctx context.Context,
	tx *sql.Tx,
	indikatorId int64,
	data []domain.TargetIndikatorTujuanPenetapanOpd,
) (int, error) {

	if len(data) == 0 {
		return 0, nil
	}

	var (
		valueStrings []string
		valueArgs    []any
	)

	for i, item := range data {

		base := i * 6

		valueStrings = append(
			valueStrings,
			fmt.Sprintf(
				"($%d, $%d, $%d, $%d, $%d, $%d)",
				base+1,
				base+2,
				base+3,
				base+4,
				base+5,
				base+6,
			),
		)

		valueArgs = append(
			valueArgs,
			indikatorId,
			item.Tahun,
			item.Target,
			item.Satuan,
			item.CreatedBy,
			item.PenetapanId,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO tb_target_indikator_tujuan_penetapan_opd
		(
			indikator_tujuan_id,
			tahun,
			target,
			satuan,
			created_by,
			penetapan_id
		)
		VALUES %s
	`,
		strings.Join(valueStrings, ","),
	)

	result, err := tx.ExecContext(
		ctx,
		query,
		valueArgs...,
	)

	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}

func (r *PenetapanOpdRepository) SaveTujuanPenetapanOpdBatch(
	ctx context.Context,
	tx *sql.Tx,
	data []domain.TujuanPenetapanOpd,
) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}

	var (
		valueStrings []string
		valueArgs    []any
	)
	for i, item := range data {
		base := i * 7
		valueStrings = append(
			valueStrings,
			fmt.Sprintf(
				"($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				base+1,
				base+2,
				base+3,
				base+4,
				base+5,
				base+6,
				base+7,
			),
		)
		valueArgs = append(
			valueArgs,
			item.KodeOpd,
			item.KodeTujuanOpd,
			item.TujuanOpd,
			item.Periode,
			item.TahunAktif,
			item.CreatedBy,
			item.PenetapanId,
		)
	}
	query := fmt.Sprintf(`
		INSERT INTO tb_tujuan_penetapan_opd
		(
			kode_opd,
			kode_tujuan_opd,
			tujuan_opd,
			periode,
			tahun_aktif,
			created_by,
			penetapan_id
		)
		VALUES %s
	`,
		strings.Join(valueStrings, ","),
	)
	result, err := tx.ExecContext(
		ctx,
		query,
		valueArgs...,
	)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(rowsAffected), nil
}

func (r *PenetapanOpdRepository) GetActiveSnapshot(
	ctx context.Context,
	kodeOpd string,
	jenisPenetapan string,
	tahun int,
) (int64, error) {
	query := `SELECT id
                  FROM penetapan_opd
 		  WHERE kode_opd = $1
		    AND jenis_penetapan = $2
	            AND tahun = $3
                    AND is_active = TRUE
                   LIMIT 1
		`
	var id int64
	err := r.DB.QueryRowContext(
		ctx,
		query,
		kodeOpd,
		jenisPenetapan,
		tahun,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PenetapanOpdRepository) DeactivateOldSnapshot(
	ctx context.Context,
	tx *sql.Tx,
	kodeOpd string,
	jenisPenetapan string,
	tahun int,
) error {
	snapshotStatus := domain.SnapshotStatusArchived
	query := `UPDATE penetapan_opd
		  SET
		    snapshot_status = $4,
		    is_active = FALSE
 		  WHERE kode_opd = $1
		    AND jenis_penetapan = $2
	            AND tahun = $3
		`
	_, err := tx.ExecContext(
		ctx,
		query,
		kodeOpd,
		jenisPenetapan,
		tahun,
		snapshotStatus,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *PenetapanOpdRepository) GetPenetapanNextVersion(
	ctx context.Context,
	tx *sql.Tx,
	kodeOpd string,
	jenisPenetapan string,
	tahun int,
) (int, error) {
	query := `SELECT COALESCE(MAX(versi), 0) + 1
		  FROM penetapan_opd
 		  WHERE kode_opd = $1
		    AND jenis_penetapan = $2
	            AND tahun = $3
		`
	var versi int
	err := tx.QueryRowContext(
		ctx,
		query,
		kodeOpd,
		jenisPenetapan,
		tahun,
	).Scan(&versi)
	if err != nil {
		return 0, err
	}

	return versi, nil
}

func (r *PenetapanOpdRepository) SavePenetapanOpd(
	ctx context.Context,
	tx *sql.Tx,
	penetapan domain.PenetapanOpd,
) (int64, error) {
	query := `INSERT INTO
		penetapan_opd
		(kode_opd, tahun, jenis_penetapan, versi, snapshot_status, generated_by, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
		`
	var id int64

	err := tx.QueryRowContext(
		ctx,
		query,
		penetapan.KodeOpd,
		penetapan.Tahun,
		penetapan.JenisPenetapan,
		penetapan.Versi,
		penetapan.SnapshotStatus,
		penetapan.GeneratedBy,
		penetapan.IsActive,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PenetapanOpdRepository) InsertMetadata(
	ctx context.Context,
	metadata domain.SyncPenetapanMetadataOpd,
) (int64, error) {
	query := `INSERT INTO
		  sync_penetapan_metadata_opd
             	  (kode_opd, tahun, jenis_penetapan, status, started_at)
                  VALUES ($1, $2, $3, $4, $5)
	   	  RETURNING id`

	var id int64
	err := r.DB.QueryRowContext(
		ctx,
		query,
		metadata.KodeOpd,
		metadata.Tahun,
		metadata.JenisPenetapan,
		metadata.Status,
		metadata.StartedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PenetapanOpdRepository) UpdateMetadataStatus(
	ctx context.Context,
	metadataId int64,
	status string,
	errMessage *string,
) error {
	var query string
	if status == domain.SyncStatusSuccess ||
		status == domain.SyncStatusFailed {
		query = `
		UPDATE sync_penetapan_metadata_opd
		SET
			status = $1,
			error_message = $2,
			finished_at = NOW(),
			last_modified_date = NOW()
		WHERE id = $3
	`
	} else {
		query = `
		UPDATE sync_penetapan_metadata_opd
		SET
			status = $1,
			error_message = $2,
			last_modified_date = NOW()
		WHERE id = $3
	`
	}
	_, err := r.DB.ExecContext(
		ctx,
		query,
		status,
		errMessage,
		metadataId,
	)

	return err
}

func (r *PenetapanOpdRepository) FindTujuanBySnapshot(
	ctx context.Context,
	req domain.PenetapanOpdRequest,
) ([]domain.TujuanPenetapanOpd, error) {

	query := `
		SELECT
			tj.id,
			tj.kode_opd,
			tj.kode_tujuan_opd,
			tj.tujuan_opd,
			tj.periode,
			tj.tahun_aktif,
			tj.created_date,
			tj.last_modified_date,
			tj.created_by,
                        pn.versi
		FROM tb_tujuan_penetapan_opd tj
		JOIN penetapan_opd pn ON pn.id = tj.penetapan_id
		WHERE tj.kode_opd = $1
		AND tj.tahun_aktif = $2
                AND tj.penetapan_id = $3
		ORDER BY tj.id ASC
	`

	rows, err := r.DB.QueryContext(
		ctx,
		query,
		req.KodeOpd,
		req.Tahun,
		req.SnapshotId,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []domain.TujuanPenetapanOpd

	for rows.Next() {

		var item domain.TujuanPenetapanOpd

		err := rows.Scan(
			&item.Id,
			&item.KodeOpd,
			&item.KodeTujuanOpd,
			&item.TujuanOpd,
			&item.Periode,
			&item.TahunAktif,
			&item.CreatedDate,
			&item.LastModifiedDate,
			&item.CreatedBy,
			&item.Versi,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *PenetapanOpdRepository) FindSasaran(
	ctx context.Context,
	req domain.PenetapanOpdRequest,
) ([]domain.SasaranPenetapanOpd, error) {

	query := `
		SELECT
			id,
			kode_opd,
			kode_sasaran_opd,
			sasaran_opd,
			periode,
			tahun_aktif,
			created_date,
			last_modified_date,
			created_by
		FROM tb_sasaran_penetapan_opd
		WHERE kode_opd = $1
		AND tahun_aktif = $2
		ORDER BY id ASC
	`

	rows, err := r.DB.QueryContext(
		ctx,
		query,
		req.KodeOpd,
		req.Tahun,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make(
		[]domain.SasaranPenetapanOpd,
		0,
	)

	for rows.Next() {

		var item domain.SasaranPenetapanOpd

		err := rows.Scan(
			&item.Id,
			&item.KodeOpd,
			&item.KodeSasaranOpd,
			&item.SasaranOpd,
			&item.Periode,
			&item.TahunAktif,
			&item.CreatedDate,
			&item.LastModifiedDate,
			&item.CreatedBy,
		)

		if err != nil {
			return nil, err
		}

		result = append(
			result,
			item,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *PenetapanOpdRepository) FindIndikatorTujuanByTujuanIds(
	ctx context.Context,
	tujuanIds []int64,
) ([]domain.IndikatorTujuanPenetapanOpd, error) {
	if len(tujuanIds) == 0 {

		return []domain.IndikatorTujuanPenetapanOpd{}, nil
	}
	var (
		placeholders []string
		args         []any
	)
	for i, id := range tujuanIds {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		args = append(args, id)
	}
	query := fmt.Sprintf(`
		SELECT
			id,
			id_tujuan_opd,
			kode_indikator,
			kode_opd,
			indikator,
			rumus_perhitungan,
			sumber_data,
			definisi_operasional,
			tahun_aktif,
			created_date,
			last_modified_date,
			created_by,
			penetapan_id
		FROM tb_indikator_tujuan_penetapan_opd
		WHERE id_tujuan_opd IN (%s)
		ORDER BY id ASC
	`,
		strings.Join(placeholders, ","),
	)
	rows, err := r.DB.QueryContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(
		[]domain.IndikatorTujuanPenetapanOpd,
		0,
		len(tujuanIds),
	)
	for rows.Next() {
		var item domain.IndikatorTujuanPenetapanOpd
		err := rows.Scan(
			&item.Id,
			&item.IdTujuanOpd,
			&item.KodeIndikator,
			&item.KodeOpd,
			&item.Indikator,
			&item.RumusPerhitungan,
			&item.SumberData,
			&item.DefinisiOperasional,
			&item.TahunAktif,
			&item.CreatedDate,
			&item.LastModifiedDate,
			&item.CreatedBy,
			&item.PenetapanId,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
func (r *PenetapanOpdRepository) FindTargetIndikatorTujuanByIndikatorIds(
	ctx context.Context,
	indikatorIds []int64,
) ([]domain.TargetIndikatorTujuanPenetapanOpd, error) {

	if len(indikatorIds) == 0 {
		return []domain.TargetIndikatorTujuanPenetapanOpd{}, nil
	}

	var (
		placeholders []string
		args         []any
	)

	for i, id := range indikatorIds {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		args = append(args, id)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			indikator_tujuan_id,
			tahun,
			target,
			satuan,
			created_date,
			last_modified_date,
			created_by,
			penetapan_id
		FROM tb_target_indikator_tujuan_penetapan_opd
		WHERE indikator_tujuan_id IN (%s)
		ORDER BY id ASC
	`,
		strings.Join(placeholders, ","),
	)

	rows, err := r.DB.QueryContext(
		ctx,
		query,
		args...,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make(
		[]domain.TargetIndikatorTujuanPenetapanOpd,
		0,
		len(indikatorIds),
	)

	for rows.Next() {

		var item domain.TargetIndikatorTujuanPenetapanOpd

		err := rows.Scan(
			&item.Id,
			&item.IndikatorTujuanId,
			&item.Tahun,
			&item.Target,
			&item.Satuan,
			&item.CreatedDate,
			&item.LastModifiedDate,
			&item.CreatedBy,
			&item.PenetapanId,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
