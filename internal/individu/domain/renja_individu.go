package domain

import (
	"time"
)

type RenjaIndividu struct {
	Id int64

	PenetapanIndividuId int64
	KodePk              string

	PegawaiId  string
	KodeOpd    string
	TahunAktif int

	KodeProgram string
	NamaProgram string
	PaguProgram int64

	KodeKegiatan string
	NamaKegiatan string
	PaguKegiatan int64

	KodeSubkegiatan string
	NamaSubkegiatan string
	PaguSubkegiatan int64

	CreatedDate      time.Time
	LastModifiedDate time.Time
	CreatedBy        *string
}
