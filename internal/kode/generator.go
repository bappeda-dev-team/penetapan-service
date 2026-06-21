package kode

import (
	"fmt"
)

func KodeSasaranOpd(sasaranId string) string {
	kodeSasaran := fmt.Sprintf("SAS-OPD-%s", sasaranId)

	return kodeSasaran
}
