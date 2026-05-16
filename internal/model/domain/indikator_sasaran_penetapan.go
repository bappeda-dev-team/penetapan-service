package domain

import "time"

type IndikatorSasaranPenetapanOpd struct {
	Id                  int64
	IdSasaranOpd        int64
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
	Target              []TargetIndikatorSasaranPenetapanOpd
}
