package domain

import "time"

type IndikatorTujuanPenetapanOpd struct {
	Id                  int64
	IdTujuanOpd         int64
	KodeIndikator       string
	KodeOpd             string
	Indikator           string
	RumusPerhitungan    *string
	SumberData          *string
	DefinisiOperasional *string
	TahunAktif          int
	CreatedDate         time.Time
	LastModifiedDate    time.Time
	CreatedBy           *string
	PenetapanId         int64
	Target              []TargetIndikatorTujuanPenetapanOpd
}
