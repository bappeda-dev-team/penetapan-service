package domain

import "time"

type SasaranPenetapanOpd struct {
	Id               int64
	KodeOpd          string
	KodeSasaranOpd   *string
	SasaranOpd       string
	Periode          string
	TahunAktif       int
	CreatedDate      time.Time
	LastModifiedDate time.Time
	CreatedBy        *string
}
