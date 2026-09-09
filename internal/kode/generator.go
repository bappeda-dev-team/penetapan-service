package kode

import (
	"fmt"
	"strings"
)

func KodeSasaranOpd(sasaranId string) string {
	kodeSasaran := fmt.Sprintf("SAS-OPD-%s", sasaranId)

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
