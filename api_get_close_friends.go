package zca

import (
	"context"
	"net/http"
)

// GetCloseFriends gets close friends.
func (a *API) GetCloseFriends(ctx context.Context) ([]User, error) {
	return call[[]User](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("profile")+"/api/social/friend/getclosedfriends", nil, true), map[string]any{})
}
