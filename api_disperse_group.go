package zca

import (
	"context"
	"net/http"
)

// DisperseGroup disbands a group.
func (a *API) DisperseGroup(ctx context.Context, groupID string) (string, error) {
	return call[string](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("group")+"/api/group/disperse", nil, true), map[string]any{"grid": groupID, "imei": a.IMEI})
}
