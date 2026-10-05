package zca

import (
	"context"
	"net/http"
)

type GetPendingGroupMembersUserInfo struct {
	UID        string `json:"uid"`
	Dpn        string `json:"dpn"`
	Avatar     string `json:"avatar"`
	UserSubmit any    `json:"user_submit"`
}

type GetPendingGroupMembersResponse struct {
	Time  int64                            `json:"time"`
	Users []GetPendingGroupMembersUserInfo `json:"users"`
}

// GetPendingGroupMembers lists pending members of a group (leader/deputy only).
func (a *API) GetPendingGroupMembers(ctx context.Context, groupID string) (*GetPendingGroupMembersResponse, error) {
	params := map[string]any{"grid": groupID, "imei": a.IMEI}
	return call[*GetPendingGroupMembersResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group")+"/api/group/pending-mems/list", nil, true), params)
}
