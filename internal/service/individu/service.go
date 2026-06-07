package individu

import (
	"context"
	"log/slog"

	"github.com/bappeda-dev-team/penetapan-service/internal/client/perencanaan"
	web "github.com/bappeda-dev-team/penetapan-service/internal/model/web/individu"
	"github.com/bappeda-dev-team/penetapan-service/internal/repository"
	"github.com/bappeda-dev-team/penetapan-service/internal/service/opd/sync"
)

type PenetapanIndividuService struct {
	Repo         *repository.PenetapanIndividuRepository
	Perencanaan  *perencanaan.PerencanaanClient
	SyncExecutor *sync.Registry
	Logger       *slog.Logger
}

func NewPenetapanIndividuService(
	repo *repository.PenetapanIndividuRepository,
	perencanaanClient *perencanaan.PerencanaanClient,
	syncExecutor *sync.Registry,
	// logger taruh paling bawah
	logger *slog.Logger,
) *PenetapanIndividuService {
	return &PenetapanIndividuService{
		SyncExecutor: syncExecutor,
		Logger:       logger,
	}
}

func (s *PenetapanIndividuService) FindRekinsIndividu(
	ctx context.Context,
	req web.SyncPenetapanIndividuRequest,
) (web.RekinPenetapanIndividuResponse, error) {
	return web.RekinPenetapanIndividuResponse{}, nil
}
