package zca

import (
	"context"
	"net/http"
)

type MiniProfile struct {
	UserID   string           `json:"userId"`
	Avatar   string           `json:"avatar"`
	GlobalID string           `json:"globalId"`
	Type     int64            `json:"type"`
	BizPkg   ZBusinessPackage `json:"bizPkg"`
	OAInfo   any              `json:"oaInfo"`
	ZaloName string           `json:"zaloName"`
}

// GetMiniProfilesResponse maps user id to its mini profile.
type GetMiniProfilesResponse = map[string]MiniProfile

// GetMiniProfiles gets mini profiles of users. avatarSize 0 means AvatarSizeMedium.
func (a *API) GetMiniProfiles(ctx context.Context, userIDs []string, avatarSize AvatarSize) (GetMiniProfilesResponse, error) {
	if len(userIDs) == 0 {
		return nil, newError("Missing user id")
	}
	if avatarSize == 0 {
		avatarSize = AvatarSizeMedium
	}
	return call[GetMiniProfilesResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("profile")+"/api/social/friend/getminiprofiles", nil, true),
		map[string]any{"friend_ids": userIDs, "avatar_size": avatarSize, "incInvalid": 1, "srcReq": 5})
}
