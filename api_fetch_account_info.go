package zca

import (
	"context"
	"net/http"
)

type FetchAccountInfoResponse struct {
	Profile User `json:"profile"`
}

// FetchAccountInfo gets the logged-in account's profile.
func (a *API) FetchAccountInfo(ctx context.Context) (*FetchAccountInfoResponse, error) {
	resp, err := a.Request(ctx, http.MethodGet, a.MakeURL(a.svc("profile")+"/api/social/profile/me-v2", nil, true), nil, nil)
	if err != nil {
		return nil, err
	}
	raw, err := a.Resolve(resp, true)
	if err != nil {
		return nil, err
	}
	return decodeData[*FetchAccountInfoResponse](raw)
}
