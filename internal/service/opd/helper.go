package opd

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
		Versi:         tujuanOpd.Versi,
		Indikator:     []web.IndikatorTujuanPenetapanResponse{},
	}
}

func ToIndikatorTujuanOpdResponse(indikator domain.IndikatorTujuanPenetapanOpd) web.IndikatorTujuanPenetapanResponse {
	return web.IndikatorTujuanPenetapanResponse{
		Id:                  indikator.Id,
		IdTujuanOpd:         indikator.IdTujuanOpd,
		Indikator:           indikator.Indikator,
		RumusPerhitungan:    indikator.RumusPerhitungan,
		SumberData:          indikator.SumberData,
		DefinisiOperasional: indikator.DefinisiOperasional,
		TahunAktif:          indikator.TahunAktif,
		Target:              []web.TargetIndikatorResponse{},
	}
}

func ToTargetIndikatorTujuanOpdResponse(target domain.TargetIndikatorTujuanPenetapanOpd) web.TargetIndikatorResponse {
	return web.TargetIndikatorResponse{
		Id:          target.Id,
		IndikatorId: target.IndikatorTujuanId,
		Tahun:       target.Tahun,
		Target:      target.Target,
		Satuan:      target.Satuan,
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
