package zca

import (
	"context"
	"net/http"
)

// RemoveFriend removes a friend.
func (a *API) RemoveFriend(ctx context.Context, friendID string) (string, error) {
	return call[string](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("friend")+"/api/friend/remove", nil, true), map[string]any{"fid": friendID, "imei": a.IMEI})
}
