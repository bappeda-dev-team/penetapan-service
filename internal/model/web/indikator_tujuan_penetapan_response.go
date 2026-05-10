package web

type IndikatorTujuanPenetapanResponse struct {
	Id               int64                     `json:"id"`
	Indikator        string                    `json:"indikator"`
	RumusPerhitungan *string                   `json:"rumus_perhitungan"`
	SumberData       *string                   `json:"sumber_data"`
	TahunAktif       int                       `json:"tahun_aktif"`
	Target           []TargetIndikatorResponse `json:"target"`
}
