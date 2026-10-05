package zca

import (
	"context"
	"net/http"
)

// AcceptFriendRequest accepts a friend request from a user.
func (a *API) AcceptFriendRequest(ctx context.Context, friendID string) (string, error) {
	return call[string](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("friend")+"/api/friend/accept", nil, true), map[string]any{"fid": friendID, "language": a.Language})
}
