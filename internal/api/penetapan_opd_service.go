package api

import (
	"context"

	"github.com/bappeda-dev-team/penetapan-service/internal/model/domain"
	"github.com/bappeda-dev-team/penetapan-service/internal/model/web"
)

type PenetapanOpdService struct {
	Repo *PenetapanOpdRepository
}

func (s *PenetapanOpdService) FindTujuan(
	ctx context.Context,
	req domain.PenetapanOpdRequest,
) ([]web.TujuanPenetapanOpdResponse, error) {
	tujuanOpd, err := s.Repo.FindTujuan(ctx, req)
	if err != nil {
		return nil, err
	}

	result := make(
		[]web.TujuanPenetapanOpdResponse,
		0,
		len(tujuanOpd),
	)
	for _, tujuan := range tujuanOpd {
		result = append(result, ToTujuanOpdResponse(tujuan))
	}

	return result, nil
}

func (s *PenetapanOpdService) FindSasaran(
	ctx context.Context,
	req domain.PenetapanOpdRequest,
) ([]web.SasaranPenetapanOpdResponse, error) {
	sasaranOpd, err := s.Repo.FindSasaran(ctx, req)
	if err != nil {
		return nil, err
	}

	result := make(
		[]web.SasaranPenetapanOpdResponse,
		0,
		len(sasaranOpd),
	)
	for _, sasaran := range sasaranOpd {
		result = append(result, ToSasaranOpdResponse(sasaran))
	}

	return result, nil
}
