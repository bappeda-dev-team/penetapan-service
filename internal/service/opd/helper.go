package opd

import (
	"context"

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

func (s *PenetapanOpdService) findIndikatorRenjaProgram(
	ctx context.Context,
	programs *[]domain.RenjaProgram,
) error {

	if len(*programs) == 0 {
		return nil
	}

	// ambil seluruh subkegiatan id
	subIds := make(
		[]int64,
		0,
		len(*programs),
	)

	for _, sub := range *programs {
		subIds = append(
			subIds,
			sub.Id,
		)
	}

	// query sekali
	indikators, err := s.Repo.FindIndikatorRenjaProgram(
		ctx,
		subIds,
	)
	if err != nil {
		return err
	}

	if len(indikators) == 0 {
		return nil
	}

	// ambil indikator id
	indikatorIds := make(
		[]int64,
		0,
		len(indikators),
	)

	for _, ind := range indikators {
		indikatorIds = append(
			indikatorIds,
			ind.Id,
		)
	}

	targets, err := s.Repo.FindTargetIndikatorRenjaProgramBatch(
		ctx,
		indikatorIds,
	)
	if err != nil {
		return err
	}

	// indikator_id -> targets
	targetMap := make(
		map[int64][]domain.TargetIndikatorRenjaProgram,
	)

	for _, tgt := range targets {

		targetMap[tgt.IndikatorProgramId] =
			append(
				targetMap[tgt.IndikatorProgramId],
				tgt,
			)
	}

	// isi target ke indikator
	for i := range indikators {
		indikators[i].Targets = targetMap[indikators[i].Id]
	}

	// subkegiatan_id -> indikator
	indikatorMap := make(
		map[int64][]domain.IndikatorRenjaProgram,
	)

	for _, ind := range indikators {

		indikatorMap[ind.ProgramId] =
			append(
				indikatorMap[ind.ProgramId],
				ind,
			)
	}

	// tempel ke subkegiatan
	for i := range *programs {

		prg := &(*programs)[i]

		prg.Indikators = indikatorMap[prg.Id]
	}

	return nil
}

func (s *PenetapanOpdService) findIndikatorRenjaKegiatan(
	ctx context.Context,
	kegiatans *[]domain.RenjaKegiatan,
) error {

	if len(*kegiatans) == 0 {
		return nil
	}

	// ambil seluruh subkegiatan id
	subIds := make(
		[]int64,
		0,
		len(*kegiatans),
	)

	for _, sub := range *kegiatans {
		subIds = append(
			subIds,
			sub.Id,
		)
	}

	// query sekali
	indikators, err := s.Repo.FindIndikatorRenjaKegiatan(
		ctx,
		subIds,
	)
	if err != nil {
		return err
	}

	if len(indikators) == 0 {
		return nil
	}

	// ambil indikator id
	indikatorIds := make(
		[]int64,
		0,
		len(indikators),
	)

	for _, ind := range indikators {
		indikatorIds = append(
			indikatorIds,
			ind.Id,
		)
	}

	targets, err := s.Repo.FindTargetIndikatorRenjaKegiatanBatch(
		ctx,
		indikatorIds,
	)
	if err != nil {
		return err
	}

	// indikator_id -> targets
	targetMap := make(
		map[int64][]domain.TargetIndikatorRenjaKegiatan,
	)

	for _, tgt := range targets {

		targetMap[tgt.IndikatorKegiatanId] =
			append(
				targetMap[tgt.IndikatorKegiatanId],
				tgt,
			)
	}

	// isi target ke indikator
	for i := range indikators {
		indikators[i].Targets = targetMap[indikators[i].Id]
	}

	// subkegiatan_id -> indikator
	indikatorMap := make(
		map[int64][]domain.IndikatorRenjaKegiatan,
	)

	for _, ind := range indikators {

		indikatorMap[ind.KegiatanId] =
			append(
				indikatorMap[ind.KegiatanId],
				ind,
			)
	}

	// tempel ke subkegiatan
	for i := range *kegiatans {

		keg := &(*kegiatans)[i]

		keg.Indikators = indikatorMap[keg.Id]
	}

	return nil
}

func (s *PenetapanOpdService) findIndikatorRenjaSubkegiatan(
	ctx context.Context,
	subkegiatans *[]domain.RenjaSubkegiatan,
) error {

	if len(*subkegiatans) == 0 {
		return nil
	}

	// ambil seluruh subkegiatan id
	subIds := make(
		[]int64,
		0,
		len(*subkegiatans),
	)

	for _, sub := range *subkegiatans {
		subIds = append(
			subIds,
			sub.Id,
		)
	}

	// query sekali
	indikators, err := s.Repo.FindIndikatorRenjaSubkegiatan(
		ctx,
		subIds,
	)
	if err != nil {
		return err
	}

	if len(indikators) == 0 {
		return nil
	}

	// ambil indikator id
	indikatorIds := make(
		[]int64,
		0,
		len(indikators),
	)

	for _, ind := range indikators {
		indikatorIds = append(
			indikatorIds,
			ind.Id,
		)
	}

	targets, err := s.Repo.FindTargetIndikatorRenjaSubkegiatanBatch(
		ctx,
		indikatorIds,
	)
	if err != nil {
		return err
	}

	// indikator_id -> targets
	targetMap := make(
		map[int64][]domain.TargetIndikatorRenjaSubkegiatan,
	)

	for _, tgt := range targets {

		targetMap[tgt.IndikatorSubkegiatanId] =
			append(
				targetMap[tgt.IndikatorSubkegiatanId],
				tgt,
			)
	}

	// isi target ke indikator
	for i := range indikators {

		indikators[i].Targets = targetMap[indikators[i].Id]
	}

	// subkegiatan_id -> indikator
	indikatorMap := make(
		map[int64][]domain.IndikatorRenjaSubkegiatan,
	)

	for _, ind := range indikators {

		indikatorMap[ind.SubkegiatanId] =
			append(
				indikatorMap[ind.SubkegiatanId],
				ind,
			)
	}

	// tempel ke subkegiatan
	for i := range *subkegiatans {

		sub := &(*subkegiatans)[i]

		sub.Indikators = indikatorMap[sub.Id]
	}

	return nil
}
