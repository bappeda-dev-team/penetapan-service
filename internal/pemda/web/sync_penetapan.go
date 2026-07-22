package web

import (
	"time"

	"github.com/bappeda-dev-team/penetapan-service/internal/pemda/domain"
	"github.com/bappeda-dev-team/penetapan-service/internal/validator"
)

type SyncPenetapanResponse struct {
	SyncId           int64                 `json:"sync_id" example:"155"`
	Status           string                `json:"status" example:"SUCCESS"`
	Tahun            int                   `json:"tahun" example:"2025"`
	JenisPenetapan   domain.JenisPenetapan `json:"jenis_penetapan" example:"TUJUAN-PEMDA"`
	ProcessedAt      time.Time             `json:"processed_at"`
	ProcessedSummary SyncPenetapanSummary  `json:"processed_summary"`
}

type SyncPenetapanRequest struct {
	Tahun int `json:"tahun" example:"2025"`
}

type SyncPenetapanSummary struct {
	TujuanPemda  *int `json:"tujuan_pemda,omitempty" example:"1"`
	SasaranPemda *int `json:"sasaran_pemda,omitempty" example:"1"`

	Indikator int `json:"indikator" example:"3"`
	Target    int `json:"target" example:"4"`
}

func ValidateSyncPenetapanRequest(v *validator.Validator, syncRequest *SyncPenetapanRequest) {
	v.Check(syncRequest.Tahun != 0, "tahun", "tidak boleh kosong")
	v.Check(syncRequest.Tahun >= 2020, "tahun", "tidak valid")
	v.Check(syncRequest.Tahun <= 2080, "tahun", "tidak valid")
}
