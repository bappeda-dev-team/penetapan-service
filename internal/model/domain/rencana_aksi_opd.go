package domain

import (
	"time"
)

type RencanaAksiOpd struct {
	Id                 int64
	KodeRencanaAksiOpd string
	PenetapanId        int64
	KodeOpd            string
	KodeSasaranOpd     string
	KodePk             string
	// NAMA RENAKSI
	NamaPk           string
	PegawaiId        string
	NamaSubkegiatan  string
	KodeSubkegiatan  string
	AnggaranRenaksi  int64
	Tahun            int
	Tw1              int
	Tw2              int
	Tw3              int
	Tw4              int
	CreatedDate      time.Time
	LastModifiedDate time.Time
	CreatedBy        *string
	Versi            int
}
