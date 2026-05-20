package web

type RenjaPenetapanOpdResponse struct {
	KodeOpd       string                      `json:"kode_opd" example:"1.02.0.00.0.00.01.0000"`
	TahunAktif    int                         `json:"tahun_aktif" example:"2025"`
	Versi         int                         `json:"versi" example:"1"`
	IsLocked      bool                        `json:"is_locked" example:"true"`
	Urusans       []RenjaUrusanResponse       `json:"urusans"`
	BidangUrusans []RenjaBidangUrusanResponse `json:"bidang_urusans"`
	Programs      []RenjaProgramResponse      `json:"programs"`
	Kegiatans     []RenjaKegiatanResponse     `json:"kegiatans"`
	Subkegiatans  []RenjaSubkegiatanResponse  `json:"subkegiatans"`
}

type RenjaUrusanResponse struct {
	Id         int64  `json:"id"`
	KodeUrusan string `json:"kode_urusan" example:"1"`
	Urusan     string `json:"urusan" example:"URUSAN PENUNJANG"`
	IsLocked   bool   `json:"is_locked" example:"true"`
}

type RenjaBidangUrusanResponse struct {
	Id               int64  `json:"id"`
	KodeBidangUrusan string `json:"kode_bidang_urusan" example:"1.01"`
	BidangUrusan     string `json:"bidang_urusan" example:"URUSAN PEMERINTAHAN BIDANG PENDIDIKAN"`
	IsLocked         bool   `json:"is_locked" example:"true"`
}

type RenjaProgramResponse struct {
	Id          int64                    `json:"id"`
	KodeProgram string                   `json:"kode_program" example:"1.01.01.01"`
	Program     string                   `json:"program" example:"PROGRAM X"`
	IsLocked    bool                     `json:"is_locked" example:"true"`
	Indikators  []IndikatorRenjaResponse `json:"indikators"`
}

type RenjaKegiatanResponse struct {
	Id           int64                    `json:"id"`
	KodeKegiatan string                   `json:"kode_kegiatan" example:"1.01.01.01.01.01"`
	Kegiatan     string                   `json:"kegiatan" example:"KEGIATAN X"`
	IsLocked     bool                     `json:"is_locked" example:"true"`
	Indikators   []IndikatorRenjaResponse `json:"indikators"`
}

type RenjaSubkegiatanResponse struct {
	Id              int64                    `json:"id"`
	KodeSubkegiatan string                   `json:"kode_subkegiatan" example:"1.01.01.01.01.01.0001"`
	Subkegiatan     string                   `json:"subkegiatan" example:"SUBKEGIATAN X"`
	IsLocked        bool                     `json:"is_locked" example:"true"`
	Indikators      []IndikatorRenjaResponse `json:"indikators"`
}

type IndikatorRenjaResponse struct {
	Id            int64                     `json:"id"`
	KodeIndikator string                    `json:"kode_indikator" example:"IND-123"`
	Indikator     string                    `json:"indikator" example:"TEST-INDIKATOR"`
	Targets       []TargetIndikatorResponse `json:"targets"`
}
