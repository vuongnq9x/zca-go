package zca

import (
	"context"
	"net/http"
)

type LogoutResponse = any

// Logout logs out of the Zalo account (Zalo Web logoutV2). Session cookies are
// invalid afterwards; the caller must stop any running Listener itself.
func (a *API) Logout(ctx context.Context) (LogoutResponse, error) {
	return call[LogoutResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL("https://wpa.chat.zalo.me/api/v2/login/logOut", nil, true),
		map[string]any{"time": nowMs(), "imei": a.IMEI, "computer_name": "Web", "language": a.Language})
}
