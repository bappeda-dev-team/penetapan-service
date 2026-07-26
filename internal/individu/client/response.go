package client

type PkPenetapanResponse struct {
	Id               string        `json:"id"`
	KodeOpd          string        `json:"kode_opd"`
	NamaOpd          string        `json:"nama_opd"`
	LevelPk          int           `json:"level_pk"`
	NipPemilikPk     string        `json:"nip_pemilik_pk"`
	NamaPemilikPk    string        `json:"nama_pemilik_pk"`
	IdRekinPemilikPk string        `json:"id_rekin_pemilik_pk"`
	SasaranOpdId     int64         `json:"id_sasaran_opd"`
	RekinPemilikPk   string        `json:"rekin_pemilik_pk"`
	Tahun            int           `json:"tahun"`
	Keterangan       string        `json:"keterangan"`
	AnggaranPk       int           `json:"anggaran_pk"`
	Indikators       []IndikatorPk `json:"indikators"`
	Renaksis         []RenaksiItem `json:"renaksi"`
}

type IndikatorPk struct {
	IdRekin     string        `json:"id_rekin"`
	IdIndikator string        `json:"id_indikator"`
	Indikator   string        `json:"indikator"`
	Targets     []TargetIndPk `json:"targets"`
}

type TargetIndPk struct {
	IdIndikator string `json:"id_indikator"`
	IdTarget    string `json:"id_target"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
}

type RenaksiItem struct {
	Id               string         `json:"id_renaksi"`
	RencanaKinerjaId string         `json:"rekin_id"`
	KodeOpd          string         `json:"kode_opd,omitempty"`
	Urutan           int            `json:"urutan"`
	NamaRencanaAksi  string         `json:"nama_rencana_aksi"`
	Anggaran         int            `json:"anggaran"`
	Pelaksanaan      []BobotBulanan `json:"pelaksanaan"`
}

type BobotBulanan struct {
	Id    string `json:"id_pelaksanaan"`
	Bulan int    `json:"bulan"`
	Bobot int    `json:"bobot"`
}

type RenjaIndividuResponse struct {
	RekinId     string `json:"rekin_id"`
	KodeProgram string `json:"kode_program"`
	NamaProgram string `json:"nama_program"`
	PaguProgram int64  `json:"pagu_program"`

	KodeKegiatan string `json:"kode_kegiatan"`
	NamaKegiatan string `json:"nama_kegiatan"`
	PaguKegiatan int64  `json:"pagu_kegiatan"`

	KodeSubkegiatan string `json:"kode_subkegiatan"`
	NamaSubkegiatan string `json:"nama_subkegiatan"`
	PaguSubkegiatan int64  `json:"pagu_subkegiatan"`
}
