package zca

import (
	"context"
	"net/http"
)

// ChangeFriendAlias changes a friend's alias (nickname).
func (a *API) ChangeFriendAlias(ctx context.Context, alias, friendID string) (string, error) {
	return call[string](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("alias")+"/api/alias/update", nil, true), map[string]any{"friendId": friendID, "alias": alias, "imei": a.IMEI})
}
