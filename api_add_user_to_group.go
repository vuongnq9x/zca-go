package zca

import (
	"context"
	"net/http"
)

type AddUserToGroupResponse struct {
	ErrorMembers []string            `json:"errorMembers"`
	ErrorData    map[string][]string `json:"error_data"`
}

// AddUserToGroup adds user(s) to an existing group.
func (a *API) AddUserToGroup(ctx context.Context, memberIDs []string, groupID string) (*AddUserToGroupResponse, error) {
	types := make([]int, len(memberIDs))
	for i := range types {
		types[i] = -1
	}
	params := map[string]any{
		"grid":        groupID,
		"members":     memberIDs,
		"memberTypes": types,
		"imei":        a.IMEI,
		"clientLang":  a.Language,
	}
	return call[*AddUserToGroupResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("group")+"/api/group/invite/v2", nil, true), params)
}
