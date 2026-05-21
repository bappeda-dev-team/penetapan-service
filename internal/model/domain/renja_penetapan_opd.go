package domain

import "time"

type RenjaPenetapanOpd struct {
	KodeOpd    string
	TahunAktif int
	Urusans    []RenjaUrusan
}

// URUSAN

type RenjaUrusan struct {
	Id               int64
	PenetapanId      int64
	KodeOpd          string
	KodeUrusan       string
	Urusan           string
	TahunAktif       int
	CreatedDate      time.Time
	LastModifiedDate time.Time
	CreatedBy        *string
	Versi            int
	BidangUrusans    []RenjaBidangUrusan
	PaguAnggaran     []AnggaranRenja
}

// BIDANG URUSAN

type RenjaBidangUrusan struct {
	Id               int64
	PenetapanId      int64
	KodeOpd          string
	KodeUrusan       string
	KodeBidangUrusan string
	BidangUrusan     string
	TahunAktif       int
	CreatedDate      time.Time
	LastModifiedDate time.Time
	CreatedBy        *string
	Versi            int
	Programs         []RenjaProgram
	PaguAnggaran     []AnggaranRenja
}

// PROGRAM

type RenjaProgram struct {
	Id               int64
	PenetapanId      int64
	KodeOpd          string
	KodeBidangUrusan string
	KodeProgram      string
	Program          string
	TahunAktif       int
	CreatedDate      time.Time
	LastModifiedDate time.Time
	CreatedBy        *string
	Versi            int
	Kegiatans        []RenjaKegiatan
	Indikators       []IndikatorRenjaProgram
	PaguAnggaran     []AnggaranRenja
}

type IndikatorRenjaProgram struct {
	Id               int64
	KodeOpd          string
	ProgramId        int64
	KodeIndikator    string
	Indikator        string
	Tahun            int
	CreatedDate      time.Time
	LastModifiedDate time.Time
	CreatedBy        *string
	PenetapanId      int64
	Targets          []TargetIndikatorRenjaProgram
	PaguAnggaran     []AnggaranRenja
}

type TargetIndikatorRenjaProgram struct {
	Id                 int64
	IndikatorProgramId int64
	KodeTarget         string
	Target             float64
	Satuan             string
	Tahun              int
	CreatedDate        time.Time
	LastModifiedDate   time.Time
	CreatedBy          *string
	PaguAnggaran       []AnggaranRenja
}

// KEGIATAN

type RenjaKegiatan struct {
	Id               int64
	PenetapanId      int64
	KodeOpd          string
	KodeProgram      string
	KodeKegiatan     string
	Kegiatan         string
	TahunAktif       int
	CreatedDate      time.Time
	LastModifiedDate time.Time
	CreatedBy        *string
	Versi            int
	SubKegiatans     []RenjaSubkegiatan
	Indikators       []IndikatorRenjaKegiatan
	PaguAnggaran     []AnggaranRenja
}

type IndikatorRenjaKegiatan struct {
	Id               int64
	KodeOpd          string
	KegiatanId       int64
	KodeIndikator    string
	Indikator        string
	Tahun            int
	CreatedDate      time.Time
	LastModifiedDate time.Time
	CreatedBy        *string
	PenetapanId      int64
	Targets          []TargetIndikatorRenjaKegiatan
}

type TargetIndikatorRenjaKegiatan struct {
	Id                  int64
	IndikatorKegiatanId int64
	KodeTarget          string
	Target              float64
	Satuan              string
	Tahun               int
	CreatedDate         time.Time
	LastModifiedDate    time.Time
	CreatedBy           *string
}

// SUBKEGIATAN

type RenjaSubkegiatan struct {
	Id               int64
	PenetapanId      int64
	KodeOpd          string
	KodeKegiatan     string
	KodeSubkegiatan  string
	Subkegiatan      string
	TahunAktif       int
	CreatedDate      time.Time
	LastModifiedDate time.Time
	CreatedBy        *string
	Versi            int
	Indikators       []IndikatorRenjaSubkegiatan
	PaguAnggaran     []AnggaranRenja
}

type IndikatorRenjaSubkegiatan struct {
	Id               int64
	KodeOpd          string
	SubkegiatanId    int64
	KodeIndikator    string
	Indikator        string
	Tahun            int
	CreatedDate      time.Time
	LastModifiedDate time.Time
	CreatedBy        *string
	Targets          []TargetIndikatorRenjaSubkegiatan
}

type TargetIndikatorRenjaSubkegiatan struct {
	Id                     int64
	IndikatorSubkegiatanId int64
	KodeTarget             string
	Target                 float64
	Satuan                 string
	Tahun                  int
	CreatedDate            time.Time
	LastModifiedDate       time.Time
	CreatedBy              *string
}

type AnggaranRenja struct {
	Id               int64
	KodePagu         string
	PaguAnggaran     int64
	UrusanId         *int64
	BidangUrusanId   *int64
	ProgramId        *int64
	KegiatanId       *int64
	SubkegiatanId    *int64
	Tahun            int
	JenisPagu        string
	CreatedDate      time.Time
	LastModifiedDate time.Time
	CreatedBy        *string
}
