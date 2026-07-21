package repository

import (
	"database/sql"
)

type PenetapanPemdaRepository struct {
	DB *sql.DB
}

func NewPenetapanPemdaRepository(db *sql.DB) *PenetapanPemdaRepository {
	return &PenetapanPemdaRepository{
		DB: db,
	}
}
