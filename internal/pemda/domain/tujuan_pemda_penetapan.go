package domain

import "time"

type TujuanPemdaPenetapan struct {
	Id               int64
	Visi             string
	Misi             string
	KodeTujuanPemda  string
	TujuanPemda      string
	Periode          string
	TahunAktif       int
	CreatedDate      time.Time
	LastModifiedDate time.Time
	CreatedBy        *string
	PenetapanPemdaId int64
	Versi            int
	Indikator        []IndikatorTujuanPemdaPenetapan
}

type IndikatorTujuanPemdaPenetapan struct {
	Id                  int64
	IdTujuanPemda       int64
	KodeIndikator       string
	Indikator           string
	RumusPerhitungan    *string
	SumberData          *string
	DefinisiOperasional *string
	TahunAktif          int
	CreatedDate         time.Time
	LastModifiedDate    time.Time
	CreatedBy           *string
	Target              []TargetIndikatorTujuanPemdaPenetapan
}

type TargetIndikatorTujuanPemdaPenetapan struct {
	Id                     int64
	IndikatorTujuanPemdaId int64
	KodeTarget             string
	Tahun                  int
	Target                 float64
	Satuan                 string
	CreatedDate            time.Time
	LastModifiedDate       time.Time
	CreatedBy              *string
}
