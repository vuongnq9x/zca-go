package zca

import (
	"context"
	"net/http"
)

type ExtraInfo struct {
	Purpose string   `json:"purpose"`
	Role    []string `json:"role"`
}

// GetExtraFriendInfoResponse maps friend id to its extra info.
type GetExtraFriendInfoResponse = map[string]ExtraInfo

// GetExtraFriendInfo gets extra info of friends.
func (a *API) GetExtraFriendInfo(ctx context.Context, friendIDs []string) (GetExtraFriendInfoResponse, error) {
	if len(friendIDs) == 0 {
		return nil, newError("Missing friend id")
	}
	return call[GetExtraFriendInfoResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("friend")+"/api/friend/context/multiget", nil, true),
		map[string]any{"fids": friendIDs, "imei": a.IMEI})
}
