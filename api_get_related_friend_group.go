package zca

import (
	"context"
	"net/http"
)

type GetRelatedFriendGroupResponse struct {
	GroupRelateds map[string][]string `json:"groupRelateds"`
}

// GetRelatedFriendGroup gets groups related to friends (zBusiness).
func (a *API) GetRelatedFriendGroup(ctx context.Context, friendIDs []string) (*GetRelatedFriendGroupResponse, error) {
	if friendIDs == nil {
		friendIDs = []string{}
	}
	params := map[string]any{"friend_ids": mustJSON(friendIDs), "imei": a.IMEI}
	return call[*GetRelatedFriendGroupResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("friend")+"/api/friend/group/related", nil, true), params)
}
