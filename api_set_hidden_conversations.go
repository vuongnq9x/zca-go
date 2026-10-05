package zca

import (
	"context"
	"net/http"
)

// SetHiddenConversations hides or unhides conversations.
func (a *API) SetHiddenConversations(ctx context.Context, hidden bool, threadIDs []string, threadType ThreadType) (string, error) {
	if len(threadIDs) == 0 {
		return "", newError("threadId is required")
	}
	isGroup := 0
	if threadType == ThreadTypeGroup {
		isGroup = 1
	}
	type thread struct {
		ThreadID string `json:"thread_id"`
		IsGroup  int    `json:"is_group"`
	}
	threads := make([]thread, len(threadIDs))
	for i, id := range threadIDs {
		threads[i] = thread{id, isGroup}
	}
	set, other := "add_threads", "del_threads"
	if !hidden {
		set, other = other, set
	}
	params := map[string]any{set: mustJSON(threads), other: "[]", "imei": a.IMEI}
	return call[string](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("conversation")+"/api/hiddenconvers/add-remove", nil, true), params)
}
