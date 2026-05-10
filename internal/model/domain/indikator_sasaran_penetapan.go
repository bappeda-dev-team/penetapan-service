package domain

import "time"

type IndikatorSasaranPenetapanOpd struct {
	Id               int64
	IdSasaranOpd     int64
	KodeOpd          string
	Indikator        string
	RumusPerhitungan *string
	SumberData       *string
	TahunAktif       int
	CreatedDate      time.Time
	LastModifiedDate time.Time
	CreatedBy        *string
}
