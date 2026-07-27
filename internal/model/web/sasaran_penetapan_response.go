package web

type SasaranPenetapanOpdResponse struct {
	KodeOpd    string               `json:"kode_opd" example:"1.02.0.00.0.00.01.0000"`
	TahunAktif int                  `json:"tahun_aktif" example:"2025"`
	Versi      int                  `json:"versi" example:"1"`
	IsLocked   bool                 `json:"is_locked" example:"true"`
	Sasarans   []SasaranOpdResponse `json:"sasaran_opds"`
}

type SasaranOpdResponse struct {
	Id             int64                               `json:"id"`
	KodeTujuanOpd  string                              `json:"kode_tujuan_opd" example:"TUJ-OPD-123"`
	KodeSasaranOpd string                              `json:"kode_sasaran_opd" example:"SAS-001"`
	SasaranOpd     string                              `json:"sasaran_opd" example:"Meningkatnya Nilai SAKIP"`
	Periode        string                              `json:"periode" example:"2025-2029"`
	Indikator      []IndikatorSasaranPenetapanResponse `json:"indikators"`
}
