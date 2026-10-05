package zca

import (
	"context"
	"net/http"
)

// ChatTTL is the auto-delete time to live in milliseconds.
type ChatTTL int64

const (
	ChatTTLNoDelete     ChatTTL = 0
	ChatTTLOneDay       ChatTTL = 86400000
	ChatTTLSevenDays    ChatTTL = 604800000
	ChatTTLFourteenDays ChatTTL = 1209600000
)

// UpdateAutoDeleteChat sets the auto-delete TTL of a thread.
func (a *API) UpdateAutoDeleteChat(ctx context.Context, ttl ChatTTL, threadID string, threadType ThreadType) (string, error) {
	isGroup := 0
	if threadType == ThreadTypeGroup {
		isGroup = 1
	}
	params := map[string]any{"threadId": threadID, "isGroup": isGroup, "ttl": ttl, "clientLang": a.Language}
	return call[string](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("conversation")+"/api/conv/autodelete/updateConvers", nil, true), params)
}
