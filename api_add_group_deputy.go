package zca

import (
	"context"
	"net/http"
)

// AddGroupDeputy promotes member(s) to group deputy.
func (a *API) AddGroupDeputy(ctx context.Context, memberIDs []string, groupID string) (string, error) {
	return call[string](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group")+"/api/group/admins/add", nil, true), map[string]any{"grid": groupID, "members": memberIDs, "imei": a.IMEI})
}
