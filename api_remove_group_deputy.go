package zca

import (
	"context"
	"net/http"
)

// RemoveGroupDeputy removes deputy role from member(s) of a group.
func (a *API) RemoveGroupDeputy(ctx context.Context, memberIDs []string, groupID string) (string, error) {
	params := map[string]any{"grid": groupID, "members": memberIDs, "imei": a.IMEI}
	return call[string](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group")+"/api/group/admins/remove", nil, true), params)
}
