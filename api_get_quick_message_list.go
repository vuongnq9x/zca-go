package zca

import (
	"context"
	"net/http"
)

type GetQuickMessageListResponse struct {
	Cursor  int64          `json:"cursor"`
	Version int64          `json:"version"`
	Items   []QuickMessage `json:"items"`
}

// GetQuickMessageList gets the quick message list.
func (a *API) GetQuickMessageList(ctx context.Context) (*GetQuickMessageListResponse, error) {
	params := map[string]any{"version": 0, "lang": 0, "imei": a.IMEI}
	return call[*GetQuickMessageListResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("quick_message")+"/api/quickmessage/list", nil, true), params)
}
