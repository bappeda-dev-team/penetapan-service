package domain

import "time"

type TargetIndikatorTujuanPenetapanOpd struct {
	Id                int64
	IndikatorTujuanId int64
	KodeTarget        string
	Tahun             int
	Target            float64
	Satuan            string
	CreatedDate       time.Time
	LastModifiedDate  time.Time
	CreatedBy         *string
}
