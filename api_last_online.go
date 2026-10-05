package zca

import (
	"context"
	"net/http"
)

type LastOnlineResponse struct {
	Settings struct {
		ShowOnlineStatus bool `json:"show_online_status"`
	} `json:"settings"`
	LastOnline int64 `json:"lastOnline"`
}

// LastOnline gets a user's last online time.
func (a *API) LastOnline(ctx context.Context, uid string) (*LastOnlineResponse, error) {
	params := map[string]any{"uid": uid, "conv_type": 1, "imei": a.IMEI}
	return call[*LastOnlineResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("profile")+"/api/social/profile/lastOnline", nil, true), params)
}
