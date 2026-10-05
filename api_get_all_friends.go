package zca

import (
	"context"
	"net/http"
)

// GetAllFriends gets all friends. count 0 means 20000, page 0 means 1, avatarSize 0 means AvatarSizeSmall.
func (a *API) GetAllFriends(ctx context.Context, count, page int, avatarSize AvatarSize) ([]User, error) {
	if count == 0 {
		count = 20000
	}
	if page == 0 {
		page = 1
	}
	if avatarSize == 0 {
		avatarSize = AvatarSizeSmall
	}
	params := map[string]any{
		"incInvalid":  1,
		"page":        page,
		"count":       count,
		"avatar_size": avatarSize,
		"actiontime":  0,
		"imei":        a.IMEI,
	}
	return call[[]User](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("profile")+"/api/social/friend/getfriends", nil, true), params)
}
