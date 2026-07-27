package web

type TujuanPenetapanOpdResponse struct {
	KodeOpd    string              `json:"kode_opd" example:"1.02.0.00.0.00.01.0000"`
	TahunAktif int                 `json:"tahun_aktif" example:"2025"`
	Versi      int                 `json:"versi" example:"1"`
	IsLocked   bool                `json:"is_locked" example:"true"`
	Tujuans    []TujuanOpdResponse `json:"tujuan_opds"`
}

type TujuanOpdResponse struct {
	Id            int64                              `json:"id" example:"1"`
	KodeTujuanOpd string                             `json:"kode_tujuan_opd" example:"TUJ-001"`
	TujuanOpd     string                             `json:"tujuan_opd" example:"Meningkatkan kualitas pelayanan publik"`
	Periode       string                             `json:"periode" example:"2025-2029"`
	Indikator     []IndikatorTujuanPenetapanResponse `json:"indikators"`
}

type TujuanSasaranPenetapanOpdResponse struct {
	KodeOpd    string                     `json:"kode_opd" example:"1.02.0.00.0.00.01.0000"`
	TahunAktif int                        `json:"tahun_aktif" example:"2025"`
	Versi      int                        `json:"versi" example:"1"`
	IsLocked   bool                       `json:"is_locked" example:"true"`
	Tujuans    []TujuanSasaranOpdResponse `json:"tujuan_opds"`
}

type TujuanSasaranOpdResponse struct {
	Id            int64                              `json:"id" example:"1"`
	KodeTujuanOpd string                             `json:"kode_tujuan_opd" example:"TUJ-001"`
	TujuanOpd     string                             `json:"tujuan_opd" example:"Meningkatkan kualitas pelayanan publik"`
	Periode       string                             `json:"periode" example:"2025-2029"`
	Indikator     []IndikatorTujuanPenetapanResponse `json:"indikators"`
	SasaranOpds   []SasaranOpdResponse               `json:"sasaran_opds"`
}
