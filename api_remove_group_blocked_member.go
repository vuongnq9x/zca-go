package zca

import (
	"context"
	"net/http"
)

// RemoveGroupBlockedMember unblocks member(s) of a group.
func (a *API) RemoveGroupBlockedMember(ctx context.Context, memberIDs []string, groupID string) (string, error) {
	return call[string](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group")+"/api/group/blockedmems/remove", nil, true), map[string]any{"grid": groupID, "members": memberIDs})
}
