package zca

import (
	"context"
	"net/http"
)

type UpdateActiveStatusResponse struct {
	Status bool `json:"status"`
}

// UpdateActiveStatus pings (active) or deactivates the online status.
func (a *API) UpdateActiveStatus(ctx context.Context, active bool) (*UpdateActiveStatusResponse, error) {
	status, path := 0, "/api/social/profile/deactive"
	if active {
		status, path = 1, "/api/social/profile/ping"
	}
	return call[*UpdateActiveStatusResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("profile")+path, nil, true), map[string]any{"status": status, "imei": a.IMEI})
}
