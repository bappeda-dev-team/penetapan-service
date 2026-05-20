package perencanaan

type UrusanDetailResponse struct {
	KodeOpd           string                      `json:"kode_opd"`
	TahunAwal         string                      `json:"tahun_awal"`
	TahunAkhir        string                      `json:"tahun_akhir"`
	PaguAnggaranTotal []PaguAnggaranTotalResponse `json:"pagu_total"`
	Urusan            []UrusanResponse            `json:"urusan"`
	Tahun             string                      `json:"tahun,omitempty"`
}

type UrusanResponse struct {
	Kode         string                      `json:"kode"`
	Nama         string                      `json:"nama"`
	Jenis        string                      `json:"jenis"`
	Anggaran     []PaguAnggaranTotalResponse `json:"anggaran,omitempty"`
	Indikator    []IndikatorMatrixResponse   `json:"indikator"`
	BidangUrusan []BidangUrusanResponse      `json:"bidang_urusan"`
}

type BidangUrusanResponse struct {
	Kode      string                      `json:"kode"`
	Nama      string                      `json:"nama"`
	Jenis     string                      `json:"jenis"`
	Anggaran  []PaguAnggaranTotalResponse `json:"anggaran,omitempty"`
	Indikator []IndikatorMatrixResponse   `json:"indikator"`
	Program   []ProgramResponse           `json:"program"`
}

type ProgramResponse struct {
	Kode      string                      `json:"kode"`
	Nama      string                      `json:"nama"`
	Jenis     string                      `json:"jenis"`
	Anggaran  []PaguAnggaranTotalResponse `json:"anggaran,omitempty"`
	Indikator []IndikatorMatrixResponse   `json:"indikator"`
	Kegiatan  []KegiatanResponse          `json:"kegiatan"`
}

type KegiatanResponse struct {
	Kode        string                      `json:"kode"`
	Nama        string                      `json:"nama"`
	Jenis       string                      `json:"jenis"`
	Anggaran    []PaguAnggaranTotalResponse `json:"anggaran,omitempty"`
	Indikator   []IndikatorMatrixResponse   `json:"indikator"`
	SubKegiatan []SubKegiatanResponse       `json:"subkegiatan"`
}

type SubKegiatanResponse struct {
	Kode          string                      `json:"kode"`
	Nama          string                      `json:"nama"`
	Jenis         string                      `json:"jenis"`
	Tahun         string                      `json:"tahun,omitempty"`
	PegawaiId     string                      `json:"pegawai_id"`
	NamaPegawai   string                      `json:"nama_pegawai"`
	Anggaran      []PaguAnggaranTotalResponse `json:"anggaran,omitempty"`
	TotalAnggaran int64                       `json:"total_anggaran,omitempty"`
	Indikator     []IndikatorMatrixResponse   `json:"indikator"`
}

type PaguAnggaranTotalResponse struct {
	Tahun        string `json:"tahun"`
	PaguAnggaran int64  `json:"pagu_indikatif"`
	JenisPagu    string `json:"jenis_pagu"`
}

type IndikatorMatrixResponse struct {
	Id            string `json:"id,omitempty"`
	KodeIndikator string `json:"kode_indikator"`
	Kode          string `json:"kode,omitempty"`
	KodeOpd       string `json:"kode_opd,omitempty"`
	ProgramId     string `json:"program_id,omitempty"`
	Indikator     string `json:"indikator"`
	Tahun         string `json:"tahun"`
	TargetId      string `json:"target_id"`
	Target        string `json:"target"`
	Satuan        string `json:"satuan"`
	StatusTarget  bool   `json:"status_target_renja,omitempty"`
}
