package web

import (
	"time"

	"github.com/bappeda-dev-team/penetapan-service/internal/individu/domain"
	"github.com/bappeda-dev-team/penetapan-service/internal/validator"
)

type SyncPenetapanResponse struct {
	SyncId           int64                 `json:"sync_id" example:"155"`
	Status           string                `json:"status" example:"SUCCESS"`
	PegawaiId        string                `json:"pegawai_id" example:"1998877665544"`
	KodeOpd          string                `json:"kode_opd" example:"1.22.33"`
	Tahun            int                   `json:"tahun" example:"2025"`
	JenisPenetapan   domain.JenisPenetapan `json:"jenis_penetapan" example:"REKIN"`
	ProcessedAt      time.Time             `json:"processed_at"`
	ProcessedSummary SyncPenetapanSummary  `json:"processed_summary"`
}

type SyncPenetapanSummary struct {
	Rekin *int `json:"rencana_kinerja,omitempty" example:"1"`

	Indikator     int `json:"indikator" example:"3"`
	Target        int `json:"target" example:"4"`
	Renaksi       int `json:"renaksi" example:"3"`
	RenjaIndividu int `json:"renja_individu" example:"6"`
}

type SyncPenetapanRequest struct {
	PegawaiId string `json:"pegawai_id" example:"199887766"`
	KodeOpd   string `json:"kode_opd" example:"1.22.33"`
	Tahun     int    `json:"tahun" example:"2025"`
}

func ValidateSyncPenetapanRequest(v *validator.Validator, syncRequest *SyncPenetapanRequest) {
	v.Check(syncRequest.PegawaiId != "", "pegawai_id", "tidak boleh kosong")
	v.Check(validator.Matches(syncRequest.PegawaiId, validator.PegawaiIdRX), "pegawai_id", "format pegawai_id tidak valid")
	v.Check(syncRequest.KodeOpd != "", "kode_opd", "tidak boleh kosong")
	v.Check(validator.Matches(syncRequest.KodeOpd, validator.KodeOpdRX), "kode_opd", "format kode tidak valid")

	v.Check(syncRequest.Tahun != 0, "tahun", "tidak boleh kosong")
	v.Check(syncRequest.Tahun >= 2020, "tahun", "tidak valid")
	v.Check(syncRequest.Tahun <= 2080, "tahun", "tidak valid")
}
