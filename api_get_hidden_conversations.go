package zca

import (
	"context"
	"net/http"
)

type GetHiddenConversationsResponse struct {
	Pin     string `json:"pin"`
	Threads []struct {
		// IsGroup: 1 true, 0 false.
		IsGroup  int64  `json:"is_group"`
		ThreadID string `json:"thread_id"`
	} `json:"threads"`
}

// GetHiddenConversations gets hidden conversations.
func (a *API) GetHiddenConversations(ctx context.Context) (*GetHiddenConversationsResponse, error) {
	return call[*GetHiddenConversationsResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("conversation")+"/api/hiddenconvers/get-all", nil, true), map[string]any{"imei": a.IMEI})
}
