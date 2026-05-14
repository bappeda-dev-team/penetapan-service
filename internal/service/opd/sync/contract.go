package sync

import (
	"context"

	"github.com/bappeda-dev-team/penetapan-service/internal/model/web"
)

type PenetapanSyncExecutor interface {
	Sync(
		ctx context.Context,
		syncId int64,
		req *web.SyncPenetapanOpdRequest,
		currentUser string,
	) (web.SyncPenetapanOpdSummary, error)
}
