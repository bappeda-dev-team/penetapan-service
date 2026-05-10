package domain

import "time"

type IndikatorTujuanPenetapanOpd struct {
	Id               int64
	IdTujuanOpd      int64
	KodeOpd          string
	Indikator        string
	RumusPerhitungan *string
	SumberData       *string
	TahunAktif       int
	CreatedDate      time.Time
	LastModifiedDate time.Time
	CreatedBy        *string
}
