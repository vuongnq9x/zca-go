package zca

import (
	"context"
	"net/http"
)

// RejectFriendRequest rejects a friend request from a user.
func (a *API) RejectFriendRequest(ctx context.Context, friendID string) (string, error) {
	return call[string](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("friend")+"/api/friend/reject", nil, true), map[string]any{"fid": friendID})
}
