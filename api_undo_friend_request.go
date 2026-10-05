package zca

import (
	"context"
	"net/http"
)

// UndoFriendRequest cancels a sent friend request.
func (a *API) UndoFriendRequest(ctx context.Context, friendID string) (string, error) {
	return call[string](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("friend")+"/api/friend/undo", nil, true), map[string]any{"fid": friendID})
}
