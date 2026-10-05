package zca

import (
	"context"
	"net/http"
)

type GetFullAvatarResponse struct {
	BkFullAvatar string `json:"bk_full_avatar"`
	FullAvatar   string `json:"full_avatar"`
}

// GetFullAvatar gets the full-size avatar of a user.
func (a *API) GetFullAvatar(ctx context.Context, friendID string) (*GetFullAvatarResponse, error) {
	return call[*GetFullAvatarResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("profile")+"/api/social/profile/avatar", nil, true), map[string]any{"fid": friendID, "imei": a.IMEI})
}
