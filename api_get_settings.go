package zca

import (
	"context"
	"net/http"
)

type GetSettingsResponse = UserSetting

// GetSettings gets my account settings.
func (a *API) GetSettings(ctx context.Context) (*GetSettingsResponse, error) {
	return call[*GetSettingsResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL("https://wpa.chat.zalo.me/api/setting/me", nil, true), map[string]any{})
}
