package zca

import (
	"context"
	"net/http"
	"strconv"
)

type DeleteChatResponse struct {
	Status int64 `json:"status"`
}

// DeleteChatLastMessage is the last message to delete backwards from.
type DeleteChatLastMessage struct {
	OwnerID     string `json:"ownerId"`
	CliMsgID    string `json:"cliMsgId"`
	GlobalMsgID string `json:"globalMsgId"`
}

// DeleteChat deletes a conversation (for you only).
func (a *API) DeleteChat(ctx context.Context, lastMessage DeleteChatLastMessage, threadID string, threadType ThreadType) (*DeleteChatResponse, error) {
	base, key := a.svc("chat")+"/api/message/deleteconver", "toid"
	if threadType == ThreadTypeGroup {
		base, key = a.svc("group")+"/api/group/deleteconver", "grid"
	}
	params := map[string]any{
		key:        threadID,
		"cliMsgId": strconv.FormatInt(nowMs(), 10),
		"conver":   lastMessage,
		"onlyMe":   1,
		"imei":     a.IMEI,
	}
	return call[*DeleteChatResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(base, map[string]any{"nretry": 0}, true), params)
}
