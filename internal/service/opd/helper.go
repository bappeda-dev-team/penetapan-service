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

func (s *PenetapanOpdService) findPaguRenjaUrusan(
	ctx context.Context,
	urusans *[]domain.RenjaUrusan,
) error {

	if len(*urusans) == 0 {
		return nil
	}

	// ambil seluruh subkegiatan id
	urIds := make(
		[]int64,
		0,
		len(*urusans),
	)

	for _, ur := range *urusans {
		urIds = append(
			urIds,
			ur.Id,
		)
	}

	paguAnggs, err := s.Repo.FindPaguRenjaUrusan(
		ctx,
		urIds,
	)
	if err != nil {
		return err
	}
	if len(paguAnggs) == 0 {
		return nil
	}
	paguMap := make(
		map[int64][]domain.AnggaranRenja,
	)
	for _, pagu := range paguAnggs {
		urusanId := *pagu.UrusanId
		paguMap[urusanId] =
			append(
				paguMap[urusanId],
				pagu,
			)
	}
	// tempel ke urusan
	for i := range *urusans {
		ur := &(*urusans)[i]
		ur.PaguAnggaran = paguMap[ur.Id]
	}

	return nil
}

func (s *PenetapanOpdService) findPaguRenjaBidangUrusan(
	ctx context.Context,
	bidangUrusans *[]domain.RenjaBidangUrusan,
) error {

	if len(*bidangUrusans) == 0 {
		return nil
	}

	ids := make(
		[]int64,
		0,
		len(*bidangUrusans),
	)

	for _, b := range *bidangUrusans {
		ids = append(
			ids,
			b.Id,
		)
	}

	pagus, err := s.Repo.FindPaguRenjaBidangUrusan(
		ctx,
		ids,
	)

	if err != nil {
		return err
	}

	if len(pagus) == 0 {
		return nil
	}

	paguMap := make(
		map[int64][]domain.AnggaranRenja,
	)

	for _, pagu := range pagus {

		if pagu.BidangUrusanId == nil {
			continue
		}

		paguMap[*pagu.BidangUrusanId] =
			append(
				paguMap[*pagu.BidangUrusanId],
				pagu,
			)
	}

	for i := range *bidangUrusans {

		b := &(*bidangUrusans)[i]

		b.PaguAnggaran =
			paguMap[b.Id]
	}

	return nil
}

func (s *PenetapanOpdService) findPaguRenjaProgram(
	ctx context.Context,
	programs *[]domain.RenjaProgram,
) error {

	if len(*programs) == 0 {
		return nil
	}

	ids := make(
		[]int64,
		0,
		len(*programs),
	)

	for _, p := range *programs {
		ids = append(
			ids,
			p.Id,
		)
	}

	pagus, err := s.Repo.FindPaguRenjaProgram(
		ctx,
		ids,
	)

	if err != nil {
		return err
	}

	if len(pagus) == 0 {
		return nil
	}

	paguMap := make(
		map[int64][]domain.AnggaranRenja,
	)

	for _, pagu := range pagus {

		if pagu.ProgramId == nil {
			continue
		}

		paguMap[*pagu.ProgramId] =
			append(
				paguMap[*pagu.ProgramId],
				pagu,
			)
	}

	for i := range *programs {

		p := &(*programs)[i]

		p.PaguAnggaran =
			paguMap[p.Id]
	}

	return nil
}

func (s *PenetapanOpdService) findPaguRenjaKegiatan(
	ctx context.Context,
	kegiatans *[]domain.RenjaKegiatan,
) error {

	if len(*kegiatans) == 0 {
		return nil
	}

	ids := make(
		[]int64,
		0,
		len(*kegiatans),
	)

	for _, k := range *kegiatans {
		ids = append(
			ids,
			k.Id,
		)
	}

	pagus, err := s.Repo.FindPaguRenjaKegiatan(
		ctx,
		ids,
	)

	if err != nil {
		return err
	}

	if len(pagus) == 0 {
		return nil
	}

	paguMap := make(
		map[int64][]domain.AnggaranRenja,
	)

	for _, pagu := range pagus {

		if pagu.KegiatanId == nil {
			continue
		}

		paguMap[*pagu.KegiatanId] =
			append(
				paguMap[*pagu.KegiatanId],
				pagu,
			)
	}

	for i := range *kegiatans {

		k := &(*kegiatans)[i]

		k.PaguAnggaran =
			paguMap[k.Id]
	}

	return nil
}

func (s *PenetapanOpdService) findPaguRenjaSubkegiatan(
	ctx context.Context,
	subkegiatans *[]domain.RenjaSubkegiatan,
) error {

	if len(*subkegiatans) == 0 {
		return nil
	}

	ids := make(
		[]int64,
		0,
		len(*subkegiatans),
	)

	for _, sub := range *subkegiatans {
		ids = append(
			ids,
			sub.Id,
		)
	}

	pagus, err := s.Repo.FindPaguRenjaSubkegiatan(
		ctx,
		ids,
	)

	if err != nil {
		return err
	}

	if len(pagus) == 0 {
		return nil
	}

	paguMap := make(
		map[int64][]domain.AnggaranRenja,
	)

	for _, pagu := range pagus {

		if pagu.SubkegiatanId == nil {
			continue
		}

		paguMap[*pagu.SubkegiatanId] =
			append(
				paguMap[*pagu.SubkegiatanId],
				pagu,
			)
	}

	for i := range *subkegiatans {

		sub := &(*subkegiatans)[i]

		sub.PaguAnggaran =
			paguMap[sub.Id]
	}

	return nil
}
