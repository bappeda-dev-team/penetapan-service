package web

import (
	"time"

	"github.com/bappeda-dev-team/penetapan-service/internal/validator"
)

type SyncPenetapanOpdResponse struct {
	SyncId           int64                   `json:"sync_id" example:"155"`
	Status           string                  `json:"status" example:"SUCCESS"`
	KodeOpd          string                  `json:"kode_opd" example:"1.22.33"`
	Tahun            int                     `json:"tahun" example:"2025"`
	JenisPenetapan   string                  `json:"jenis_penetapan" example:"tujuan"`
	ProcessedAt      time.Time               `json:"processed_at"`
	ProcessedSummary SyncPenetapanOpdSummary `json:"processed_summary"`
}

type SyncPenetapanOpdSummary struct {
	Tujuan  *int `json:"tujuan,omitempty" example:"1"`
	Sasaran *int `json:"sasaran,omitempty" example:"1"`
	Renja   *int `json:"renja,omitempty" example:"1"`

	Indikator int `json:"indikator" example:"3"`
	Target    int `json:"target" example:"4"`
}

type SyncPenetapanOpdRequest struct {
	KodeOpd        string `json:"kode_opd" example:"1.22.33"`
	Tahun          int    `json:"tahun" example:"2025"`
	JenisPenetapan string `json:"jenis_penetapan" example:"tujuan"`
}

func ValidateSyncPenetapanRequest(v *validator.Validator, syncRequest *SyncPenetapanOpdRequest) {
	v.Check(syncRequest.KodeOpd != "", "kode_opd", "tidak boleh kosong")
	v.Check(validator.Matches(syncRequest.KodeOpd, validator.KodeOpdRX), "kode_opd", "format kode tidak valid")

	v.Check(syncRequest.Tahun != 0, "tahun", "tidak boleh kosong")
	v.Check(syncRequest.Tahun >= 2020, "tahun", "tidak valid")
	v.Check(syncRequest.Tahun <= 2080, "tahun", "tidak valid")

	v.Check(syncRequest.JenisPenetapan != "", "jenis_penetapan", "tidak boleh kosong")

	v.Check(
		validator.PermittedValue(
			syncRequest.JenisPenetapan,
			"tujuan",
			"sasaran",
			"renja",
		),
		"jenis_penetapan",
		"jenis_penetapan tidak diketahui",
	)
}
