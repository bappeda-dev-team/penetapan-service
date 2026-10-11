package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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

func (repo *PenetapanIndividuRepository) FindPkIndividus(
	ctx context.Context,
	req domain.SnapshotPenetapan,
) ([]domain.PkPenetapan, error) {

	const query = `
		SELECT
			pk.id,
			pk.pegawai_id,
			pk.kode_opd,
			pk.tahun_aktif,
			pk.level_pk,
			pk.kode_sasaran_opd,
			pk.kode_pk,
			pk.nama_pk,
			pk.keterangan_pk,
			pk.nama_pemilik_pk,
			pk.created_date,
			pk.last_modified_date,
			pk.created_by,
			pk.penetapan_individu_id,
			pk.anggaran_pk,
			pi.versi
		FROM pk_individu pk
		JOIN penetapan_individu pi ON pi.id = pk.penetapan_individu_id
		WHERE pk.pegawai_id = $1
		AND pk.kode_opd = $2
		AND pk.tahun_aktif = $3
		AND pk.penetapan_individu_id = $4
		ORDER BY level_pk ASC, kode_pk ASC
	`
	rows, err := repo.DB.QueryContext(ctx, query,
		req.PegawaiId,
		req.KodeOpd,
		req.Tahun,
		req.SnapshotId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	pks := make([]domain.PkPenetapan, 0)
	for rows.Next() {
		var pk domain.PkPenetapan
		err := rows.Scan(
			&pk.Id,
			&pk.PegawaiId,
			&pk.KodeOpd,
			&pk.TahunAktif,
			&pk.LevelPk,
			&pk.KodeSasaranOpd,
			&pk.KodePk,
			&pk.NamaPk,
			&pk.KeteranganPk,
			&pk.NamaPemilikPk,
			&pk.CreatedDate,
			&pk.LastModifiedDate,
			&pk.CreatedBy,
			&pk.PenetapanIndividuId,
			&pk.AnggaranPk,
			&pk.Versi,
		)
		if err != nil {
			return nil, err
		}
		pk.IndikatorPk = make([]domain.IndikatorPk, 0)
		pks = append(pks, pk)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return pks, nil

}

func (repo *PenetapanIndividuRepository) FindIndikatorPkByPkIds(
	ctx context.Context,
	pkIds []int64,
) ([]domain.IndikatorPk, error) {
	if len(pkIds) == 0 {
		return []domain.IndikatorPk{}, nil
	}

	var (
		placeholders []string
		args         []any
	)

	for i, id := range pkIds {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		args = append(args, id)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			pk_individu_id,
			kode_opd,
			tahun_aktif,
			kode_indikator_pk,
			nama_indikator_pk,
			rumus_perhitungan,
			sumber_data,
			definisi_operasional,
			created_date,
			last_modified_date,
			created_by
		FROM indikator_pk
		WHERE pk_individu_id IN (%s)
		ORDER BY pk_individu_id, kode_indikator_pk
	`, strings.Join(placeholders, ","))

	rows, err := repo.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]domain.IndikatorPk, 0)

	for rows.Next() {
		var item domain.IndikatorPk

		err := rows.Scan(
			&item.Id,
			&item.IdPk,
			&item.KodeOpd,
			&item.TahunAktif,
			&item.KodeIndikatorPk,
			&item.NamaIndikatorPk,
			&item.RumusPerhitungan,
			&item.SumberData,
			&item.DefinisiOperasional,
			&item.CreatedDate,
			&item.LastModifiedDate,
			&item.CreatedBy,
		)
		if err != nil {
			return nil, err
		}

		item.TargetPk = make([]domain.TargetPk, 0)

		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *PenetapanIndividuRepository) FindTargetPkByIndikatorIds(
	ctx context.Context,
	indikatorIds []int64,
) ([]domain.TargetPk, error) {
	if len(indikatorIds) == 0 {
		return []domain.TargetPk{}, nil
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
			indikator_pk_id,
			kode_target_pk,
			tahun,
			target,
			satuan,
			created_date,
			last_modified_date,
			created_by
		FROM target_pk
		WHERE indikator_pk_id IN (%s)
		ORDER BY indikator_pk_id, tahun
	`, strings.Join(placeholders, ","))

	rows, err := repo.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]domain.TargetPk, 0)

	for rows.Next() {
		var item domain.TargetPk

		err := rows.Scan(
			&item.Id,
			&item.IdIndikatorPk,
			&item.KodeTargetPk,
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
			kode_sasaran_opd,
			kode_pk,
			nama_pk,
			keterangan_pk,
			nama_pemilik_pk,
			anggaran_pk,
			versi,
			created_by
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13
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
		req.KodeSasaranOpd,
		req.KodePk,
		req.NamaPk,
		req.KeteranganPk,
		req.NamaPemilikPk,
		req.AnggaranPk,
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

func (repo *PenetapanIndividuRepository) SaveRenaksiIndividuPkPenetapan(
	ctx context.Context,
	tx *sql.Tx,
	req domain.RenaksiIndividu,
) (int64, error) {
	const query = `
		INSERT INTO renaksi_individu (
			pk_individu_id,
			kode_opd,
			tahun_aktif,
			kode_renaksi,
			urutan,
			nama_rencana_aksi,
			anggaran,
			created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
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
		req.KodeRenaksi,
		req.Urutan,
		req.NamaRencanaAksi,
		req.Anggaran,
		req.CreatedBy,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *PenetapanIndividuRepository) FindRenaksiIndividuByPkIds(
	ctx context.Context,
	pkIds []int64,
) ([]domain.RenaksiIndividu, error) {
	if len(pkIds) == 0 {
		return []domain.RenaksiIndividu{}, nil
	}

	var (
		placeholders []string
		args         []any
	)

	for i, id := range pkIds {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		args = append(args, id)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			pk_individu_id,
			kode_renaksi,
			urutan,
			nama_rencana_aksi,
			anggaran,
			created_by
		FROM renaksi_individu
		WHERE pk_individu_id IN (%s)
		ORDER BY pk_individu_id, id
	`, strings.Join(placeholders, ","))

	rows, err := repo.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.RenaksiIndividu

	for rows.Next() {
		var item domain.RenaksiIndividu

		err := rows.Scan(
			&item.Id,
			&item.IdPk,
			&item.KodeRenaksi,
			&item.Urutan,
			&item.NamaRencanaAksi,
			&item.Anggaran,
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

func (repo *PenetapanIndividuRepository) SavePelaksananRenaksiIndividuPkPenetapan(
	ctx context.Context,
	tx *sql.Tx,
	req domain.PelaksanaanRenaksi,
) (int64, error) {
	const query = `
		INSERT INTO renaksi_individu_pelaksanaan (
			renaksi_individu_id,
			kode_pelaksanaan,
			bulan,
			bobot,
			created_by
		) VALUES (
			$1, $2, $3, $4, $5
		)
		RETURNING id
	`

	var id int64

	err := tx.QueryRowContext(
		ctx,
		query,
		req.IdRenaksiIndividu,
		req.KodePelaksanaan,
		req.Bulan,
		req.Bobot,
		req.CreatedBy,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *PenetapanIndividuRepository) FindPelaksanaanByRenaksiIndividuIds(
	ctx context.Context,
	renaksiIds []int64,
) ([]domain.PelaksanaanRenaksi, error) {
	if len(renaksiIds) == 0 {
		return []domain.PelaksanaanRenaksi{}, nil
	}

	var (
		placeholders []string
		args         []any
	)

	for i, id := range renaksiIds {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		args = append(args, id)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			kode_pelaksanaan,
			renaksi_individu_id,
			bulan,
			bobot
		FROM renaksi_individu_pelaksanaan
		WHERE renaksi_individu_id IN (%s)
		ORDER BY renaksi_individu_id, bulan
	`, strings.Join(placeholders, ","))

	rows, err := repo.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]domain.PelaksanaanRenaksi, 0)

	for rows.Next() {
		var item domain.PelaksanaanRenaksi

		err := rows.Scan(
			&item.Id,
			&item.KodePelaksanaan,
			&item.IdRenaksiIndividu,
			&item.Bulan,
			&item.Bobot,
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

func (repo *PenetapanIndividuRepository) FindPelaksanaanByRenaksiIndividuIdsAndBulan(
	ctx context.Context,
	renaksiIds []int64,
	bulan int,
) ([]domain.PelaksanaanRenaksi, error) {
	if len(renaksiIds) == 0 {
		return []domain.PelaksanaanRenaksi{}, nil
	}

	placeholders := make([]string, len(renaksiIds))
	args := make([]any, 0, len(renaksiIds)+1)

	for i, id := range renaksiIds {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args = append(args, id)
	}

	// bulan menjadi parameter terakhir
	args = append(args, bulan)
	bulanPlaceholder := fmt.Sprintf("$%d", len(args))

	query := fmt.Sprintf(`
		SELECT
			id,
			kode_pelaksanaan,
			renaksi_individu_id,
			bulan,
			bobot
		FROM renaksi_individu_pelaksanaan
		WHERE renaksi_individu_id IN (%s)
		  AND bulan = %s
		ORDER BY renaksi_individu_id, bulan
	`, strings.Join(placeholders, ","), bulanPlaceholder)

	rows, err := repo.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]domain.PelaksanaanRenaksi, 0)

	for rows.Next() {
		var item domain.PelaksanaanRenaksi

		if err := rows.Scan(
			&item.Id,
			&item.KodePelaksanaan,
			&item.IdRenaksiIndividu,
			&item.Bulan,
			&item.Bobot,
		); err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
