package zca

import (
	"context"
	"net/http"
)

// UpgradeGroupToCommunity upgrades a group to a community (verified 18+ account; Zalo error 185 = limit reached).
func (a *API) UpgradeGroupToCommunity(ctx context.Context, groupID string) (string, error) {
	return call[string](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group")+"/api/group/upgrade/community", nil, true),
		map[string]any{"grId": groupID, "language": a.Language})
}
