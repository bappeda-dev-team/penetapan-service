package web

type RenaksiOpdPenetapanResponse struct {
	KodeOpd     string               `json:"kode_opd" example:"1.02.0.00.0.00.01.0000"`
	TahunAktif  int                  `json:"tahun_aktif" example:"2025"`
	Versi       int                  `json:"versi" example:"1"`
	IsLocked    bool                 `json:"is_locked" example:"true"`
	RenaksiOpds []RenaksiOpdResponse `json:""`
}

type RenaksiOpdResponse struct {
	KodeRencanaAksiOpd string `json:"kode_rencana_aksi_opd" example:"RA-2026-000001"`
	KodeOpd            string `json:"kode_opd" example:"5.01.5.05.0.00.01.0000"`
	KodeSasaranOpd     string `json:"kode_sasaran_opd" example:"96623"`
	KodePk             string `json:"kode_pk" example:"REKIN-PEG-2026-88731"`
	NamaPk             string `json:"nama_pk" example:"Pendampingan penyusunan Renstra Perangkat Daerah"`
	PegawaiId          string `json:"pegawai_id" example:"198610022015052001"`
	KodeSubkegiatan    string `json:"kode_subkegiatan" example:"5.01.03.2.03.0002"`
	NamaSubkegiatan    string `json:"nama_subkegiatan" example:"5.01.03.2.03.0002"`
	AnggaranRenaksi    int64  `json:"anggaran_renaksi" example:"24000000"`
	Tahun              int    `json:"tahun" example:"2026"`
	Tw1                int    `json:"tw1" example:"0"`
	Tw2                int    `json:"tw2" example:"30"`
	Tw3                int    `json:"tw3" example:"54"`
	Tw4                int    `json:"tw4" example:"10"`
}
