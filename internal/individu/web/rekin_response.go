package web

type RekinPenetapanIndividuResponse struct {
	IdPegawai  string                  `json:"pegawai_id" example:"12345"`
	Nama       string                  `json:"nama" example:"Pegawai X"`
	KodeOpd    string                  `json:"kode_opd" example:"1.23.456"`
	TahunAktif int                     `json:"tahun_aktif" example:"2025"`
	Rekins     []RekinIndividuResponse `json:"rekins"`
}

type RekinIndividuResponse struct {
	Id              int64                 `json:"id" example:"1"`
	LevelPk         int                   `json:"level_pk" example:"2"`
	KodeSasaranOpd  string                `json:"kode_sasaran_opd" example:"SAS-OPD-1"`
	KodePk          string                `json:"kode_pk" example:"RK-001"`
	Rekin           string                `json:"rekin" example:"Terwujudnya pelayanan publik yang berkualitas"`
	KeteranganPk    string                `json:"keterangan_pk,omitempty" example:"Kinerja utama individu"`
	NamaPemilikPk   string                `json:"nama_pemilik_pk" example:"Budi Santoso"`
	AnggaranPk      int                   `json:"anggaran_pk" example:"500000"`
	Versi           int                   `json:"versi" example:"1"`
	IndikatorPkList []IndikatorPkResponse `json:"indikator_pk"`
	Renaksis        []RenaksiPkResponse   `json:"renaksis"`
}

type IndikatorPkResponse struct {
	Id                  int64              `json:"id" example:"1"`
	KodeIndikatorPk     string             `json:"kode_indikator_pk" example:"IKU-001"`
	NamaIndikatorPk     string             `json:"nama_indikator_pk" example:"Persentase kepuasan masyarakat"`
	RumusPerhitungan    *string            `json:"rumus_perhitungan,omitempty" example:"Jumlah puas / total responden x 100%"`
	SumberData          *string            `json:"sumber_data,omitempty" example:"Survei Kepuasan Masyarakat"`
	DefinisiOperasional *string            `json:"definisi_operasional,omitempty" example:"Persentase responden yang menyatakan puas"`
	TargetPkList        []TargetPkResponse `json:"target_pk"`
}

type TargetPkResponse struct {
	Id           int64   `json:"id" example:"1"`
	KodeTargetPk string  `json:"kode_target_pk" example:"TGT-001"`
	Tahun        int     `json:"tahun" example:"2025"`
	Target       float64 `json:"target" example:"95"`
	Satuan       string  `json:"satuan" example:"persen"`
}

type RenaksiPkResponse struct {
	Id              int64                          `json:"id" example:"1"`
	UrutanRenaksi   int                            `json:"urutan_renaksi" example:"1"`
	KodeRenaksi     string                         `json:"kode_renaksi" example:"REN-001"`
	NamaRenaksi     string                         `json:"nama_renaksi" example:"Rapat Koordinasi"`
	AnggaranRenaksi int                            `json:"anggaran_renaksi" example:"100000"`
	Pelaksanaans    []PelaksanaanRenaksiPkResponse `json:"pelaksanaans"`
}

type PelaksanaanRenaksiPkResponse struct {
	Id               int64  `json:"id" example:"1"`
	KodePelaksanaan  string `json:"kode_pelaksanaan" example:"PEL-001"`
	BulanPelaksanaan int    `json:"bulan_pelaksanaan" example:"12"`
	BobotPelaksanaan int    `json:"bobot_pelaksanaan" example:"10"`
}

type RenjaIndividuResponse struct {
	IdPegawai  string          `json:"pegawai_id" example:"12345777"`
	Nama       string          `json:"nama" example:"Pegawai X"`
	KodeOpd    string          `json:"kode_opd" example:"1.23.456"`
	TahunAktif int             `json:"tahun_aktif" example:"2025"`
	Renjas     []RenjaIndividu `json:"renjas"`
}

type RenjaIndividu struct {
	Id          int64  `json:"id" example:"1"`
	KodePk      string `json:"kode_pk" example:"RK-001"`
	LevelPk     int    `json:"level_pk" example:"5"`
	IdPegawai   string `json:"pegawai_id" example:"12345"`
	NamaPegawai string `json:"nama_pegawai" example:"Pegawai X"`

	KodeProgram string `json:"kode_program" example:"8.01.05"`
	NamaProgram string `json:"nama_program" example:"PROGRAM ABC"`
	PaguProgram int64  `json:"pagu_program" example:"2500000000"`

	KodeKegiatan string `json:"kode_kegiatan" example:"8.01.05.2.01"`
	NamaKegiatan string `json:"nama_kegiatan" example:"KEGIATAN XX"`
	PaguKegiatan int64  `json:"pagu_kegiatan" example:"1250000000"`

	KodeSubkegiatan string `json:"kode_subkegiatan" example:"8.01.05.2.01.0003"`
	NamaSubkegiatan string `json:"nama_subkegiatan" example:"SUBKEGIATAN X"`
	PaguSubkegiatan int64  `json:"pagu_subkegiatan" example:"500000000"`
}
