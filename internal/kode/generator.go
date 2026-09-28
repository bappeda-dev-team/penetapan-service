package kode

import (
	"fmt"
	"strings"
)

func KodeSasaranOpd(sasaranId string) string {
	kodeSasaran := fmt.Sprintf("SAS-OPD-%s", sasaranId)

	return kodeSasaran
}

func KodeSasaranOpdInt(sasaranId int) string {
	kodeSasaran := fmt.Sprintf("SAS-OPD-%d", sasaranId)

	return kodeSasaran
}

func KodeIndikatorSasaranOpd(indikatorSasaranId string) string {
	if strings.TrimSpace(indikatorSasaranId) == "" {
		return ""
	}
	return fmt.Sprintf("IND-SAS-%s", indikatorSasaranId)
}

func KodeTargetSasaranOpd(targetSasaranId string) string {
	if strings.TrimSpace(targetSasaranId) == "" {
		return ""
	}
	return fmt.Sprintf("TGT-SAS-%s", targetSasaranId)
}

func KodeRenaksiOpd(renaksiId int) string {
	return fmt.Sprintf("REN-OPD-%d", renaksiId)
}
