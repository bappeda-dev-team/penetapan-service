package domain

import "time"

type SnapshotPenetapan struct {
	Id             int64
	JenisSnapshot  JenisPenetapan
	PegawaiId      string
	KodeOpd        string
	Tahun          int
	SnapshotId     *int64
	Versi          int
	SnapshotStatus string
	GeneratedAt    time.Time
	GeneratedBy    *string
	IsActive       bool
	Bulan          int
}

type ActiveSnapshot struct {
	Id    int64
	Versi int
}
