package domain

import "time"

type PenetapanOpd struct {
	Id             int64
	KodeOpd        string
	Tahun          int
	JenisPenetapan string
	Versi          int
	SnapshotStatus string
	GeneratedAt    time.Time
	GeneratedBy    *string
	IsActive       bool
}
