package zca

import (
	"context"
	"net/http"
	"strings"
)

type GroupMemberProfile struct {
	DisplayName    string `json:"displayName"`
	ZaloName       string `json:"zaloName"`
	Avatar         string `json:"avatar"`
	AccountStatus  int64  `json:"accountStatus"`
	Type           int64  `json:"type"`
	LastUpdateTime int64  `json:"lastUpdateTime"`
	GlobalID       string `json:"globalId"`
	ID             string `json:"id"`
}

type GetGroupMembersInfoResponse struct {
	Profiles          map[string]GroupMemberProfile `json:"profiles"`
	UnchangedsProfile []any                         `json:"unchangeds_profile"`
}

// GetGroupMembersInfo gets profiles of group members.
func (a *API) GetGroupMembersInfo(ctx context.Context, memberIDs []string) (*GetGroupMembersInfoResponse, error) {
	ids := make([]string, len(memberIDs))
	for i, id := range memberIDs {
		if !strings.HasSuffix(id, "_0") {
			id += "_0"
		}
		ids[i] = id
	}
	return call[*GetGroupMembersInfoResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("profile")+"/api/social/group/members", nil, true), map[string]any{"friend_pversion_map": ids})
}
