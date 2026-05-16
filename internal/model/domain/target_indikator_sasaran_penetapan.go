package domain

import "time"

type TargetIndikatorSasaranPenetapanOpd struct {
	Id                 int64
	IndikatorSasaranId int64
	KodeTarget         string
	Tahun              int
	Target             float64
	Satuan             string
	CreatedDate        time.Time
	LastModifiedDate   time.Time
	CreatedBy          *string
}
