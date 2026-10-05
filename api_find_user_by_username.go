package zca

import (
	"context"
	"net/http"
)

// FindUserByUsername finds a user by username. avatarSize 0 means AvatarSizeLarge.
func (a *API) FindUserByUsername(ctx context.Context, username string, avatarSize AvatarSize) (*UserBasic, error) {
	if avatarSize == 0 {
		avatarSize = AvatarSizeLarge
	}
	params := map[string]any{"user_name": username, "avatar_size": avatarSize}
	return call[*UserBasic](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("friend")+"/api/friend/search/by-user-name", nil, true), params)
}
