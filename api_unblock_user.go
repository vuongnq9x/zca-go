package zca

import (
	"context"
	"net/http"
)

// UnblockUser unblocks a user.
func (a *API) UnblockUser(ctx context.Context, userID string) (string, error) {
	return call[string](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("friend")+"/api/friend/unblock", nil, true), map[string]any{"fid": userID, "imei": a.IMEI})
}
