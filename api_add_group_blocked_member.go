package zca

import (
	"context"
	"net/http"
)

// AddGroupBlockedMember blocks member(s) from a group.
func (a *API) AddGroupBlockedMember(ctx context.Context, memberIDs []string, groupID string) (string, error) {
	return call[string](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group")+"/api/group/blockedmems/add", nil, true), map[string]any{"grid": groupID, "members": memberIDs})
}
