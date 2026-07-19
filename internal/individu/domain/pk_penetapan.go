package domain

import "time"

type PkPenetapan struct {
	Id         int64
	PegawaiId  string
	KodeOpd    string
	TahunAktif int

	LevelPk        int
	KodeSasaranOpd string
	KodePk         string // Snapshot Rekin
	NamaPk         string // Nama Rekin
	KeteranganPk   string // Keterangan PK
	NamaPemilikPk  string // Nama Pegawai

	CreatedDate         time.Time
	LastModifiedDate    time.Time
	CreatedBy           *string
	PenetapanIndividuId int64
	Versi               int

	IndikatorPk []IndikatorPk
	Renaksi     []RenaksiIndividu
}

type IndikatorPk struct {
	Id         int64
	IdPk       int64
	KodeOpd    string
	TahunAktif int

	KodeIndikatorPk string
	NamaIndikatorPk string

	RumusPerhitungan    *string
	SumberData          *string
	DefinisiOperasional *string

	CreatedDate      time.Time
	LastModifiedDate time.Time
	CreatedBy        *string

	TargetPk []TargetPk
}

type TargetPk struct {
	Id            int64
	IdIndikatorPk int64

	KodeTargetPk string
	Tahun        int
	Target       float64
	Satuan       string

	CreatedDate      time.Time
	LastModifiedDate time.Time
	CreatedBy        *string
}

type RenaksiIndividu struct {
	Id         int64
	IdPk       int64
	KodeOpd    string
	TahunAktif int

	KodeRenaksi     string
	Urutan          int
	NamaRencanaAksi string

	CreatedDate      time.Time
	LastModifiedDate time.Time
	CreatedBy        *string

	Pelaksanaan []PelaksanaanRenaksi
}

type PelaksanaanRenaksi struct {
	Id                int64
	IdRenaksiIndividu int64

	Bulan int
	Bobot int

	CreatedDate      time.Time
	LastModifiedDate time.Time
	CreatedBy        *string
}
