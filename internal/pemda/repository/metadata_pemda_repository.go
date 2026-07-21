package repository

import (
	"context"

	"github.com/bappeda-dev-team/penetapan-service/internal/pemda/domain"
)

func (repo *PenetapanPemdaRepository) InsertMetadata(
	ctx context.Context,
	metadata domain.SyncPenetapanMetadata,
) (int64, error) {

	query := `INSERT INTO
		  sync_penetapan_metadata_pemda
             	  (tahun,
 		   jenis_penetapan, status,
 		   started_at)
                  VALUES ($1, $2, $3, $4)
	   	  RETURNING id`

	var id int64
	err := repo.DB.QueryRowContext(
		ctx,
		query,
		metadata.Tahun,
		metadata.JenisPenetapan,
		metadata.Status,
		metadata.StartedAt,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *PenetapanPemdaRepository) UpdateMetadataStatus(
	ctx context.Context,
	metadataId int64,
	status string,
	errMessage *string,
) error {
	var query string
	if status == domain.SyncStatusSuccess ||
		status == domain.SyncStatusFailed {
		query = `
		UPDATE sync_penetapan_metadata_pemda
		SET
			status = $1,
			error_message = $2,
			finished_at = NOW(),
			last_modified_date = NOW()
		WHERE id = $3
	`
	} else {
		query = `
		UPDATE sync_penetapan_metadata_pemda
		SET
			status = $1,
			error_message = $2,
			last_modified_date = NOW()
		WHERE id = $3
	`
	}
	_, err := repo.DB.ExecContext(
		ctx,
		query,
		status,
		errMessage,
		metadataId,
	)

	return err
}
