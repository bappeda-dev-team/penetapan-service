package web

type TujuanPenetapanOpdResponse struct {
	Id            int64                              `json:"id"`
	KodeOpd       string                             `json:"kode_opd"`
	KodeTujuanOpd *string                            `json:"kode_tujuan_opd,omitempty"`
	TujuanOpd     string                             `json:"tujuan_opd"`
	Periode       string                             `json:"periode"`
	TahunAktif    int                                `json:"tahun_aktif"`
	Indikator     []IndikatorTujuanPenetapanResponse `json:"indikator"`
}
