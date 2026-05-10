package api

import (
	"github.com/bappeda-dev-team/penetapan-service/internal/model/domain"
	"github.com/bappeda-dev-team/penetapan-service/internal/model/web"
)

func ToTujuanOpdResponse(tujuanOpd domain.TujuanPenetapanOpd) web.TujuanPenetapanOpdResponse {
	return web.TujuanPenetapanOpdResponse{
		Id:            tujuanOpd.Id,
		KodeOpd:       tujuanOpd.KodeOpd,
		KodeTujuanOpd: tujuanOpd.KodeTujuanOpd,
		TujuanOpd:     tujuanOpd.TujuanOpd,
		Periode:       tujuanOpd.Periode,
		TahunAktif:    tujuanOpd.TahunAktif,
	}
}

func ToSasaranOpdResponse(sasaranOpd domain.SasaranPenetapanOpd) web.SasaranPenetapanOpdResponse {
	return web.SasaranPenetapanOpdResponse{
		Id:             sasaranOpd.Id,
		KodeOpd:        sasaranOpd.KodeOpd,
		KodeSasaranOpd: sasaranOpd.KodeSasaranOpd,
		SasaranOpd:     sasaranOpd.SasaranOpd,
		Periode:        sasaranOpd.Periode,
		TahunAktif:     sasaranOpd.TahunAktif,
	}
}
