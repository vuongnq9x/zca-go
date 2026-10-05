package zca

import (
	"context"
	"net/http"
	"net/url"
)

// UpdateProfileBio updates the profile bio (status). Params go both in the query and the body.
func (a *API) UpdateProfileBio(ctx context.Context, status string) (string, error) {
	enc, err := a.EncodeAES(mustJSON(map[string]any{"status": status}))
	if err != nil {
		return "", newError("Failed to encrypt params")
	}
	u := a.MakeURL(a.MakeURL(a.svc("profile")+"/api/social/profile/status", nil, true), map[string]any{"params": enc}, true)
	resp, err := a.Request(ctx, http.MethodPost, u, []byte(url.Values{"params": {enc}}.Encode()), nil)
	if err != nil {
		return "", err
	}
	raw, err := a.Resolve(resp, true)
	if err != nil {
		return "", err
	}
	return decodeData[string](raw)
}
