package repository

import (
	"context"
	"database/sql"

	"github.com/bappeda-dev-team/penetapan-service/internal/individu/domain"
)

func (repo *PenetapanIndividuRepository) SaveRenjaIndividu(
	ctx context.Context,
	tx *sql.Tx,
	req domain.RenjaIndividu,
) (int64, error) {
	const query = `
		INSERT INTO renja_individu (
			penetapan_individu_id,
			kode_pk,
			pegawai_id,
			kode_opd,
			tahun_aktif,
			kode_program,
			nama_program,
			pagu_program,
			kode_kegiatan,
			nama_kegiatan,
			pagu_kegiatan,
			kode_subkegiatan,
			nama_subkegiatan,
			pagu_subkegiatan,
			created_by
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14,
			$15
		)
		RETURNING id
	`

	var id int64

	err := tx.QueryRowContext(
		ctx,
		query,
		req.PenetapanIndividuId,
		req.KodePk,
		req.PegawaiId,
		req.KodeOpd,
		req.TahunAktif,
		req.KodeProgram,
		req.NamaProgram,
		req.PaguProgram,
		req.KodeKegiatan,
		req.NamaKegiatan,
		req.PaguKegiatan,
		req.KodeSubkegiatan,
		req.NamaSubkegiatan,
		req.PaguSubkegiatan,
		req.CreatedBy,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
