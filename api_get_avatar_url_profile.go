package zca

import (
	"context"
	"net/http"
)

// GetAvatarUrlProfileResponse maps user id -> avatar.
type GetAvatarUrlProfileResponse map[string]struct {
	Avatar string `json:"avatar"`
}

// GetAvatarUrlProfile gets avatar urls of friends. avatarSize 0 means AvatarSizeLarge.
func (a *API) GetAvatarUrlProfile(ctx context.Context, friendIDs []string, avatarSize AvatarSize) (GetAvatarUrlProfileResponse, error) {
	if avatarSize == 0 {
		avatarSize = AvatarSizeLarge
	}
	params := map[string]any{"friend_ids": friendIDs, "avatar_size": avatarSize, "srcReq": -1}
	return call[GetAvatarUrlProfileResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("profile")+"/api/social/profile/avatar-url", nil, true), params)
}
