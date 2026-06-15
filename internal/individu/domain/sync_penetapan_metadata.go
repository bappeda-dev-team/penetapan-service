package domain

import (
	"time"
)

type SyncPenetapanMetadata struct {
	Id               int64
	PegawaiId        string
	KodeOpd          string
	Tahun            int
	JenisPenetapan   JenisPenetapan
	Status           string
	StartedAt        time.Time
	FinishedAt       *time.Time
	SyncBy           *string
	ErrorMessage     *string
	CreatedDate      time.Time
	LastModifiedDate time.Time
}
