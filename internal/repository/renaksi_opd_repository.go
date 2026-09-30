package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/bappeda-dev-team/penetapan-service/internal/model/domain"
)

func (r *PenetapanOpdRepository) SaveRencanaAksiOpdBatch(
	ctx context.Context,
	tx *sql.Tx,
	data []domain.RencanaAksiOpd,
) (int, error) {

	if len(data) == 0 {
		return 0, nil
	}

	const columnCount = 18

	valueStrings := make([]string, 0, len(data))
	valueArgs := make([]any, 0, len(data)*columnCount)

	for i, item := range data {
		base := i * columnCount

		valueStrings = append(
			valueStrings,
			fmt.Sprintf(
				"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				base+1,
				base+2,
				base+3,
				base+4,
				base+5,
				base+6,
				base+7,
				base+8,
				base+9,
				base+10,
				base+11,
				base+12,
				base+13,
				base+14,
				base+15,
				base+16,
				base+17,
				base+18,
			),
		)

		valueArgs = append(
			valueArgs,
			item.KodeRencanaAksiOpd,
			item.PenetapanId,
			item.KodeOpd,
			item.KodeSasaranOpd,
			item.KodePk,
			item.NamaPk,
			item.PegawaiId,
			item.KodeSubkegiatan,
			item.NamaSubkegiatan,
			item.AnggaranRenaksi,
			item.Tahun,
			item.Tw1,
			item.Tw2,
			item.Tw3,
			item.Tw4,
			item.CreatedDate,
			item.LastModifiedDate,
			item.CreatedBy,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO rencana_aksi_opd
		(
			kode_rencana_aksi_opd,
			penetapan_id,
			kode_opd,
			kode_sasaran_opd,
			kode_pk,
			nama_pk,
			pegawai_id,
			kode_subkegiatan,
			nama_subkegiatan,
			anggaran_renaksi,
			tahun,
			tw1,
			tw2,
			tw3,
			tw4,
			created_date,
			last_modified_date,
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

func (r *PenetapanOpdRepository) FindRencanaAksiOpdBySnapshot(
	ctx context.Context,
	req domain.PenetapanOpdRequest,
) ([]domain.RencanaAksiOpd, error) {
	const op = "penetapan_opd_repository.FindRencanaAksiOpdBySnapshot"

	query := `
		SELECT
			id,
			kode_rencana_aksi_opd,
			penetapan_id,
			kode_opd,
			kode_sasaran_opd,
			kode_pk,
			nama_pk,
			pegawai_id,
			kode_subkegiatan,
			nama_subkegiatan,
			anggaran_renaksi,
			tahun,
			tw1,
			tw2,
			tw3,
			tw4,
			created_date,
			last_modified_date,
			created_by,
			versi
		FROM rencana_aksi_opd
		WHERE penetapan_id = $1
		ORDER BY id
	`

	rows, err := r.DB.QueryContext(ctx, query, req.SnapshotId)
	if err != nil {
		return nil, fmt.Errorf("%s: query: %w", op, err)
	}
	defer rows.Close()

	renaksis := make([]domain.RencanaAksiOpd, 0)

	for rows.Next() {
		var item domain.RencanaAksiOpd

		if err := rows.Scan(
			&item.Id,
			&item.KodeRencanaAksiOpd,
			&item.PenetapanId,
			&item.KodeOpd,
			&item.KodeSasaranOpd,
			&item.KodePk,
			&item.NamaPk,
			&item.PegawaiId,
			&item.KodeSubkegiatan,
			&item.NamaSubkegiatan,
			&item.AnggaranRenaksi,
			&item.Tahun,
			&item.Tw1,
			&item.Tw2,
			&item.Tw3,
			&item.Tw4,
			&item.CreatedDate,
			&item.LastModifiedDate,
			&item.CreatedBy,
			&item.Versi,
		); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}

		renaksis = append(renaksis, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows: %w", op, err)
	}

	return renaksis, nil
}
