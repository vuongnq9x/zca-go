package zca

import (
	"context"
	"net/http"
)

type GetFriendRequestStatusResponse struct {
	AddFriendPrivacy int64 `json:"addFriendPrivacy"`
	IsSeenFriendReq  bool  `json:"isSeenFriendReq"`
	IsFriend         int64 `json:"is_friend"`
	IsRequested      int64 `json:"is_requested"`
	IsRequesting     int64 `json:"is_requesting"`
}

// GetFriendRequestStatus gets the friend request status with a user.
func (a *API) GetFriendRequestStatus(ctx context.Context, friendID string) (*GetFriendRequestStatusResponse, error) {
	return call[*GetFriendRequestStatusResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("friend")+"/api/friend/reqstatus", nil, true), map[string]any{"fid": friendID, "imei": a.IMEI})
}
