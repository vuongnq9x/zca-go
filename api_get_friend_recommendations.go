package zca

import (
	"context"
	"net/http"
)

type FriendRecommendationsType int

const (
	FriendRecommendationsTypeRecommendedFriend     FriendRecommendationsType = 1
	FriendRecommendationsTypeReceivedFriendRequest FriendRecommendationsType = 2
)

type FriendRecommendationsCollapseMsgListConfig struct {
	CollapseID    int64 `json:"collapseId"`
	CollapseXItem int64 `json:"collapseXItem"`
	CollapseYItem int64 `json:"collapseYItem"`
}

type FriendRecommendationsDataInfo struct {
	UserID      string                    `json:"userId"`
	ZaloName    string                    `json:"zaloName"`
	DisplayName string                    `json:"displayName"`
	Avatar      string                    `json:"avatar"`
	PhoneNumber string                    `json:"phoneNumber"`
	Status      string                    `json:"status"`
	Gender      Gender                    `json:"gender"`
	Dob         int64                     `json:"dob"`
	Type        int64                     `json:"type"`
	RecommType  FriendRecommendationsType `json:"recommType"`
	RecommSrc   int64                     `json:"recommSrc"`
	RecommTime  int64                     `json:"recommTime"`
	RecommInfo  struct {
		SuggestWay int64   `json:"suggestWay"`
		Source     int64   `json:"source"`
		Message    string  `json:"message"`
		CustomText *string `json:"customText"`
	} `json:"recommInfo"`
	BizPkg          ZBusinessPackage `json:"bizPkg"`
	IsSeenFriendReq bool             `json:"isSeenFriendReq"`
}

type FriendRecommendationsRecommItem struct {
	RecommItemType int64                         `json:"recommItemType"`
	DataInfo       FriendRecommendationsDataInfo `json:"dataInfo"`
}

type GetFriendRecommendationsResponse struct {
	ExpiredDuration       int64                                      `json:"expiredDuration"`
	CollapseMsgListConfig FriendRecommendationsCollapseMsgListConfig `json:"collapseMsgListConfig"`
	RecommItems           []FriendRecommendationsRecommItem          `json:"recommItems"`
}

// GetFriendRecommendations gets friend recommendations and received friend requests.
func (a *API) GetFriendRecommendations(ctx context.Context) (*GetFriendRecommendationsResponse, error) {
	return call[*GetFriendRecommendationsResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("friend")+"/api/friend/recommendsv2/list", nil, true), map[string]any{"imei": a.IMEI})
}
