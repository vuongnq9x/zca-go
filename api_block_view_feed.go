package zca

import (
	"context"
	"net/http"
)

// BlockViewFeed blocks/unblocks a friend from viewing your feed.
func (a *API) BlockViewFeed(ctx context.Context, isBlockFeed bool, userID string) (string, error) {
	block := 0
	if isBlockFeed {
		block = 1
	}
	return call[string](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("friend")+"/api/friend/feed/block", nil, true),
		map[string]any{"fid": userID, "isBlockFeed": block, "imei": a.IMEI})
}
