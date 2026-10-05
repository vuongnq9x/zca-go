package zca

import (
	"context"
	"net/http"
)

type LostFocusResponse struct {
	Status bool `json:"status"`
}

// LostFocus tells Zalo the client went to background.
func (a *API) LostFocus(ctx context.Context) (*LostFocusResponse, error) {
	return call[*LostFocusResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("profile")+"/api/social/profile/changefgtobg", nil, true), map[string]any{})
}
