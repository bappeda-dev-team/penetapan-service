package web

type RenaksiOpdPenetapanResponse struct {
	KodeOpd     string               `json:"kode_opd" example:"1.02.0.00.0.00.01.0000"`
	TahunAktif  int                  `json:"tahun_aktif" example:"2025"`
	Versi       int                  `json:"versi" example:"1"`
	IsLocked    bool                 `json:"is_locked" example:"true"`
	RenaksiOpds []RenaksiOpdResponse `json:""`
}

type RenaksiOpdResponse struct {
	Renaksi string `json:"renaksi" example:"terwujudnya xx"`
}
