package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/bappeda-dev-team/penetapan-service/internal/model/domain"
)

type PenetapanOpdRepository struct {
	DB *sql.DB
}

func NewPenetapanOpdRepository(db *sql.DB) *PenetapanOpdRepository {
	return &PenetapanOpdRepository{
		DB: db,
	}
}

// METADATA SNAPSHOT

func (r *PenetapanOpdRepository) GetActiveSnapshot(
	ctx context.Context,
	kodeOpd string,
	jenisPenetapan string,
	tahun int,
) (*domain.ActiveSnapshot, error) {
	query := `
	    SELECT id, versi
	    FROM penetapan_opd
	    WHERE
	    	kode_opd = $1
	    	AND jenis_penetapan = $2
	    	AND tahun = $3
	    	AND is_active = TRUE
	    ORDER BY versi DESC
	    LIMIT 1`

	var result domain.ActiveSnapshot
	err := r.DB.QueryRowContext(
		ctx,
		query,
		kodeOpd,
		jenisPenetapan,
		tahun,
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

// SASARAN OPD

func (r *PenetapanOpdRepository) SaveSasaranPenetapanOpd(
	ctx context.Context,
	tx *sql.Tx,
	data domain.SasaranPenetapanOpd,
) (int64, error) {

	query := `
		INSERT INTO sasaran_opd
		(
			kode_opd,
			kode_sasaran_opd,
			sasaran_opd,
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
		data.KodeSasaranOpd,
		data.SasaranOpd,
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

func (r *PenetapanOpdRepository) SaveIndikatorSasaranPenetapanOpd(
	ctx context.Context,
	tx *sql.Tx,
	data domain.IndikatorSasaranPenetapanOpd,
	sasaranId int64,
) (int64, error) {

	query := `
		INSERT INTO indikator_sasaran_opd
		(
			sasaran_opd_id,
			kode_opd,
			kode_indikator,
			indikator,
			rumus_perhitungan,
			sumber_data,
			definisi_operasional,
			tahun_aktif,
			created_by
		)
		VALUES
		(
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
		RETURNING id
	`

	var id int64

	err := tx.QueryRowContext(
		ctx,
		query,
		sasaranId,
		data.KodeOpd,
		data.KodeIndikator,
		data.Indikator,
		data.RumusPerhitungan,
		data.SumberData,
		data.DefinisiOperasional,
		data.TahunAktif,
		data.CreatedBy,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PenetapanOpdRepository) SaveTargetIndikatorSasaranBatch(
	ctx context.Context,
	tx *sql.Tx,
	indikatorId int64,
	data []domain.TargetIndikatorSasaranPenetapanOpd,
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
			item.KodeTarget,
			item.Tahun,
			item.Target,
			item.Satuan,
			item.CreatedBy,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO target_indikator_sasaran_opd
		(
			indikator_sasaran_id,
			kode_target,
			tahun,
			target,
			satuan,
			created_by
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

// TUJUAN OPD

func (r *PenetapanOpdRepository) SaveTujuanPenetapanOpd(
	ctx context.Context,
	tx *sql.Tx,
	data domain.TujuanPenetapanOpd,
) (int64, error) {

	query := `
		INSERT INTO tujuan_opd
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
		INSERT INTO indikator_tujuan_opd
		(
			tujuan_opd_id,
			kode_opd,
			kode_indikator,
			indikator,
			rumus_perhitungan,
			sumber_data,
			definisi_operasional,
			tahun_aktif,
			created_by
		)
		VALUES
		(
			$1, $2, $3, $4, $5, $6, $7, $8, $9
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
			item.KodeTarget,
			item.Tahun,
			item.Target,
			item.Satuan,
			item.CreatedBy,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO target_indikator_tujuan_opd
		(
			indikator_tujuan_id,
			kode_target,
			tahun,
			target,
			satuan,
			created_by
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
		FROM tujuan_opd tj
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
			tujuan_opd_id,
			kode_indikator,
			kode_opd,
			indikator,
			rumus_perhitungan,
			sumber_data,
			definisi_operasional,
			tahun_aktif,
			created_date,
			last_modified_date,
			created_by
		FROM indikator_tujuan_opd
		WHERE tujuan_opd_id IN (%s)
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
			kode_target,
			tahun,
			target,
			satuan,
			created_date,
			last_modified_date,
			created_by
		FROM target_indikator_tujuan_opd
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
			&item.KodeTarget,
			&item.Tahun,
			&item.Target,
			&item.Satuan,
			&item.CreatedDate,
			&item.LastModifiedDate,
			&item.CreatedBy,
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

func (r *PenetapanOpdRepository) FindSasaranBySnapshot(
	ctx context.Context,
	req domain.PenetapanOpdRequest,
) ([]domain.SasaranPenetapanOpd, error) {

	query := `
		SELECT
			sas.id,
			sas.kode_opd,
			sas.kode_sasaran_opd,
			sas.sasaran_opd,
			sas.periode,
			sas.tahun_aktif,
			sas.created_date,
			sas.last_modified_date,
			sas.created_by,
                        pn.versi
		FROM sasaran_opd sas
		JOIN penetapan_opd pn ON pn.id = sas.penetapan_id
		WHERE sas.kode_opd = $1
		AND sas.tahun_aktif = $2
                AND sas.penetapan_id = $3
		ORDER BY sas.id ASC
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

	var result []domain.SasaranPenetapanOpd

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

func (r *PenetapanOpdRepository) FindIndikatorSasaranBySasarands(
	ctx context.Context,
	sasaranids []int64,
) ([]domain.IndikatorSasaranPenetapanOpd, error) {
	if len(sasaranids) == 0 {

		return []domain.IndikatorSasaranPenetapanOpd{}, nil
	}
	var (
		placeholders []string
		args         []any
	)
	for i, id := range sasaranids {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		args = append(args, id)
	}
	query := fmt.Sprintf(`
		SELECT
			id,
			sasaran_opd_id,
			kode_indikator,
			kode_opd,
			indikator,
			rumus_perhitungan,
			sumber_data,
			definisi_operasional,
			tahun_aktif,
			created_date,
			last_modified_date,
			created_by
		FROM indikator_sasaran_opd
		WHERE sasaran_opd_id IN (%s)
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
		[]domain.IndikatorSasaranPenetapanOpd,
		0,
		len(sasaranids),
	)
	for rows.Next() {
		var item domain.IndikatorSasaranPenetapanOpd
		err := rows.Scan(
			&item.Id,
			&item.IdSasaranOpd,
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

func (r *PenetapanOpdRepository) FindTargetIndikatorSasaranByIndikatorIds(
	ctx context.Context,
	indikatorIds []int64,
) ([]domain.TargetIndikatorSasaranPenetapanOpd, error) {

	if len(indikatorIds) == 0 {
		return []domain.TargetIndikatorSasaranPenetapanOpd{}, nil
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
			indikator_sasaran_id,
			kode_target,
			tahun,
			target,
			satuan,
			created_date,
			last_modified_date,
			created_by
		FROM target_indikator_sasaran_opd
		WHERE indikator_sasaran_id IN (%s)
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
		[]domain.TargetIndikatorSasaranPenetapanOpd,
		0,
		len(indikatorIds),
	)

	for rows.Next() {

		var item domain.TargetIndikatorSasaranPenetapanOpd

		err := rows.Scan(
			&item.Id,
			&item.IndikatorSasaranId,
			&item.KodeTarget,
			&item.Tahun,
			&item.Target,
			&item.Satuan,
			&item.CreatedDate,
			&item.LastModifiedDate,
			&item.CreatedBy,
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

func (r *PenetapanOpdRepository) SaveRenjaUrusan(
	ctx context.Context,
	tx *sql.Tx,
	urusan domain.RenjaUrusan,
) (int64, error) {

	query := `
	INSERT INTO renja_urusan(
		penetapan_id,
		kode_opd,
		kode_urusan,
		urusan,
		tahun_aktif,
		created_by
	)
	VALUES($1,$2,$3,$4,$5,$6)
	RETURNING id
	`

	var id int64
	err := tx.QueryRowContext(
		ctx,
		query,
		urusan.PenetapanId,
		urusan.KodeOpd,
		urusan.KodeUrusan,
		urusan.Urusan,
		urusan.TahunAktif,
		urusan.CreatedBy,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PenetapanOpdRepository) SaveRenjaBidangUrusan(
	ctx context.Context,
	tx *sql.Tx,
	bidang domain.RenjaBidangUrusan,
) (int64, error) {

	query := `
	INSERT INTO renja_bidang_urusan(
		penetapan_id,
		kode_opd,
		kode_urusan,
		kode_bidang_urusan,
		bidang_urusan,
		tahun_aktif,
		created_by
	)
	VALUES($1,$2,$3,$4,$5,$6,$7)
	RETURNING id
	`

	var id int64
	err := tx.QueryRowContext(
		ctx,
		query,
		bidang.PenetapanId,
		bidang.KodeOpd,
		bidang.KodeUrusan,
		bidang.KodeBidangUrusan,
		bidang.BidangUrusan,
		bidang.TahunAktif,
		bidang.CreatedBy,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (r *PenetapanOpdRepository) SaveRenjaProgram(
	ctx context.Context,
	tx *sql.Tx,
	program domain.RenjaProgram,
) (int64, error) {

	query := `
	INSERT INTO renja_program(
		penetapan_id,
		kode_opd,
		kode_bidang_urusan,
		kode_program,
		program,
		tahun_aktif,
		created_by
	)
	VALUES($1,$2,$3,$4,$5,$6,$7)
	RETURNING id
	`

	var id int64

	err := tx.QueryRowContext(
		ctx,
		query,
		program.PenetapanId,
		program.KodeOpd,
		program.KodeBidangUrusan,
		program.KodeProgram,
		program.Program,
		program.TahunAktif,
		program.CreatedBy,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PenetapanOpdRepository) SaveRenjaKegiatan(
	ctx context.Context,
	tx *sql.Tx,
	kegiatan domain.RenjaKegiatan,
) (int64, error) {

	query := `
	INSERT INTO renja_kegiatan(
		penetapan_id,
		kode_opd,
		kode_program,
		kode_kegiatan,
		kegiatan,
		tahun_aktif,
		created_by
	)
	VALUES($1,$2,$3,$4,$5,$6,$7)
	RETURNING id
	`

	var id int64

	err := tx.QueryRowContext(
		ctx,
		query,
		kegiatan.PenetapanId,
		kegiatan.KodeOpd,
		kegiatan.KodeProgram,
		kegiatan.KodeKegiatan,
		kegiatan.Kegiatan,
		kegiatan.TahunAktif,
		kegiatan.CreatedBy,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PenetapanOpdRepository) SaveRenjaSubkegiatan(
	ctx context.Context,
	tx *sql.Tx,
	sub domain.RenjaSubkegiatan,
) (int64, error) {

	query := `
	INSERT INTO renja_subkegiatan(
		penetapan_id,
		kode_opd,
		kode_kegiatan,
		kode_subkegiatan,
		subkegiatan,
		tahun_aktif,
		created_by
	)
	VALUES($1,$2,$3,$4,$5,$6,$7)
	RETURNING id
	`

	var id int64

	err := tx.QueryRowContext(
		ctx,
		query,
		sub.PenetapanId,
		sub.KodeOpd,
		sub.KodeKegiatan,
		sub.KodeSubkegiatan,
		sub.Subkegiatan,
		sub.TahunAktif,
		sub.CreatedBy,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PenetapanOpdRepository) SaveIndikatorRenjaProgram(
	ctx context.Context,
	tx *sql.Tx,
	indikator domain.IndikatorRenjaProgram,
	prgId int64,
) (int64, error) {

	query := `
	INSERT INTO indikator_renja_program(
		program_id,
		kode_indikator,
		indikator,
		tahun,
		created_by
	)
	VALUES($1,$2,$3,$4,$5)
	RETURNING id
	`

	var id int64

	err := tx.QueryRowContext(
		ctx,
		query,
		prgId,
		indikator.KodeIndikator,
		indikator.Indikator,
		indikator.Tahun,
		indikator.CreatedBy,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PenetapanOpdRepository) SaveTargetIndikatorRenjaProgramBatch(
	ctx context.Context,
	tx *sql.Tx,
	indId int64,
	targets []domain.TargetIndikatorRenjaProgram,
) (int, error) {

	if len(targets) == 0 {
		return 0, nil
	}

	query := `
	INSERT INTO target_indikator_renja_program(
		indikator_program_id,
		kode_target,
		target,
		satuan,
		tahun,
		created_by
	)
	VALUES
	`

	args := []any{}
	values := []string{}

	for i, t := range targets {

		offset := i*6 + 1

		values = append(
			values,
			fmt.Sprintf(
				"($%d,$%d,$%d,$%d,$%d,$%d)",
				offset,
				offset+1,
				offset+2,
				offset+3,
				offset+4,
				offset+5,
			),
		)

		args = append(
			args,
			indId,
			t.KodeTarget,
			t.Target,
			t.Satuan,
			t.Tahun,
			t.CreatedBy,
		)
	}

	query += strings.Join(values, ",")

	result, err := tx.ExecContext(
		ctx,
		query,
		args...,
	)

	if err != nil {
		return 0, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rows), nil
}

func (r *PenetapanOpdRepository) SaveIndikatorRenjaKegiatan(
	ctx context.Context,
	tx *sql.Tx,
	indikator domain.IndikatorRenjaKegiatan,
	kegId int64,
) (int64, error) {

	query := `
	INSERT INTO indikator_renja_kegiatan(
		kegiatan_id,
		kode_indikator,
		indikator,
		tahun,
		created_by
	)
	VALUES($1,$2,$3,$4,$5)
	RETURNING id
	`

	var id int64

	err := tx.QueryRowContext(
		ctx,
		query,
		kegId,
		indikator.KodeIndikator,
		indikator.Indikator,
		indikator.Tahun,
		indikator.CreatedBy,
	).Scan(&id)

	return id, err
}

func (r *PenetapanOpdRepository) SaveTargetIndikatorRenjaKegiatanBatch(
	ctx context.Context,
	tx *sql.Tx,
	indId int64,
	targets []domain.TargetIndikatorRenjaKegiatan,
) (int, error) {

	if len(targets) == 0 {
		return 0, nil
	}

	query := `
	INSERT INTO target_indikator_renja_kegiatan(
		indikator_kegiatan_id,
		kode_target,
		target,
		satuan,
		tahun,
		created_by
	)
	VALUES
	`

	var (
		args   []any
		values []string
	)

	for i, target := range targets {

		offset := (i * 6) + 1

		values = append(
			values,
			fmt.Sprintf(
				"($%d,$%d,$%d,$%d,$%d,$%d)",
				offset,
				offset+1,
				offset+2,
				offset+3,
				offset+4,
				offset+5,
			),
		)

		args = append(
			args,
			indId,
			target.KodeTarget,
			target.Target,
			target.Satuan,
			target.Tahun,
			target.CreatedBy,
		)
	}

	query += strings.Join(values, ",")

	result, err := tx.ExecContext(
		ctx,
		query,
		args...,
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

func (r *PenetapanOpdRepository) SaveIndikatorRenjaSubkegiatan(
	ctx context.Context,
	tx *sql.Tx,
	indikator domain.IndikatorRenjaSubkegiatan,
	subId int64,
) (int64, error) {

	query := `
	INSERT INTO indikator_renja_subkegiatan(
		subkegiatan_id,
		kode_indikator,
		indikator,
		tahun,
		created_by
	)
	VALUES($1,$2,$3,$4,$5)
	RETURNING id
	`

	var id int64

	err := tx.QueryRowContext(
		ctx,
		query,
		subId,
		indikator.KodeIndikator,
		indikator.Indikator,
		indikator.Tahun,
		indikator.CreatedBy,
	).Scan(&id)

	return id, err
}

func (r *PenetapanOpdRepository) SaveTargetIndikatorRenjaSubkegiatanBatch(
	ctx context.Context,
	tx *sql.Tx,
	indId int64,
	targets []domain.TargetIndikatorRenjaSubkegiatan,
) (int, error) {

	if len(targets) == 0 {
		return 0, nil
	}

	query := `
	INSERT INTO target_indikator_renja_subkegiatan(
		indikator_subkegiatan_id,
		kode_target,
		target,
		satuan,
		tahun,
		created_by
	)
	VALUES
	`

	var (
		args   []any
		values []string
	)

	for i, target := range targets {

		offset := (i * 6) + 1

		values = append(
			values,
			fmt.Sprintf(
				"($%d,$%d,$%d,$%d,$%d,$%d)",
				offset,
				offset+1,
				offset+2,
				offset+3,
				offset+4,
				offset+5,
			),
		)

		args = append(
			args,
			indId,
			target.KodeTarget,
			target.Target,
			target.Satuan,
			target.Tahun,
			target.CreatedBy,
		)
	}

	query += strings.Join(values, ",")

	result, err := tx.ExecContext(
		ctx,
		query,
		args...,
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

func (r *PenetapanOpdRepository) FindRenjaUrusanBySnapshot(
	ctx context.Context,
	req domain.PenetapanOpdRequest,
) ([]domain.RenjaUrusan, error) {

	query := `
	SELECT
		id,
		penetapan_id,
		kode_opd,
		kode_urusan,
		urusan,
		tahun_aktif,
		created_date,
		last_modified_date,
		created_by
	FROM renja_urusan
	WHERE penetapan_id=$1
	ORDER BY kode_urusan
	`

	rows, err := r.DB.QueryContext(
		ctx,
		query,
		req.SnapshotId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.RenjaUrusan

	for rows.Next() {

		var data domain.RenjaUrusan

		err := rows.Scan(
			&data.Id,
			&data.PenetapanId,
			&data.KodeOpd,
			&data.KodeUrusan,
			&data.Urusan,
			&data.TahunAktif,
			&data.CreatedDate,
			&data.LastModifiedDate,
			&data.CreatedBy,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, data)
	}

	return result, rows.Err()
}

func (r *PenetapanOpdRepository) FindRenjaBidangUrusanBySnapshot(
	ctx context.Context,
	req domain.PenetapanOpdRequest,
) ([]domain.RenjaBidangUrusan, error) {

	query := `
	SELECT
		id,
		penetapan_id,
		kode_opd,
		kode_urusan,
		kode_bidang_urusan,
		bidang_urusan,
		tahun_aktif,
		created_date,
		last_modified_date,
		created_by
	FROM renja_bidang_urusan
	WHERE penetapan_id=$1
	ORDER BY kode_bidang_urusan
	`

	rows, err := r.DB.QueryContext(ctx, query, req.SnapshotId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.RenjaBidangUrusan

	for rows.Next() {

		var data domain.RenjaBidangUrusan

		err := rows.Scan(
			&data.Id,
			&data.PenetapanId,
			&data.KodeOpd,
			&data.KodeUrusan,
			&data.KodeBidangUrusan,
			&data.BidangUrusan,
			&data.TahunAktif,
			&data.CreatedDate,
			&data.LastModifiedDate,
			&data.CreatedBy,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, data)
	}

	return result, rows.Err()
}

func (r *PenetapanOpdRepository) FindRenjaProgramBySnapshot(
	ctx context.Context,
	req domain.PenetapanOpdRequest,
) ([]domain.RenjaProgram, error) {

	query := `
	SELECT
		id,
		penetapan_id,
		kode_opd,
		kode_bidang_urusan,
		kode_program,
		program,
		tahun_aktif,
		created_date,
		last_modified_date,
		created_by
	FROM renja_program
	WHERE penetapan_id=$1
	ORDER BY kode_program
	`

	rows, err := r.DB.QueryContext(
		ctx,
		query,
		req.SnapshotId,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []domain.RenjaProgram

	for rows.Next() {

		var data domain.RenjaProgram

		err := rows.Scan(
			&data.Id,
			&data.PenetapanId,
			&data.KodeOpd,
			&data.KodeBidangUrusan,
			&data.KodeProgram,
			&data.Program,
			&data.TahunAktif,
			&data.CreatedDate,
			&data.LastModifiedDate,
			&data.CreatedBy,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, data)
	}

	return result, rows.Err()
}

func (r *PenetapanOpdRepository) FindRenjaKegiatanBySnapshot(
	ctx context.Context,
	req domain.PenetapanOpdRequest,
) ([]domain.RenjaKegiatan, error) {

	query := `
	SELECT
		id,
		penetapan_id,
		kode_opd,
		kode_program,
		kode_kegiatan,
		kegiatan,
		tahun_aktif,
		created_date,
		last_modified_date,
		created_by
	FROM renja_kegiatan
	WHERE penetapan_id=$1
	ORDER BY kode_kegiatan
	`

	rows, err := r.DB.QueryContext(
		ctx,
		query,
		req.SnapshotId,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []domain.RenjaKegiatan

	for rows.Next() {

		var data domain.RenjaKegiatan

		err := rows.Scan(
			&data.Id,
			&data.PenetapanId,
			&data.KodeOpd,
			&data.KodeProgram,
			&data.KodeKegiatan,
			&data.Kegiatan,
			&data.TahunAktif,
			&data.CreatedDate,
			&data.LastModifiedDate,
			&data.CreatedBy,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, data)
	}

	return result, rows.Err()
}

func (r *PenetapanOpdRepository) FindRenjaSubkegiatanBySnapshot(
	ctx context.Context,
	req domain.PenetapanOpdRequest,
) ([]domain.RenjaSubkegiatan, error) {

	query := `
	SELECT
		id,
		penetapan_id,
		kode_opd,
		kode_kegiatan,
		kode_subkegiatan,
		subkegiatan,
		tahun_aktif,
		created_date,
		last_modified_date,
		created_by
	FROM renja_subkegiatan
	WHERE penetapan_id=$1
	ORDER BY kode_subkegiatan
	`

	rows, err := r.DB.QueryContext(
		ctx,
		query,
		req.SnapshotId,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []domain.RenjaSubkegiatan

	for rows.Next() {

		var data domain.RenjaSubkegiatan

		err := rows.Scan(
			&data.Id,
			&data.PenetapanId,
			&data.KodeOpd,
			&data.KodeKegiatan,
			&data.KodeSubkegiatan,
			&data.Subkegiatan,
			&data.TahunAktif,
			&data.CreatedDate,
			&data.LastModifiedDate,
			&data.CreatedBy,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, data)
	}

	return result, rows.Err()
}

func (r *PenetapanOpdRepository) FindIndikatorRenjaProgram(
	ctx context.Context,
	programIds []int64,
) ([]domain.IndikatorRenjaProgram, error) {

	if len(programIds) == 0 {
		return []domain.IndikatorRenjaProgram{}, nil
	}

	var (
		placeholders []string
		args         []any
	)

	for i, id := range programIds {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		args = append(args, id)
	}

	query := fmt.Sprintf(`
	SELECT
		id,
		program_id,
		kode_indikator,
		indikator,
		tahun,
		created_date,
		last_modified_date,
		created_by
	FROM indikator_renja_program
	WHERE program_id IN (%s)
	ORDER BY id
	`, strings.Join(placeholders, ","))

	rows, err := r.DB.QueryContext(
		ctx,
		query,
		args...,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []domain.IndikatorRenjaProgram

	for rows.Next() {

		var item domain.IndikatorRenjaProgram

		err := rows.Scan(
			&item.Id,
			&item.ProgramId,
			&item.KodeIndikator,
			&item.Indikator,
			&item.Tahun,
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

	return result, rows.Err()
}

func (r *PenetapanOpdRepository) FindTargetIndikatorRenjaProgramBatch(
	ctx context.Context,
	indikatorIds []int64,
) ([]domain.TargetIndikatorRenjaProgram, error) {

	if len(indikatorIds) == 0 {
		return []domain.TargetIndikatorRenjaProgram{}, nil
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
		indikator_program_id,
		kode_target,
		target,
		satuan,
		tahun,
		created_date,
		last_modified_date,
		created_by
	FROM target_indikator_renja_program
	WHERE indikator_program_id IN (%s)
	ORDER BY id
	`, strings.Join(placeholders, ","))

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
		[]domain.TargetIndikatorRenjaProgram,
		0,
		len(indikatorIds),
	)

	for rows.Next() {

		var item domain.TargetIndikatorRenjaProgram

		err := rows.Scan(
			&item.Id,
			&item.IndikatorProgramId,
			&item.KodeTarget,
			&item.Target,
			&item.Satuan,
			&item.Tahun,
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

	return result, rows.Err()
}

func (r *PenetapanOpdRepository) FindIndikatorRenjaKegiatan(
	ctx context.Context,
	kegiatanIds []int64,
) ([]domain.IndikatorRenjaKegiatan, error) {

	if len(kegiatanIds) == 0 {
		return []domain.IndikatorRenjaKegiatan{}, nil
	}

	var (
		placeholders []string
		args         []any
	)

	for i, id := range kegiatanIds {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		args = append(args, id)
	}

	query := fmt.Sprintf(`
	SELECT
		id,
		kegiatan_id,
		kode_indikator,
		indikator,
		tahun,
		created_date,
		last_modified_date,
		created_by
	FROM indikator_renja_kegiatan
	WHERE kegiatan_id IN (%s)
	ORDER BY id
	`, strings.Join(placeholders, ","))

	rows, err := r.DB.QueryContext(
		ctx,
		query,
		args...,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []domain.IndikatorRenjaKegiatan

	for rows.Next() {

		var item domain.IndikatorRenjaKegiatan

		err := rows.Scan(
			&item.Id,
			&item.KegiatanId,
			&item.KodeIndikator,
			&item.Indikator,
			&item.Tahun,
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

	return result, rows.Err()
}

func (r *PenetapanOpdRepository) FindTargetIndikatorRenjaKegiatanBatch(
	ctx context.Context,
	indikatorIds []int64,
) ([]domain.TargetIndikatorRenjaKegiatan, error) {

	if len(indikatorIds) == 0 {
		return []domain.TargetIndikatorRenjaKegiatan{}, nil
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
		indikator_kegiatan_id,
		kode_target,
		target,
		satuan,
		tahun,
		created_date,
		last_modified_date,
		created_by
	FROM target_indikator_renja_kegiatan
	WHERE indikator_kegiatan_id IN (%s)
	ORDER BY id
	`, strings.Join(placeholders, ","))

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
		[]domain.TargetIndikatorRenjaKegiatan,
		0,
		len(indikatorIds),
	)

	for rows.Next() {

		var item domain.TargetIndikatorRenjaKegiatan

		err := rows.Scan(
			&item.Id,
			&item.IndikatorKegiatanId,
			&item.KodeTarget,
			&item.Target,
			&item.Satuan,
			&item.Tahun,
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

	return result, rows.Err()
}

func (r *PenetapanOpdRepository) FindIndikatorRenjaSubkegiatan(
	ctx context.Context,
	subkegiatanIds []int64,
) ([]domain.IndikatorRenjaSubkegiatan, error) {

	if len(subkegiatanIds) == 0 {

		return []domain.IndikatorRenjaSubkegiatan{}, nil
	}
	var (
		placeholders []string
		args         []any
	)
	for i, id := range subkegiatanIds {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		args = append(args, id)
	}

	query := fmt.Sprintf(`
	SELECT
		id,
		kegiatan_id,
		kode_indikator,
		indikator,
		tahun,
		created_date,
		last_modified_date,
		created_by
	FROM indikator_renja_kegiatan
	WHERE kegiatan_id IN (%s)
	ORDER BY id
	`, strings.Join(placeholders, ","))

	rows, err := r.DB.QueryContext(
		ctx,
		query,
		args...,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []domain.IndikatorRenjaSubkegiatan

	for rows.Next() {

		var item domain.IndikatorRenjaSubkegiatan

		err := rows.Scan(
			&item.Id,
			&item.SubkegiatanId,
			&item.KodeIndikator,
			&item.Indikator,
			&item.Tahun,
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

	return result, rows.Err()
}

func (r *PenetapanOpdRepository) FindTargetIndikatorRenjaSubkegiatanBatch(
	ctx context.Context,
	indikatorIds []int64,
) ([]domain.TargetIndikatorRenjaSubkegiatan, error) {

	if len(indikatorIds) == 0 {
		return []domain.TargetIndikatorRenjaSubkegiatan{}, nil
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
		indikator_subkegiatan_id,
		kode_target,
		target,
		satuan,
		tahun,
		created_date,
		last_modified_date,
		created_by
	FROM target_indikator_renja_subkegiatan
	WHERE indikator_subkegiatan_id IN (%s)
	ORDER BY id
	`, strings.Join(placeholders, ","))

	rows, err := r.DB.QueryContext(
		ctx,
		query,
		args...,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []domain.TargetIndikatorRenjaSubkegiatan

	for rows.Next() {

		var item domain.TargetIndikatorRenjaSubkegiatan

		err := rows.Scan(
			&item.Id,
			&item.IndikatorSubkegiatanId,
			&item.KodeTarget,
			&item.Target,
			&item.Satuan,
			&item.Tahun,
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

	return result, rows.Err()
}

func (r *PenetapanOpdRepository) SavePaguRenjaUrusan(
	ctx context.Context,
	tx *sql.Tx,
	urusanId int64,
	data []domain.AnggaranRenja,
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
			urusanId,
			item.KodePagu,
			item.Tahun,
			item.PaguAnggaran,
			item.JenisPagu,
			item.CreatedBy,
		)
	}

	query := fmt.Sprintf(`
	INSERT INTO pagu_renja_urusan
	(
		urusan_id,
		kode_pagu,
		tahun,
		pagu,
		jenis_pagu,
		created_by
	)
	VALUES %s
	`, strings.Join(valueStrings, ","),
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

func (r *PenetapanOpdRepository) SavePaguRenjaBidangUrusan(
	ctx context.Context,
	tx *sql.Tx,
	bidangUrusanId int64,
	data []domain.AnggaranRenja,
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
			bidangUrusanId,
			item.KodePagu,
			item.Tahun,
			item.PaguAnggaran,
			item.JenisPagu,
			item.CreatedBy,
		)
	}

	query := fmt.Sprintf(`
	INSERT INTO pagu_renja_bidang_urusan
	(
		bidang_urusan_id,
		kode_pagu,
		tahun,
		pagu,
		jenis_pagu,
		created_by
	)
	VALUES %s
	`, strings.Join(valueStrings, ","),
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

func (r *PenetapanOpdRepository) SavePaguRenjaProgram(
	ctx context.Context,
	tx *sql.Tx,
	programId int64,
	data []domain.AnggaranRenja,
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
			programId,
			item.KodePagu,
			item.Tahun,
			item.PaguAnggaran,
			item.JenisPagu,
			item.CreatedBy,
		)
	}

	query := fmt.Sprintf(`
	INSERT INTO pagu_renja_program
	(
		program_id,
		kode_pagu,
		tahun,
		pagu,
		jenis_pagu,
		created_by
	)
	VALUES %s
	`, strings.Join(valueStrings, ","),
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

func (r *PenetapanOpdRepository) SavePaguRenjaKegiatan(
	ctx context.Context,
	tx *sql.Tx,
	kegiatanId int64,
	data []domain.AnggaranRenja,
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
			kegiatanId,
			item.KodePagu,
			item.Tahun,
			item.PaguAnggaran,
			item.JenisPagu,
			item.CreatedBy,
		)
	}

	query := fmt.Sprintf(`
	INSERT INTO pagu_renja_kegiatan
	(
		kegiatan_id,
		kode_pagu,
		tahun,
		pagu,
		jenis_pagu,
		created_by
	)
	VALUES %s
	`, strings.Join(valueStrings, ","),
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

func (r *PenetapanOpdRepository) SavePaguRenjaSubkegiatan(
	ctx context.Context,
	tx *sql.Tx,
	subkegiatanId int64,
	data []domain.AnggaranRenja,
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
			subkegiatanId,
			item.KodePagu,
			item.Tahun,
			item.PaguAnggaran,
			item.JenisPagu,
			item.CreatedBy,
		)
	}

	query := fmt.Sprintf(`
	INSERT INTO pagu_renja_subkegiatan
	(
		subkegiatan_id,
		kode_pagu,
		tahun,
		pagu,
		jenis_pagu,
		created_by
	)
	VALUES %s
	`, strings.Join(valueStrings, ","),
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

func (r *PenetapanOpdRepository) FindPaguRenjaUrusan(
	ctx context.Context,
	urusanIds []int64,
) ([]domain.AnggaranRenja, error) {

	if len(urusanIds) == 0 {
		return []domain.AnggaranRenja{}, nil
	}

	var (
		placeholders []string
		args         []any
	)

	for i, id := range urusanIds {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		args = append(args, id)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			urusan_id,
			kode_pagu,
			tahun,
			pagu,
			jenis_pagu,
			created_date,
			last_modified_date,
			created_by
		FROM pagu_renja_urusan
		WHERE urusan_id IN (%s)
		ORDER BY urusan_id,tahun
	`,
		strings.Join(
			placeholders,
			",",
		),
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
	var pagus []domain.AnggaranRenja
	for rows.Next() {
		var pagu domain.AnggaranRenja
		err := rows.Scan(
			&pagu.Id,
			&pagu.UrusanId,
			&pagu.KodePagu,
			&pagu.Tahun,
			&pagu.PaguAnggaran,
			&pagu.JenisPagu,
			&pagu.CreatedDate,
			&pagu.LastModifiedDate,
			&pagu.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		pagus = append(
			pagus,
			pagu,
		)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return pagus, nil
}

func (r *PenetapanOpdRepository) FindPaguRenjaBidangUrusan(
	ctx context.Context,
	bidangIds []int64,
) ([]domain.AnggaranRenja, error) {

	if len(bidangIds) == 0 {
		return []domain.AnggaranRenja{}, nil
	}

	var (
		placeholders []string
		args         []any
	)

	for i, id := range bidangIds {
		placeholders = append(
			placeholders,
			fmt.Sprintf("$%d", i+1),
		)
		args = append(args, id)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			bidang_urusan_id,
			kode_pagu,
			tahun,
			pagu,
			jenis_pagu,
			created_date,
			last_modified_date,
			created_by
		FROM pagu_renja_bidang_urusan
		WHERE bidang_urusan_id IN (%s)
		ORDER BY bidang_urusan_id,tahun
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

	var pagus []domain.AnggaranRenja

	for rows.Next() {

		var pagu domain.AnggaranRenja

		err := rows.Scan(
			&pagu.Id,
			&pagu.BidangUrusanId,
			&pagu.KodePagu,
			&pagu.Tahun,
			&pagu.PaguAnggaran,
			&pagu.JenisPagu,
			&pagu.CreatedDate,
			&pagu.LastModifiedDate,
			&pagu.CreatedBy,
		)

		if err != nil {
			return nil, err
		}

		pagus = append(
			pagus,
			pagu,
		)
	}

	return pagus, rows.Err()
}

func (r *PenetapanOpdRepository) FindPaguRenjaProgram(
	ctx context.Context,
	programIds []int64,
) ([]domain.AnggaranRenja, error) {

	if len(programIds) == 0 {
		return []domain.AnggaranRenja{}, nil
	}

	var (
		placeholders []string
		args         []any
	)

	for i, id := range programIds {
		placeholders = append(
			placeholders,
			fmt.Sprintf("$%d", i+1),
		)
		args = append(args, id)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			program_id,
			kode_pagu,
			tahun,
			pagu,
			jenis_pagu,
			created_date,
			last_modified_date,
			created_by
		FROM pagu_renja_program
		WHERE program_id IN (%s)
		ORDER BY program_id,tahun
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

	var pagus []domain.AnggaranRenja

	for rows.Next() {

		var pagu domain.AnggaranRenja

		err := rows.Scan(
			&pagu.Id,
			&pagu.ProgramId,
			&pagu.KodePagu,
			&pagu.Tahun,
			&pagu.PaguAnggaran,
			&pagu.JenisPagu,
			&pagu.CreatedDate,
			&pagu.LastModifiedDate,
			&pagu.CreatedBy,
		)

		if err != nil {
			return nil, err
		}

		pagus = append(
			pagus,
			pagu,
		)
	}

	return pagus, rows.Err()
}

func (r *PenetapanOpdRepository) FindPaguRenjaKegiatan(
	ctx context.Context,
	kegiatanIds []int64,
) ([]domain.AnggaranRenja, error) {

	if len(kegiatanIds) == 0 {
		return []domain.AnggaranRenja{}, nil
	}

	var (
		placeholders []string
		args         []any
	)

	for i, id := range kegiatanIds {
		placeholders = append(
			placeholders,
			fmt.Sprintf("$%d", i+1),
		)
		args = append(args, id)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			kegiatan_id,
			kode_pagu,
			tahun,
			pagu,
			jenis_pagu,
			created_date,
			last_modified_date,
			created_by
		FROM pagu_renja_kegiatan
		WHERE kegiatan_id IN (%s)
		ORDER BY kegiatan_id,tahun
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

	var pagus []domain.AnggaranRenja

	for rows.Next() {

		var pagu domain.AnggaranRenja

		err := rows.Scan(
			&pagu.Id,
			&pagu.KegiatanId,
			&pagu.KodePagu,
			&pagu.Tahun,
			&pagu.PaguAnggaran,
			&pagu.JenisPagu,
			&pagu.CreatedDate,
			&pagu.LastModifiedDate,
			&pagu.CreatedBy,
		)

		if err != nil {
			return nil, err
		}

		pagus = append(
			pagus,
			pagu,
		)
	}

	return pagus, rows.Err()
}

func (r *PenetapanOpdRepository) FindPaguRenjaSubkegiatan(
	ctx context.Context,
	subIds []int64,
) ([]domain.AnggaranRenja, error) {

	if len(subIds) == 0 {
		return []domain.AnggaranRenja{}, nil
	}

	var (
		placeholders []string
		args         []any
	)

	for i, id := range subIds {
		placeholders = append(
			placeholders,
			fmt.Sprintf("$%d", i+1),
		)
		args = append(args, id)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			subkegiatan_id,
			kode_pagu,
			tahun,
			pagu,
			jenis_pagu,
			created_date,
			last_modified_date,
			created_by
		FROM pagu_renja_subkegiatan
		WHERE subkegiatan_id IN (%s)
		ORDER BY subkegiatan_id,tahun
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

	var pagus []domain.AnggaranRenja

	for rows.Next() {

		var pagu domain.AnggaranRenja

		err := rows.Scan(
			&pagu.Id,
			&pagu.SubkegiatanId,
			&pagu.KodePagu,
			&pagu.Tahun,
			&pagu.PaguAnggaran,
			&pagu.JenisPagu,
			&pagu.CreatedDate,
			&pagu.LastModifiedDate,
			&pagu.CreatedBy,
		)

		if err != nil {
			return nil, err
		}

		pagus = append(
			pagus,
			pagu,
		)
	}

	return pagus, rows.Err()
}
