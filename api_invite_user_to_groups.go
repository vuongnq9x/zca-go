package zca

import (
	"context"
	"net/http"
)

type InviteUserToGroupsResponse struct {
	GridMessageMap map[string]struct {
		ErrorCode    int64   `json:"error_code"`
		ErrorMessage string  `json:"error_message"`
		Data         *string `json:"data"`
	} `json:"grid_message_map"`
}

// InviteUserToGroups invites a user to one or more groups.
func (a *API) InviteUserToGroups(ctx context.Context, userID string, groupIDs []string) (*InviteUserToGroupsResponse, error) {
	params := map[string]any{
		"grids":          groupIDs,
		"member":         userID,
		"memberType":     -1,
		"srcInteraction": 2,
		"clientLang":     a.Language,
	}
	return call[*InviteUserToGroupsResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group")+"/api/group/invite/multi", nil, true), params)
}
