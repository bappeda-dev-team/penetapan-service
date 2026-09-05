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
	IdRekin               string        `json:"id_rekin"`
	IdIndikator           string        `json:"id_indikator"`
	Indikator             string        `json:"indikator"`
	IdIndikatorSasaranOpd string        `json:"id_indikator_sasaran_opd"`
	Targets               []TargetIndPk `json:"targets"`
}

type TargetIndPk struct {
	IdIndikator        string `json:"id_indikator"`
	IdTarget           string `json:"id_target"`
	IdTargetSasaranOpd string `json:"id_target_sasaran_opd"`
	Target             string `json:"target"`
	Satuan             string `json:"satuan"`
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
	NamaPemilikPk    string `json:"nama_pemilik_pk"`
	NipPemilikPk     string `json:"nip_pemilik_pk"`
	IdRekinPemilikPk string `json:"id_rekin_pemilik_pk"`
	LevelPk          int    `json:"level_pk"`

	KodeProgram       string           `json:"kode_program"`
	NamaProgram       string           `json:"nama_program"`
	PaguProgram       int64            `json:"pagu_program"`
	IndikatorPrograms []IndikatorRenja `json:"indikator_programs"`

	KodeKegiatan       string           `json:"kode_kegiatan"`
	NamaKegiatan       string           `json:"nama_kegiatan"`
	PaguKegiatan       int64            `json:"pagu_kegiatan"`
	IndikatorKegiatans []IndikatorRenja `json:"indikator_kegiatans"`

	KodeSubkegiatan       string           `json:"kode_subkegiatan"`
	NamaSubkegiatan       string           `json:"nama_subkegiatan"`
	PaguSubkegiatan       int64            `json:"pagu_subkegiatan"`
	IndikatorSubkegiatans []IndikatorRenja `json:"indikator_subkegiatans"`
}

type IndikatorRenja struct {
	Id        string        `json:"id"`
	Indikator string        `json:"indikator"`
	Targets   []TargetRenja `json:"targets"`
}

type TargetRenja struct {
	Id     string `json:"id"`
	Target string `json:"target"`
	Satuan string `json:"satuan"`
	Tahun  int    `json:"tahun"`
}
