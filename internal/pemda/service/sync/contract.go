package sync

import (
	"context"
	"github.com/bappeda-dev-team/penetapan-service/internal/pemda/web"
)

type PenetapanSyncExecutor interface {
	Sync(
		ctx context.Context,
		syncId int64,
		req *web.SyncPenetapanRequest,
		currentUser string,
	) (web.SyncPenetapanSummary, error)
}
