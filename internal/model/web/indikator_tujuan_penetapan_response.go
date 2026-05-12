package web

type IndikatorTujuanPenetapanResponse struct {
	Id               int64                     `json:"id" example:"1"`
	Indikator        string                    `json:"indikator" example:"Indikator kepuasa publik"`
	RumusPerhitungan *string                   `json:"rumus_perhitungan" example:"indeks kepuasan masyarakat"`
	SumberData       *string                   `json:"sumber_data" example:"bps"`
	TahunAktif       int                       `json:"tahun_aktif" example:"2025"`
	Target           []TargetIndikatorResponse `json:"target"`
}
