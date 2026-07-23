package helper

import (
	"strings"

	"github.com/bappeda-dev-team/penetapan-service/internal/pemda/client"
)

func IsValidTujuan(tujuan client.TujuanPemdaPenetapanResponse) bool {
	return strings.TrimSpace(TextCleaner(tujuan.TujuanPemda)) != ""
}

func IsValidSasaran(sasaran client.SasaranPemdaPenetapanResponse) bool {
	return strings.TrimSpace(TextCleaner(sasaran.SasaranPemda)) != ""
}

func IsValidIndikator(indikator client.IndikatorPenetapanDualResponse) bool {
	return strings.TrimSpace(TextCleaner(indikator.Indikator)) != "" &&
		len(indikator.TargetPenetapan) > 0
}

func IsValidTarget(target client.TargetDualResponse) bool {
	return target.Target != 0 &&
		strings.TrimSpace(target.Satuan) != ""
}
