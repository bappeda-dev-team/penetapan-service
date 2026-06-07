package individu

type SyncPenetapanIndividuRequest struct {
	IdPegawai string `json:"id_pegawai" exmaple:"12345"`
	KodeOpd   string `json:"kode_opd" example:"1.22.33"`
	Tahun     int    `json:"tahun" example:"2025"`
}
