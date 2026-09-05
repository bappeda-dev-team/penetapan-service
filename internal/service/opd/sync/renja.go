package sync

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/bappeda-dev-team/penetapan-service/internal/client/perencanaan"
	"github.com/bappeda-dev-team/penetapan-service/internal/common"
	"github.com/bappeda-dev-team/penetapan-service/internal/model/domain"
	"github.com/bappeda-dev-team/penetapan-service/internal/model/web"
	"github.com/bappeda-dev-team/penetapan-service/internal/repository"
)

type RenjaSyncExecutor struct {
	Repo              *repository.PenetapanOpdRepository
	PerencanaanClient *perencanaan.PerencanaanClient
	Logger            *slog.Logger
}

func NewRenjaSyncExecutor(
	repo *repository.PenetapanOpdRepository,
	perencanaanClient *perencanaan.PerencanaanClient,
	logger *slog.Logger,
) *RenjaSyncExecutor {
	return &RenjaSyncExecutor{
		Repo:              repo,
		PerencanaanClient: perencanaanClient,
		Logger:            logger,
	}
}

func (ex *RenjaSyncExecutor) Sync(
	ctx context.Context,
	syncId int64,
	req *web.SyncPenetapanOpdRequest,
	currentUser string,
) (web.SyncPenetapanOpdSummary, error) {
	perencanaanRequest := perencanaan.PerencanaanRequest{
		KodeOpd: req.KodeOpd,
		Tahun:   req.Tahun,
	}
	perencanaanResponse, err := ex.PerencanaanClient.GetPenetapanRenjaOpd(ctx, perencanaanRequest)
	if err != nil {
		return web.SyncPenetapanOpdSummary{}, err
	}

	if !hasValidRenja(perencanaanResponse) {
		return web.SyncPenetapanOpdSummary{}, common.NewValidation(
			"tidak ada data penetapan renja OPD yang siap disinkronkan",
		)
	}

	// mulai butuh tx
	tx, err := ex.Repo.DB.BeginTx(ctx, nil)
	if err != nil {
		return web.SyncPenetapanOpdSummary{}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// cari latest versi penetapan renja
	// versi otomatis + 1 misal 0 / belum ada -> 1
	versiPenetapan, err := ex.Repo.GetPenetapanNextVersion(ctx, tx, req.KodeOpd, domain.JenisPenetapanRenja, req.Tahun)
	if err != nil {
		return web.SyncPenetapanOpdSummary{}, err
	}
	if versiPenetapan > 1 {
		// deactivate versi lama
		errDeact := ex.Repo.DeactivateOldSnapshot(ctx, tx, req.KodeOpd, domain.JenisPenetapanRenja, req.Tahun)
		if errDeact != nil {
			return web.SyncPenetapanOpdSummary{}, err
		}
	}

	// snapshot penetapan baru
	penetapan := domain.PenetapanOpd{
		KodeOpd:        req.KodeOpd,
		Tahun:          req.Tahun,
		JenisPenetapan: domain.JenisPenetapanRenja,
		Versi:          versiPenetapan,
		SnapshotStatus: domain.SnapshotStatusActive,
		GeneratedBy:    &currentUser,
		IsActive:       true,
	}
	// save snapshot penetapan baru
	penetapanId, err := ex.Repo.SavePenetapanOpd(ctx, tx, penetapan)
	if err != nil {
		ex.Logger.Error("save penetapan opd error",
			"penetapanRenja", penetapan)
		return web.SyncPenetapanOpdSummary{}, err
	}

	// snapshot renja
	snapshotRenjas, err := ex.toRenjaSnapshots(
		perencanaanResponse,
		currentUser,
		penetapanId,
		req.KodeOpd,
		req.Tahun,
	)
	if err != nil {
		return web.SyncPenetapanOpdSummary{}, err
	}

	// simpan snapshot
	var jumlahRenja int
	var jumlahUrusan int
	var jumlahPaguUrusan int
	var jumlahBidangUrusan int
	var jumlahPaguBidangUrusan int
	var jumlahProgram int
	var jumlahPaguProgram int
	var jumlahKegiatan int
	var jumlahPaguKegiatan int
	var jumlahSubkegiatan int
	var jumlahPaguSubkegiatan int
	var jumlahIndikator int
	var jumlahTarget int
	for _, urusan := range snapshotRenjas {
		urId, err := ex.Repo.SaveRenjaUrusan(ctx, tx, urusan)
		if err != nil {
			return web.SyncPenetapanOpdSummary{}, err
		}
		jumlahRenja += 1
		jumlahUrusan += 1

		for _, bidUr := range urusan.BidangUrusans {
			bidUrId, err := ex.Repo.SaveRenjaBidangUrusan(ctx, tx, bidUr)
			if err != nil {
				return web.SyncPenetapanOpdSummary{}, err
			}
			jumlahRenja += 1
			jumlahBidangUrusan += 1

			for _, prg := range bidUr.Programs {
				prgId, err := ex.Repo.SaveRenjaProgram(ctx, tx, prg)
				if err != nil {
					return web.SyncPenetapanOpdSummary{}, err
				}
				jumlahRenja += 1
				jumlahProgram += 1

				// indikators
				for _, indikator := range prg.Indikators {

					indId, err := ex.Repo.SaveIndikatorRenjaProgram(ctx, tx, indikator, prgId)
					if err != nil {
						return web.SyncPenetapanOpdSummary{}, err
					}
					jumlahIndikator += 1

					// targets
					targets := indikator.Targets
					tgtInserted, err := ex.Repo.SaveTargetIndikatorRenjaProgramBatch(ctx, tx, indId, targets)
					if err != nil {
						return web.SyncPenetapanOpdSummary{}, err
					}
					jumlahTarget += tgtInserted
				}

				for _, keg := range prg.Kegiatans {
					kegId, err := ex.Repo.SaveRenjaKegiatan(ctx, tx, keg)
					if err != nil {
						return web.SyncPenetapanOpdSummary{}, err
					}
					jumlahRenja += 1
					jumlahKegiatan += 1
					// indikators
					for _, indikator := range keg.Indikators {

						indId, err := ex.Repo.SaveIndikatorRenjaKegiatan(ctx, tx, indikator, kegId)
						if err != nil {
							return web.SyncPenetapanOpdSummary{}, err
						}
						jumlahIndikator += 1

						// targets
						targets := indikator.Targets
						tgtInserted, err := ex.Repo.SaveTargetIndikatorRenjaKegiatanBatch(ctx, tx, indId, targets)
						if err != nil {
							return web.SyncPenetapanOpdSummary{}, err
						}
						jumlahTarget += tgtInserted
					}
					for _, sub := range keg.SubKegiatans {
						subId, err := ex.Repo.SaveRenjaSubkegiatan(ctx, tx, sub)
						if err != nil {
							return web.SyncPenetapanOpdSummary{}, err
						}
						jumlahRenja += 1
						jumlahSubkegiatan += 1
						// indikators
						for _, indikator := range sub.Indikators {

							indId, err := ex.Repo.SaveIndikatorRenjaSubkegiatan(ctx, tx, indikator, subId)
							if err != nil {
								return web.SyncPenetapanOpdSummary{}, err
							}
							jumlahIndikator += 1

							// targets
							targets := indikator.Targets
							tgtInserted, err := ex.Repo.SaveTargetIndikatorRenjaSubkegiatanBatch(ctx, tx, indId, targets)
							if err != nil {
								return web.SyncPenetapanOpdSummary{}, err
							}
							jumlahTarget += tgtInserted
						}

						// pagu subkegiatan
						jmlPagu, err := ex.Repo.SavePaguRenjaSubkegiatan(ctx, tx, subId, sub.PaguAnggaran)
						if err != nil {
							return web.SyncPenetapanOpdSummary{}, err
						}
						jumlahPaguSubkegiatan += jmlPagu

					}

					// pagu kegiatan
					jmlPagu, err := ex.Repo.SavePaguRenjaKegiatan(ctx, tx, kegId, keg.PaguAnggaran)
					if err != nil {
						return web.SyncPenetapanOpdSummary{}, err
					}
					jumlahPaguKegiatan += jmlPagu

				}

				// pagu program
				jmlPagu, err := ex.Repo.SavePaguRenjaProgram(ctx, tx, prgId, prg.PaguAnggaran)
				if err != nil {
					return web.SyncPenetapanOpdSummary{}, err
				}
				jumlahPaguProgram += jmlPagu

			}

			// pagu bidang urusan
			jmlPagu, err := ex.Repo.SavePaguRenjaBidangUrusan(ctx, tx, bidUrId, bidUr.PaguAnggaran)
			if err != nil {
				return web.SyncPenetapanOpdSummary{}, err
			}
			jumlahPaguBidangUrusan += jmlPagu
		}

		// pagu urusan
		jmlPagu, err := ex.Repo.SavePaguRenjaUrusan(ctx, tx, urId, urusan.PaguAnggaran)
		if err != nil {
			return web.SyncPenetapanOpdSummary{}, err
		}
		jumlahPaguUrusan += jmlPagu
	}

	// commit transaction penetapan tujuan
	err = tx.Commit()
	if err != nil {
		return web.SyncPenetapanOpdSummary{}, err
	}

	return web.SyncPenetapanOpdSummary{
		Renja:     &jumlahRenja,
		Indikator: jumlahIndikator,
		Target:    jumlahTarget,
	}, nil
}

func (ex *RenjaSyncExecutor) toRenjaSnapshots(
	renjaPerencanaans []perencanaan.UrusanDetailResponse,
	currentUser string,
	penetapanId int64,
	kodeOpd string,
	tahun int,
) ([]domain.RenjaUrusan, error) {

	createdBy := &currentUser

	var urusans []domain.RenjaUrusan
	for _, renja := range renjaPerencanaans {
		for _, urusan := range renja.Urusan {
			var bidangUrusans []domain.RenjaBidangUrusan
			for _, bidur := range urusan.BidangUrusan {
				var programs []domain.RenjaProgram
				for _, prog := range bidur.Program {
					var kegiatans []domain.RenjaKegiatan
					for _, keg := range prog.Kegiatan {
						var subkegiatans []domain.RenjaSubkegiatan
						for _, sub := range keg.SubKegiatan {
							snapshotSub, err := ex.toSubkegiatanSnapshot(
								sub, keg.Kode, createdBy, penetapanId, kodeOpd, tahun,
							)
							if err != nil {
								continue
							}
							subkegiatans = append(subkegiatans, snapshotSub)
						}
						kegSnapshot, err := ex.toKegiatanSnapshot(
							keg, prog.Kode, createdBy, penetapanId, kodeOpd, tahun,
						)
						if err != nil {
							continue
						}
						kegSnapshot.SubKegiatans = subkegiatans
						kegiatans = append(kegiatans, kegSnapshot)
					}
					progSnapshot, err := ex.toProgramSnapshot(
						prog, bidur.Kode, createdBy, penetapanId, kodeOpd, tahun,
					)
					if err != nil {
						continue
					}
					progSnapshot.Kegiatans = kegiatans
					programs = append(programs, progSnapshot)
				}
				bidSnapshot, err := ex.toBidangUrusanSnapshot(
					bidur, urusan.Kode, createdBy, penetapanId, kodeOpd, tahun,
				)
				if err != nil {
					continue
				}
				bidSnapshot.Programs = programs
				bidangUrusans = append(bidangUrusans, bidSnapshot)
			}
			urSnapshot, err := ex.toUrusanSnapshot(
				urusan, createdBy, penetapanId, kodeOpd, tahun,
			)
			if err != nil {
				continue
			}
			urSnapshot.BidangUrusans = bidangUrusans
			urusans = append(urusans, urSnapshot)
		}
	}
	return urusans, nil
}

func (ex *RenjaSyncExecutor) toUrusanSnapshot(
	renjaUrusan perencanaan.UrusanResponse,
	createdBy *string,
	penetapanId int64,
	kodeOpd string,
	tahun int,
) (domain.RenjaUrusan, error) {
	var paguAnggarans []domain.AnggaranRenja
	for _, angg := range renjaUrusan.Anggaran {
		kodePagu := fmt.Sprintf("PAGU-UR-%s-%s-%s", renjaUrusan.Kode, angg.Tahun, renjaUrusan.Jenis)
		paguAnggarans = append(paguAnggarans,
			domain.AnggaranRenja{
				KodePagu:     kodePagu,
				PaguAnggaran: angg.PaguAnggaran,
				Tahun:        tahun,
				JenisPagu:    angg.JenisPagu,
			})
	}

	return domain.RenjaUrusan{
		PenetapanId:  penetapanId,
		KodeOpd:      kodeOpd,
		KodeUrusan:   renjaUrusan.Kode,
		Urusan:       renjaUrusan.Nama,
		TahunAktif:   tahun,
		CreatedBy:    createdBy,
		PaguAnggaran: paguAnggarans,
	}, nil
}

func (ex *RenjaSyncExecutor) toBidangUrusanSnapshot(
	renjaBidangUrusan perencanaan.BidangUrusanResponse,
	kodeUrusan string,
	createdBy *string,
	penetapanId int64,
	kodeOpd string,
	tahun int,
) (domain.RenjaBidangUrusan, error) {
	var paguAnggarans []domain.AnggaranRenja
	for _, angg := range renjaBidangUrusan.Anggaran {
		kodePagu := fmt.Sprintf("PAGU-BIDUR-%s-%s-%s", renjaBidangUrusan.Kode, angg.Tahun, renjaBidangUrusan.Jenis)
		paguAnggarans = append(paguAnggarans,
			domain.AnggaranRenja{
				KodePagu:     kodePagu,
				PaguAnggaran: angg.PaguAnggaran,
				Tahun:        tahun,
				JenisPagu:    angg.JenisPagu,
			})
	}

	return domain.RenjaBidangUrusan{
		PenetapanId:      penetapanId,
		KodeOpd:          kodeOpd,
		KodeUrusan:       kodeUrusan,
		KodeBidangUrusan: renjaBidangUrusan.Kode,
		BidangUrusan:     renjaBidangUrusan.Nama,
		TahunAktif:       tahun,
		CreatedBy:        createdBy,
		PaguAnggaran:     paguAnggarans,
	}, nil
}

func (ex *RenjaSyncExecutor) toProgramSnapshot(
	renjaProgram perencanaan.ProgramResponse,
	kodeBidangUrusan string,
	createdBy *string,
	penetapanId int64,
	kodeOpd string,
	tahun int,
) (domain.RenjaProgram, error) {

	indikators := make([]domain.IndikatorRenjaProgram, 0, len(renjaProgram.Indikator))

	for _, ind := range renjaProgram.Indikator {
		kodeTarget := fmt.Sprintf("TGT-%s", ind.TargetId)
		target, err := strconv.ParseFloat(ind.Target, 64)
		if err != nil {
			return domain.RenjaProgram{}, fmt.Errorf(
				"invalid target indikator %q value %q",
				ind.Indikator,
				ind.Target,
			)
		}
		targets := []domain.TargetIndikatorRenjaProgram{
			{
				KodeTarget: kodeTarget,
				Target:     target,
				Satuan:     ind.Satuan,
				Tahun:      tahun,
				CreatedBy:  createdBy,
			},
		}
		kodeIndikator := fmt.Sprintf("IND-%s", ind.Id)
		indikators = append(indikators,
			domain.IndikatorRenjaProgram{
				KodeOpd:       kodeOpd,
				KodeIndikator: kodeIndikator,
				Indikator:     ind.Indikator,
				Tahun:         tahun,
				CreatedBy:     createdBy,
				Targets:       targets,
			},
		)
	}

	var paguAnggarans []domain.AnggaranRenja
	for _, angg := range renjaProgram.Anggaran {
		kodePagu := fmt.Sprintf("PAGU-PRG-%s-%s-%s", renjaProgram.Kode, angg.Tahun, renjaProgram.Jenis)
		paguAnggarans = append(paguAnggarans,
			domain.AnggaranRenja{
				KodePagu:     kodePagu,
				PaguAnggaran: angg.PaguAnggaran,
				Tahun:        tahun,
				JenisPagu:    angg.JenisPagu,
			})
	}

	return domain.RenjaProgram{
		PenetapanId:      penetapanId,
		KodeOpd:          kodeOpd,
		KodeBidangUrusan: kodeBidangUrusan,
		KodeProgram:      renjaProgram.Kode,
		Program:          renjaProgram.Nama,
		TahunAktif:       tahun,
		CreatedBy:        createdBy,
		Indikators:       indikators,
		PaguAnggaran:     paguAnggarans,
	}, nil
}

func (ex *RenjaSyncExecutor) toKegiatanSnapshot(
	renjaKegiatan perencanaan.KegiatanResponse,
	kodeProgram string,
	createdBy *string,
	penetapanId int64,
	kodeOpd string,
	tahun int,
) (domain.RenjaKegiatan, error) {

	indikators := make([]domain.IndikatorRenjaKegiatan, 0, len(renjaKegiatan.Indikator))

	for _, ind := range renjaKegiatan.Indikator {
		kodeTarget := fmt.Sprintf("TGT-%s", ind.TargetId)
		target, err := strconv.ParseFloat(ind.Target, 64)
		if err != nil {
			return domain.RenjaKegiatan{}, fmt.Errorf(
				"invalid target indikator %q value %q",
				ind.Indikator,
				ind.Target,
			)
		}
		targets := []domain.TargetIndikatorRenjaKegiatan{
			{
				KodeTarget: kodeTarget,
				Target:     target,
				Satuan:     ind.Satuan,
				Tahun:      tahun,
				CreatedBy:  createdBy,
			},
		}
		kodeIndikator := fmt.Sprintf("IND-%s", ind.Id)
		indikators = append(indikators,
			domain.IndikatorRenjaKegiatan{
				KodeOpd:       kodeOpd,
				KodeIndikator: kodeIndikator,
				Indikator:     ind.Indikator,
				Tahun:         tahun,
				CreatedBy:     createdBy,
				Targets:       targets,
			},
		)
	}
	var paguAnggarans []domain.AnggaranRenja
	for _, angg := range renjaKegiatan.Anggaran {
		kodePagu := fmt.Sprintf("PAGU-KEG-%s-%s-%s", renjaKegiatan.Kode, angg.Tahun, renjaKegiatan.Jenis)
		paguAnggarans = append(paguAnggarans,
			domain.AnggaranRenja{
				KodePagu:     kodePagu,
				PaguAnggaran: angg.PaguAnggaran,
				Tahun:        tahun,
				JenisPagu:    angg.JenisPagu,
			})
	}

	return domain.RenjaKegiatan{
		PenetapanId:  penetapanId,
		KodeOpd:      kodeOpd,
		KodeProgram:  kodeProgram,
		KodeKegiatan: renjaKegiatan.Kode,
		Kegiatan:     renjaKegiatan.Nama,
		TahunAktif:   tahun,
		CreatedBy:    createdBy,
		Indikators:   indikators,
		PaguAnggaran: paguAnggarans,
	}, nil
}

func (ex *RenjaSyncExecutor) toSubkegiatanSnapshot(
	renjaSubkegiatan perencanaan.SubKegiatanResponse,
	kodeKegiatan string,
	createdBy *string,
	penetapanId int64,
	kodeOpd string,
	tahun int,
) (domain.RenjaSubkegiatan, error) {

	indikators := make([]domain.IndikatorRenjaSubkegiatan, 0, len(renjaSubkegiatan.Indikator))

	for _, ind := range renjaSubkegiatan.Indikator {
		kodeTarget := fmt.Sprintf("TGT-%s", ind.TargetId)
		target, err := strconv.ParseFloat(ind.Target, 64)
		if err != nil {
			return domain.RenjaSubkegiatan{}, fmt.Errorf(
				"invalid target indikator %q value %q",
				ind.Indikator,
				ind.Target,
			)
		}
		targets := []domain.TargetIndikatorRenjaSubkegiatan{
			{
				KodeTarget: kodeTarget,
				Target:     target,
				Satuan:     ind.Satuan,
				Tahun:      tahun,
				CreatedBy:  createdBy,
			},
		}
		kodeIndikator := fmt.Sprintf("IND-%s", ind.Id)
		indikators = append(indikators,
			domain.IndikatorRenjaSubkegiatan{
				KodeOpd:       kodeOpd,
				KodeIndikator: kodeIndikator,
				Indikator:     ind.Indikator,
				Tahun:         tahun,
				CreatedBy:     createdBy,
				Targets:       targets,
			},
		)
	}
	var paguAnggarans []domain.AnggaranRenja
	for _, angg := range renjaSubkegiatan.Anggaran {
		kodePagu := fmt.Sprintf("PAGU-SUBKEG-%s-%s-%s", renjaSubkegiatan.Kode, angg.Tahun, renjaSubkegiatan.Jenis)
		paguAnggarans = append(paguAnggarans,
			domain.AnggaranRenja{
				KodePagu:     kodePagu,
				PaguAnggaran: angg.PaguAnggaran,
				Tahun:        tahun,
				JenisPagu:    angg.JenisPagu,
			})
	}

	return domain.RenjaSubkegiatan{
		PenetapanId:     penetapanId,
		KodeOpd:         kodeOpd,
		KodeKegiatan:    kodeKegiatan,
		KodeSubkegiatan: renjaSubkegiatan.Kode,
		Subkegiatan:     renjaSubkegiatan.Nama,
		TahunAktif:      tahun,
		CreatedBy:       createdBy,
		Indikators:      indikators,
		PaguAnggaran:    paguAnggarans,
	}, nil
}

func hasValidRenja(renjaPerencanaans []perencanaan.UrusanDetailResponse) bool {
	for _, renja := range renjaPerencanaans {
		for _, urusan := range renja.Urusan {
			if !hasRenjaContent(urusan.Nama, urusan.Indikator, urusan.Anggaran) {
				continue
			}
			for _, bidur := range urusan.BidangUrusan {
				if !hasRenjaContent(bidur.Nama, bidur.Indikator, bidur.Anggaran) {
					continue
				}
				for _, prog := range bidur.Program {
					if !hasRenjaContent(prog.Nama, prog.Indikator, prog.Anggaran) {
						continue
					}
					for _, keg := range prog.Kegiatan {
						if !hasRenjaContent(keg.Nama, keg.Indikator, keg.Anggaran) {
							continue
						}
						for _, sub := range keg.SubKegiatan {
							if hasRenjaContent(sub.Nama, sub.Indikator, sub.Anggaran) {
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}

func hasRenjaContent(nama string, indikators []perencanaan.IndikatorMatrixResponse, anggarans []perencanaan.PaguAnggaranTotalResponse) bool {
	if strings.TrimSpace(nama) == "" {
		return false
	}
	if len(anggarans) > 0 {
		return true
	}
	for _, ind := range indikators {
		if strings.TrimSpace(ind.Indikator) == "" || strings.TrimSpace(ind.Satuan) == "" {
			continue
		}
		if _, err := strconv.ParseFloat(ind.Target, 64); err != nil {
			continue
		}
		return true
	}
	return false
}
