package zca

import (
	"context"
	"net/http"
)

type GetPinConversationsResponse struct {
	Conversations []string `json:"conversations"`
	Version       int64    `json:"version"`
}

// GetPinConversations gets pinned conversations.
func (a *API) GetPinConversations(ctx context.Context) (*GetPinConversationsResponse, error) {
	return call[*GetPinConversationsResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("conversation")+"/api/pinconvers/list", nil, true), map[string]any{"imei": a.IMEI})
}
