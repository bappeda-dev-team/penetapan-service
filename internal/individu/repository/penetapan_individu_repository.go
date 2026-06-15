package repository

import (
	"context"
	"database/sql"

	"github.com/bappeda-dev-team/penetapan-service/internal/individu/domain"
)

type PenetapanIndividuRepository struct {
	DB *sql.DB
}

func NewPenetapanIndividuRepository(db *sql.DB) *PenetapanIndividuRepository {
	return &PenetapanIndividuRepository{
		DB: db,
	}
}

func (repo *PenetapanIndividuRepository) FindRekinBySnapshot(
	ctx context.Context,
	req domain.SnapshotPenetapan,
) error {
	return nil
}

func (repo *PenetapanIndividuRepository) SavePkPenetapan(
	ctx context.Context,
	tx *sql.Tx,
	req domain.PkPenetapan,
) (int64, error) {
	const query = `
		INSERT INTO pk_individu (
			penetapan_individu_id,
			pegawai_id,
			kode_opd,
			tahun_aktif,
			level_pk,
			kode_pk,
			nama_pk,
			keterangan_pk,
			nama_pemilik_pk,
			versi,
			created_by
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11
		)
		RETURNING id
	`

	var id int64

	err := tx.QueryRowContext(
		ctx,
		query,
		req.PenetapanIndividuId,
		req.PegawaiId,
		req.KodeOpd,
		req.TahunAktif,
		req.LevelPk,
		req.KodePk,
		req.NamaPk,
		req.KeteranganPk,
		req.NamaPemilikPk,
		req.Versi,
		req.CreatedBy,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *PenetapanIndividuRepository) SaveIndikatorPkPenetapan(
	ctx context.Context,
	tx *sql.Tx,
	req domain.IndikatorPk,
) (int64, error) {
	const query = `
		INSERT INTO indikator_pk (
			pk_individu_id,
			kode_opd,
			tahun_aktif,
			kode_indikator_pk,
			nama_indikator_pk,
			rumus_perhitungan,
			sumber_data,
			definisi_operasional,
			created_by
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9
		)
		RETURNING id
	`

	var id int64

	err := tx.QueryRowContext(
		ctx,
		query,
		req.IdPk,
		req.KodeOpd,
		req.TahunAktif,
		req.KodeIndikatorPk,
		req.NamaIndikatorPk,
		req.RumusPerhitungan,
		req.SumberData,
		req.DefinisiOperasional,
		req.CreatedBy,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *PenetapanIndividuRepository) SaveTargetPkPenetapan(
	ctx context.Context,
	tx *sql.Tx,
	req domain.TargetPk,
) (int64, error) {
	const query = `
		INSERT INTO target_pk (
			indikator_pk_id,
			kode_target_pk,
			tahun,
			target,
			satuan,
			created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
		RETURNING id
	`

	var id int64

	err := tx.QueryRowContext(
		ctx,
		query,
		req.IdIndikatorPk,
		req.KodeTargetPk,
		req.Tahun,
		req.Target,
		req.Satuan,
		req.CreatedBy,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
