package zca

import (
	"context"
	"net/http"
)

// JoinGroupInviteBox joins a group from the invite box.
func (a *API) JoinGroupInviteBox(ctx context.Context, groupID string) (string, error) {
	return call[string](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group")+"/api/group/inv-box/join", nil, true), map[string]any{"grid": groupID, "lang": a.Language})
}
