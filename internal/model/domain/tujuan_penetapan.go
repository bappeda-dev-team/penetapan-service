package domain

import "time"

type TujuanPenetapanOpd struct {
	Id               int64
	KodeOpd          string
	KodeTujuanOpd    *string
	TujuanOpd        string
	Periode          string
	TahunAktif       int
	CreatedDate      time.Time
	LastModifiedDate time.Time
	CreatedBy        *string
	PenetapanId      int64
	Versi            int
	Indikator        []IndikatorTujuanPenetapanOpd
}
