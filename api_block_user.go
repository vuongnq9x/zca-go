package zca

import (
	"context"
	"net/http"
)

// BlockUser blocks a user.
func (a *API) BlockUser(ctx context.Context, userID string) (string, error) {
	return call[string](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("friend")+"/api/friend/block", nil, true), map[string]any{"fid": userID, "imei": a.IMEI})
}
