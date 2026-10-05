package zca

import (
	"context"
	"net/http"
)

type GetArchivedChatListResponse struct {
	Items   []any `json:"items"`
	Version int64 `json:"version"`
}

// GetArchivedChatList gets the archived chat list.
func (a *API) GetArchivedChatList(ctx context.Context) (*GetArchivedChatListResponse, error) {
	params := map[string]any{"version": 1, "imei": a.IMEI}
	return call[*GetArchivedChatListResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("label")+"/api/archivedchat/list", nil, true), params)
}
