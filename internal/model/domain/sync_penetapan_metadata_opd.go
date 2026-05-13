package domain

import "time"

type SyncPenetapanMetadataOpd struct {
	Id               int64
	KodeOpd          string
	Tahun            int
	JenisPenetapan   string
	Status           string
	StartedAt        time.Time
	FinishedAt       *time.Time
	SyncBy           *string
	ErrorMessage     *string
	CreatedDate      time.Time
	LastModifiedDate time.Time
}
