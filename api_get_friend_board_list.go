package zca

import (
	"context"
	"net/http"
)

type GetFriendBoardListResponse struct {
	Data    []string `json:"data"`
	Version int64    `json:"version"`
}

// GetFriendBoardList gets the friend board list of a conversation.
func (a *API) GetFriendBoardList(ctx context.Context, conversationID string) (*GetFriendBoardListResponse, error) {
	params := map[string]any{"conversationId": conversationID, "version": 0, "imei": a.IMEI}
	return call[*GetFriendBoardListResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("friend_board")+"/api/friendboard/list", nil, true), params)
}
