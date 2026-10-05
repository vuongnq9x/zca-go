package zca

import (
	"context"
	"net/http"
)

type RemoveUserFromGroupResponse struct {
	ErrorMembers []string `json:"errorMembers"`
}

// RemoveUserFromGroup kicks member(s) from a group. Zalo may return code 165 (not in group) or 166 (no permission).
func (a *API) RemoveUserFromGroup(ctx context.Context, memberIDs []string, groupID string) (*RemoveUserFromGroupResponse, error) {
	params := map[string]any{"grid": groupID, "members": memberIDs, "imei": a.IMEI}
	return call[*RemoveUserFromGroupResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("group")+"/api/group/kickout", nil, true), params)
}
