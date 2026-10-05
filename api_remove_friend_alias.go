package zca

import (
	"context"
	"net/http"
)

// RemoveFriendAlias removes a friend's alias.
func (a *API) RemoveFriendAlias(ctx context.Context, friendID string) (string, error) {
	return call[string](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("alias")+"/api/alias/remove", nil, true), map[string]any{"friendId": friendID})
}
