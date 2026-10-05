package zca

import (
	"context"
	"encoding/json"
	"net/http"
)

// ReuseAvatar sets an old avatar (photoID from GetAvatarList) as current. Zalo returns null.
func (a *API) ReuseAvatar(ctx context.Context, photoID string) error {
	params := map[string]any{"photoId": photoID, "isPostSocial": 0, "imei": a.IMEI}
	_, err := call[json.RawMessage](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("profile")+"/api/social/reuse-avatar", nil, true), params)
	return err
}
