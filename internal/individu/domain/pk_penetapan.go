package domain

import "time"

type PkPenetapan struct {
	Id         int64
	PegawaiId  string
	KodeOpd    string
	TahunAktif int

	LevelPk        int
	KodeSasaranOpd string
	KodePk         string // id rekin
	NamaPk         string // rekin
	KeteranganPk   string // Ket PK
	NamaPemilikPk  string // nama pegawai

	CreatedDate         time.Time
	LastModifiedDate    time.Time
	CreatedBy           *string
	PenetapanIndividuId int64
	Versi               int
	IndikatorPk         []IndikatorPk
}

type IndikatorPk struct {
	Id         int64
	IdPk       int64
	KodeOpd    string
	TahunAktif int
	// indikator
	KodeIndikatorPk string
	NamaIndikatorPk string
	// tambahan
	RumusPerhitungan    *string
	SumberData          *string
	DefinisiOperasional *string

	CreatedDate      time.Time
	LastModifiedDate time.Time
	CreatedBy        *string
	TargetPk         []TargetPk
}

type TargetPk struct {
	Id               int64
	IdIndikatorPk    int64
	KodeTargetPk     string
	Tahun            int
	Target           float64
	Satuan           string
	CreatedDate      time.Time
	LastModifiedDate time.Time
	CreatedBy        *string
}
