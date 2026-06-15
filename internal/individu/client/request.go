package client

type SyncRequest struct {
	PegawaiId string `json:"pegawai_id"`
	KodeOpd   string `json:"kode_opd"`
	Tahun     int    `json:"tahun"`
}
