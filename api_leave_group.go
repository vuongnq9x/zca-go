package zca

import (
	"context"
	"net/http"
)

type LeaveGroupResponse struct {
	MemberError []any `json:"memberError"`
}

// LeaveGroup leaves a group, optionally silently. Zalo may return code 166 if you are not a member.
func (a *API) LeaveGroup(ctx context.Context, groupID string, silent bool) (*LeaveGroupResponse, error) {
	s := 0
	if silent {
		s = 1
	}
	params := map[string]any{
		"grids":    []string{groupID},
		"imei":     a.IMEI,
		"silent":   s,
		"language": a.Language,
	}
	return call[*LeaveGroupResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("group")+"/api/group/leave", nil, true), params)
}
