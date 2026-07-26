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
			level_pk,
			pegawai_id,
			nama_pemilik_pk,
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
			$15, $16, $17
		)
		RETURNING id
	`

	var id int64

	err := tx.QueryRowContext(
		ctx,
		query,
		req.PenetapanIndividuId,
		req.KodePk,
		req.LevelPk,
		req.PegawaiId,
		req.NamaPemilikPk,
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

func (repo *PenetapanIndividuRepository) FindRenjaIndividu(
	ctx context.Context,
	req domain.SnapshotPenetapan,
) ([]domain.RenjaIndividu, error) {

	const query = `
		SELECT
			ri.id,
			ri.penetapan_individu_id,
			ri.level_pk,
			ri.kode_pk,
			ri.nama_pemilik_pk,
			ri.pegawai_id,
			ri.kode_opd,
			ri.tahun_aktif,
			ri.kode_program,
			ri.nama_program,
			ri.pagu_program,
			ri.kode_kegiatan,
			ri.nama_kegiatan,
			ri.pagu_kegiatan,
			ri.kode_subkegiatan,
			ri.nama_subkegiatan,
			ri.pagu_subkegiatan,
			ri.created_date,
			ri.last_modified_date,
			ri.created_by
		FROM renja_individu ri
		WHERE ri.pegawai_id = $1
		AND ri.kode_opd = $2
		AND ri.tahun_aktif = $3
		AND ri.penetapan_individu_id = $4
		ORDER BY
			ri.level_pk ASC,
			ri.kode_program ASC,
			ri.kode_kegiatan ASC,
			ri.kode_subkegiatan ASC
	`

	rows, err := repo.DB.QueryContext(
		ctx,
		query,
		req.PegawaiId,
		req.KodeOpd,
		req.Tahun,
		req.SnapshotId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	renjas := make([]domain.RenjaIndividu, 0)

	for rows.Next() {
		var renja domain.RenjaIndividu

		err := rows.Scan(
			&renja.Id,
			&renja.PenetapanIndividuId,
			&renja.LevelPk,
			&renja.KodePk,
			&renja.NamaPemilikPk,
			&renja.PegawaiId,
			&renja.KodeOpd,
			&renja.TahunAktif,
			&renja.KodeProgram,
			&renja.NamaProgram,
			&renja.PaguProgram,
			&renja.KodeKegiatan,
			&renja.NamaKegiatan,
			&renja.PaguKegiatan,
			&renja.KodeSubkegiatan,
			&renja.NamaSubkegiatan,
			&renja.PaguSubkegiatan,
			&renja.CreatedDate,
			&renja.LastModifiedDate,
			&renja.CreatedBy,
		)
		if err != nil {
			return nil, err
		}

		renjas = append(renjas, renja)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return renjas, nil
}
