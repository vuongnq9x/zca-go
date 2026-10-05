package zca

import (
	"context"
	"net/http"
)

type RemoveQuickMessageResponse struct {
	ItemIDs []int64 `json:"itemIds"`
	Version int64   `json:"version"`
}

// RemoveQuickMessage removes quick message(s). Zalo may return code 212 if an item does not exist.
func (a *API) RemoveQuickMessage(ctx context.Context, itemIDs []int64) (*RemoveQuickMessageResponse, error) {
	return call[*RemoveQuickMessageResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("quick_message")+"/api/quickmessage/delete", nil, true), map[string]any{"itemIds": itemIDs})
}
