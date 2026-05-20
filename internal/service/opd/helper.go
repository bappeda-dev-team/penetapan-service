package opd

import (
	"github.com/bappeda-dev-team/penetapan-service/internal/model/domain"
	"github.com/bappeda-dev-team/penetapan-service/internal/model/web"
)

func ToTujuanOpdResponse(tujuanOpd domain.TujuanPenetapanOpd) web.TujuanOpdResponse {
	return web.TujuanOpdResponse{
		Id:            tujuanOpd.Id,
		KodeTujuanOpd: tujuanOpd.KodeTujuanOpd,
		TujuanOpd:     tujuanOpd.TujuanOpd,
		Periode:       tujuanOpd.Periode,
		Indikator:     []web.IndikatorTujuanPenetapanResponse{},
	}
}

func ToIndikatorTujuanOpdResponse(indikator domain.IndikatorTujuanPenetapanOpd) web.IndikatorTujuanPenetapanResponse {
	return web.IndikatorTujuanPenetapanResponse{
		Id:                  indikator.Id,
		KodeIndikator:       indikator.KodeIndikator,
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
		Id:         target.Id,
		KodeTarget: target.KodeTarget,
		Tahun:      target.Tahun,
		Target:     target.Target,
		Satuan:     target.Satuan,
	}
}

func ToSasaranOpdResponse(sasaranOpd domain.SasaranPenetapanOpd) web.SasaranOpdResponse {
	return web.SasaranOpdResponse{
		Id:             sasaranOpd.Id,
		KodeSasaranOpd: sasaranOpd.KodeSasaranOpd,
		SasaranOpd:     sasaranOpd.SasaranOpd,
		Periode:        sasaranOpd.Periode,
		Indikator:      []web.IndikatorSasaranPenetapanResponse{},
	}
}

func ToIndikatorSasaranOpdResponse(indikator domain.IndikatorSasaranPenetapanOpd) web.IndikatorSasaranPenetapanResponse {
	return web.IndikatorSasaranPenetapanResponse{
		Id:                  indikator.Id,
		KodeIndikator:       indikator.KodeIndikator,
		Indikator:           indikator.Indikator,
		RumusPerhitungan:    indikator.RumusPerhitungan,
		SumberData:          indikator.SumberData,
		DefinisiOperasional: indikator.DefinisiOperasional,
		TahunAktif:          indikator.TahunAktif,
		Target:              []web.TargetIndikatorResponse{},
	}
}

func ToTargetIndikatorSasaranOpdResponse(target domain.TargetIndikatorSasaranPenetapanOpd) web.TargetIndikatorResponse {
	return web.TargetIndikatorResponse{
		Id:         target.Id,
		KodeTarget: target.KodeTarget,
		Tahun:      target.Tahun,
		Target:     target.Target,
		Satuan:     target.Satuan,
	}
}
