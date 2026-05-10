package api

import (
	"context"
	"database/sql"

	"github.com/bappeda-dev-team/penetapan-service/internal/model/domain"
)

type PenetapanOpdRepository struct {
	DB *sql.DB
}

func (r *PenetapanOpdRepository) FindTujuan(
	ctx context.Context,
	req domain.PenetapanOpdRequest,
) ([]domain.TujuanPenetapanOpd, error) {

	query := `
		SELECT
			id,
			kode_opd,
			kode_tujuan_opd,
			tujuan_opd,
			periode,
			tahun_aktif,
			created_date,
			last_modified_date,
			created_by
		FROM tb_tujuan_penetapan_opd
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
