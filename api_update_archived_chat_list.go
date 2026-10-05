package zca

import (
	"context"
	"net/http"
)

type UpdateArchivedChatListTarget struct {
	ID   string     `json:"id"`
	Type ThreadType `json:"type"`
}

type UpdateArchivedChatListResponse struct {
	NeedResync bool  `json:"needResync"`
	Version    int64 `json:"version"`
}

// UpdateArchivedChatList archives (true) or unarchives conversations.
func (a *API) UpdateArchivedChatList(ctx context.Context, isArchived bool, conversations []UpdateArchivedChatListTarget) (*UpdateArchivedChatListResponse, error) {
	action := 1
	if isArchived {
		action = 0
	}
	params := map[string]any{"actionType": action, "ids": conversations, "imei": a.IMEI, "version": nowMs()}
	return call[*UpdateArchivedChatListResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("label")+"/api/archivedchat/update", nil, true), params)
}
