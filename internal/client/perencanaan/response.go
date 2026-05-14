package perencanaan

type PerencanaanTujuanOpdResponse struct {
	KodeUrusan       string              `json:"kode_urusan"`
	Urusan           string              `json:"urusan"`
	KodeBidangUrusan string              `json:"kode_bidang_urusan"`
	NamaBidangUrusan string              `json:"nama_bidang_urusan"`
	KodeOpd          string              `json:"kode_opd"`
	NamaOpd          string              `json:"nama_opd"`
	TujuanOpd        []TujuanOpdResponse `json:"tujuan_opd"`
}

type TujuanOpdResponse struct {
	Id               int                 `json:"id_tujuan_opd"`
	KodeBidangUrusan string              `json:"kode_bidang_urusan,omitempty"`
	NamaBidangUrusan string              `json:"nama_bidang_urusan,omitempty"`
	KodeOpd          string              `json:"kode_opd,omitempty"`
	NamaOpd          string              `json:"nama_opd,omitempty"`
	Tujuan           string              `json:"tujuan,omitempty"`
	RumusPerhitungan string              `json:"rumus_perhitungan,omitempty"`
	SumberData       string              `json:"sumber_data,omitempty"`
	TahunAwal        string              `json:"tahun_awal,omitempty"`
	TahunAkhir       string              `json:"tahun_akhir,omitempty"`
	JenisPeriode     string              `json:"jenis_periode,omitempty"`
	Indikator        []IndikatorResponse `json:"indikator"`
}

type PerencanaanSasaranOpdResponse struct {
	IdPohon    int    `json:"id_pohon"`
	KodeOpd    string `json:"kode_opd,omitempty"`
	NamaOpd    string `json:"nama_opd,omitempty"`
	NamaPohon  string `json:"nama_pohon"`
	JenisPohon string `json:"jenis_pohon"`
	TahunPohon string `json:"tahun_pohon"`
	LevelPohon int    `json:"level_pohon"`
	// Pelaksana  []PelaksanaOpdResponse     `json:"pelaksana"`
	SasaranOpd []SasaranOpdDetailResponse `json:"sasaran_opd"`
}

type SasaranOpdDetailResponse struct {
	Id             string                     `json:"id"`
	NamaSasaranOpd string                     `json:"nama_sasaran_opd"`
	IdTujuanOpd    int                        `json:"id_tujuan_opd"`
	NamaTujuanOpd  string                     `json:"nama_tujuan_opd"`
	TahunAwal      string                     `json:"tahun_awal"`
	TahunAkhir     string                     `json:"tahun_akhir"`
	JenisPeriode   string                     `json:"jenis_periode"`
	Indikator      []IndikatorSasaranResponse `json:"indikator"`
}

// type PelaksanaOpdResponse struct {
// 	Id          string `json:"id"`
// 	PegawaiId   string `json:"pegawai_id"`
// 	Nip         string `json:"nip"`
// 	NamaPegawai string `json:"nama_pegawai"`
// }

type IndikatorSasaranResponse struct {
	Id                  string           `json:"id"`
	KodeIndikator       string           `json:"kode_indikator"`
	NamaIndikator       string           `json:"indikator"`
	RumusPerhitungan    string           `json:"rumus_perhitungan"`
	SumberData          string           `json:"sumber_data"`
	Jenis               string           `json:"jenis"`
	DefinisiOperasional string           `json:"definisi_operasional"`
	Target              []TargetResponse `json:"target"`
}

type IndikatorResponse struct {
	Id                  string           `json:"id"`
	KodeIndikator       string           `json:"kode_indikator"`
	IdTujuanOpd         int              `json:"id_tujuan_opd"`
	NamaIndikator       string           `json:"indikator"`
	RumusPerhitungan    string           `json:"rumus_perhitungan"`
	SumberData          string           `json:"sumber_data"`
	Jenis               string           `json:"jenis"`
	DefinisiOperasional string           `json:"definisi_operasional"`
	Target              []TargetResponse `json:"target"`
}

type TargetResponse struct {
	Id              string `json:"id"`
	IndikatorId     string `json:"indikator_id"`
	Tahun           string `json:"tahun"`
	TargetIndikator string `json:"target"`
	SatuanIndikator string `json:"satuan"`
}
