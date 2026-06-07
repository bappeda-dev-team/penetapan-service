package repository

import (
	"context"
	"database/sql"
	"github.com/bappeda-dev-team/penetapan-service/internal/model/domain/individu"
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
	req individu.PenetapanOpdRequest,
) error {
	return nil
}
