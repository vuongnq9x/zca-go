package zca

import (
	"context"
	"net/http"
)

type GetAvatarListResponse struct {
	AlbumID     string `json:"albumId"`
	NextPhotoID string `json:"nextPhotoId"`
	HasMore     int64  `json:"hasMore"`
	Photos      []struct {
		PhotoID   string `json:"photoId"`
		Thumbnail string `json:"thumbnail"`
		URL       string `json:"url"`
		BkURL     string `json:"bkUrl"`
	} `json:"photos"`
}

// GetAvatarList gets the avatar list. count 0 means 50, page 0 means 1.
func (a *API) GetAvatarList(ctx context.Context, count, page int) (*GetAvatarListResponse, error) {
	if count == 0 {
		count = 50
	}
	if page == 0 {
		page = 1
	}
	params := map[string]any{"page": page, "albumId": "0", "count": count, "imei": a.IMEI}
	return call[*GetAvatarListResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("profile")+"/api/social/avatar-list", nil, true), params)
}
