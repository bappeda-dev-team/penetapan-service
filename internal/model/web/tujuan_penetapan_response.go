package web

type TujuanPenetapanOpdResponse struct {
	Id            int64                              `json:"id" example:"1"`
	KodeOpd       string                             `json:"kode_opd" example:"1.02.0.00.0.00.01.0000"`
	KodeTujuanOpd *string                            `json:"kode_tujuan_opd,omitempty" example:"TUJ-001"`
	TujuanOpd     string                             `json:"tujuan_opd" example:"Meningkatkan kualitas pelayanan publik"`
	Periode       string                             `json:"periode" example:"2025-2029"`
	TahunAktif    int                                `json:"tahun_aktif" example:"2025"`
	Indikator     []IndikatorTujuanPenetapanResponse `json:"indikator"`
}
