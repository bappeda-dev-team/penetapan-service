package kode

import (
	"fmt"
	"strings"
)

func KodeSasaranOpd(sasaranId string) string {
	sasaranId = strings.TrimSpace(sasaranId)
	if sasaranId == "" {
		return ""
	}
	if strings.HasPrefix(sasaranId, "SAS-OPD-") {
		return sasaranId
	}
	kodeSasaran := fmt.Sprintf("SAS-OPD-%s", sasaranId)

	return kodeSasaran
}

func KodeSasaranOpdInt(sasaranId int) string {
	kodeSasaran := fmt.Sprintf("SAS-OPD-%d", sasaranId)

	return kodeSasaran
}

func KodeIndikatorSasaranOpd(indikatorSasaranId string) string {
	indikatorSasaranId = strings.TrimSpace(indikatorSasaranId)
	if indikatorSasaranId == "" {
		return ""
	}
	for strings.HasPrefix(indikatorSasaranId, "IND-IND-") {
		indikatorSasaranId = strings.TrimPrefix(indikatorSasaranId, "IND-")
	}
	if strings.HasPrefix(indikatorSasaranId, "IND-") {
		return indikatorSasaranId
	}
	return fmt.Sprintf("IND-%s", indikatorSasaranId)
}

func KodeTargetSasaranOpd(targetSasaranId string) string {
	targetSasaranId = strings.TrimSpace(targetSasaranId)
	if targetSasaranId == "" {
		return ""
	}
	if strings.HasPrefix(targetSasaranId, "TGT-TRG-") {
		targetSasaranId = strings.TrimPrefix(targetSasaranId, "TGT-")
	}
	if strings.HasPrefix(targetSasaranId, "TRG-") || strings.HasPrefix(targetSasaranId, "TGT-") {
		return targetSasaranId
	}
	return fmt.Sprintf("TRG-%s", targetSasaranId)
}

func KodeRenaksiOpd(renaksiId int) string {
	return fmt.Sprintf("REN-OPD-%d", renaksiId)
}
