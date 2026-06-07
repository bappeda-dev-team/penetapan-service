package individu

type RekinPenetapanIndividuResponse struct {
	IdPegawai  string                  `json:"id_pegawai" example:"12345"`
	Nama       string                  `json:"nama" example:"Pegawai X"`
	KodeOpd    string                  `json:"kode_opd" example:"1.23.456"`
	NamaOpd    string                  `json:"nama_opd" example:"OPD X"`
	TahunAktif int                     `json:"tahun_aktif" example:"2025"`
	Rekins     []RekinIndividuResponse `json:"rekins"`
}

type RekinIndividuResponse struct {
	Rekin string `json:"rekin" example:"terwujudnya xxx"`
}
