package web

type IndikatorTujuanPenetapanResponse struct {
	Id                  int64                     `json:"id" example:"1"`
	KodeIndikator       string                    `json:"kode_indikator" example:"IND-TUJ-01"`
	Indikator           string                    `json:"indikator" example:"Indikator kepuasa publik"`
	RumusPerhitungan    *string                   `json:"rumus_perhitungan" example:"indeks kepuasan masyarakat"`
	SumberData          *string                   `json:"sumber_data" example:"bps"`
	DefinisiOperasional *string                   `json:"definisi_operasional" example:"definisi abc"`
	TahunAktif          int                       `json:"tahun_aktif" example:"2025"`
	Target              []TargetIndikatorResponse `json:"target"`
}
