package sync

import (
	"errors"

	"github.com/bappeda-dev-team/penetapan-service/internal/individu/domain"
)

type Registry struct {
	PkSyncExecutor            PenetapanSyncExecutor
	RenjaIndividuSyncExecutor PenetapanSyncExecutor
}

func (r *Registry) Get(
	jenis domain.JenisPenetapan,
) (PenetapanSyncExecutor, error) {

	switch jenis {

	case domain.JenisPenetapanPk:
		return r.PkSyncExecutor, nil
	case domain.JenisPenetapanRenjaIndividu:
		return r.RenjaIndividuSyncExecutor, nil
	}

	return nil, errors.New(
		"sync executor not found",
	)
}
