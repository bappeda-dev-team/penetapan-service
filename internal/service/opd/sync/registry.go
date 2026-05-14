package sync

import (
	"errors"

	"github.com/bappeda-dev-team/penetapan-service/internal/model/domain"
)

type Registry struct {
	TujuanSyncExecutor  PenetapanSyncExecutor
	SasaranSyncExecutor PenetapanSyncExecutor
}

func (r *Registry) Get(
	jenis domain.JenisPenetapan,
) (PenetapanSyncExecutor, error) {

	switch jenis {

	case domain.JenisPenetapanTujuan:
		return r.TujuanSyncExecutor, nil

	case domain.JenisPenetapanSasaran:
		return r.SasaranSyncExecutor, nil
	}

	return nil, errors.New(
		"sync executor not found",
	)
}
