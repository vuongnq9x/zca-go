package zca

import (
	"context"
	"net/http"
)

type GetAutoDeleteChatResponse struct {
	Convers []struct {
		DestID    string `json:"destId"`
		IsGroup   bool   `json:"isGroup"`
		TTL       int64  `json:"ttl"`
		CreatedAt int64  `json:"createdAt"`
	} `json:"convers"`
}

// GetAutoDeleteChat gets conversations with auto-delete enabled.
func (a *API) GetAutoDeleteChat(ctx context.Context) (*GetAutoDeleteChatResponse, error) {
	return call[*GetAutoDeleteChatResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("conversation")+"/api/conv/autodelete/getConvers", nil, true), map[string]any{})
}
