package sync

import (
	"errors"

	"github.com/bappeda-dev-team/penetapan-service/internal/pemda/domain"
)

type Registry struct {
	TujuanPemdaSyncExecutor PenetapanSyncExecutor
}

func (r *Registry) Get(
	jenis domain.JenisPenetapan,
) (PenetapanSyncExecutor, error) {

	switch jenis {

	case domain.JenisPenetapanTujuanPemda:
		return r.TujuanPemdaSyncExecutor, nil
	}

	return nil, errors.New(
		"sync executor not found",
	)
}
