package zca

import (
	"context"
	"net/http"
)

type SentFriendRequestInfo struct {
	UserID      string           `json:"userId"`
	ZaloName    string           `json:"zaloName"`
	DisplayName string           `json:"displayName"`
	Avatar      string           `json:"avatar"`
	GlobalID    string           `json:"globalId"`
	BizPkg      ZBusinessPackage `json:"bizPkg"`
	FReqInfo    struct {
		Message string `json:"message"`
		Src     int64  `json:"src"`
		Time    int64  `json:"time"`
	} `json:"fReqInfo"`
}

// GetSentFriendRequestResponse maps user ID -> request info.
type GetSentFriendRequestResponse map[string]SentFriendRequestInfo

// GetSentFriendRequest lists sent friend requests. Zalo may return error code 112 when there are none.
func (a *API) GetSentFriendRequest(ctx context.Context) (GetSentFriendRequestResponse, error) {
	return call[GetSentFriendRequestResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("friend")+"/api/friend/requested/list", nil, true), map[string]any{"imei": a.IMEI})
}
