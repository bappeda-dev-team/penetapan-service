package domain

import "time"

type SnapshotPenetapan struct {
	Id             int64
	JenisSnapshot  JenisPenetapan
	Tahun          int
	SnapshotId     *int64
	Versi          int
	SnapshotStatus string
	GeneratedAt    time.Time
	GeneratedBy    *string
	IsActive       bool
}

type ActiveSnapshot struct {
	Id    int64
	Versi int
}
