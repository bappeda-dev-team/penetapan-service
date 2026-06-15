package client

type PkPenetapanResponse struct {

	Id               string        `json:"id"`
	KodeOpd          string        `json:"kode_opd"`
	NamaOpd          string        `json:"nama_opd"`
	LevelPk          int           `json:"level_pk"`
	NipPemilikPk     string        `json:"nip_pemilik_pk"`
	NamaPemilikPk    string        `json:"nama_pemilik_pk"`
	IdRekinPemilikPk string        `json:"id_rekin_pemilik_pk"`
	RekinPemilikPk   string        `json:"rekin_pemilik_pk"`
	Tahun            int           `json:"tahun"`
	Keterangan       string        `json:"keterangan"`
	Indikators       []IndikatorPk `json:"indikators"`
}

type IndikatorPk struct {
	IdRekin     string        `json:"id_rekin"`
	IdIndikator string        `json:"id_indikator"`
	Indikator   string        `json:"indikator"`
	Targets     []TargetIndPk `json:"targets"`
}

type TargetIndPk struct {
	IdIndikator string `json:"id_indikator"`
	IdTarget    string `json:"id_target"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
}
