package zca

import (
	"context"
	"net/http"
)

type DeleteAvatarResponse struct {
	DelPhotoIDs []string `json:"delPhotoIds"`
	ErrMap      map[string]struct {
		Err int64 `json:"err"`
	} `json:"errMap"`
}

// DeleteAvatar deletes avatar(s) from the avatar list.
func (a *API) DeleteAvatar(ctx context.Context, photoIDs []string) (*DeleteAvatarResponse, error) {
	del := make([]map[string]string, len(photoIDs))
	for i, id := range photoIDs {
		del[i] = map[string]string{"photoId": id}
	}
	return call[*DeleteAvatarResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("profile")+"/api/social/del-avatars", nil, true),
		map[string]any{"delPhotos": mustJSON(del), "imei": a.IMEI})
}
