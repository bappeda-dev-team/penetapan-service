package perencanaan

type PerencanaanRequest struct {
	KodeOpd string `json:"kode_opd"`
	Tahun   int    `json:"tahun"`
}
