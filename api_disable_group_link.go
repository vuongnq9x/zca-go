package zca

import (
	"context"
	"net/http"
)

// DisableGroupLink disables the group invite link.
func (a *API) DisableGroupLink(ctx context.Context, groupID string) (string, error) {
	return call[string](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group")+"/api/group/link/disable", nil, true), map[string]any{"grid": groupID})
}
