package web

type SasaranPenetapanOpdResponse struct {
	Id             int64                               `json:"id"`
	KodeOpd        string                              `json:"kode_opd"`
	KodeSasaranOpd *string                             `json:"kode_sasaran_opd,omitempty"`
	SasaranOpd     string                              `json:"sasaran_opd"`
	Periode        string                              `json:"periode"`
	TahunAktif     int                                 `json:"tahun_aktif"`
	Indikator      []IndikatorSasaranPenetapanResponse `json:"indikator"`
}
