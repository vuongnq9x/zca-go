package zca

import (
	"context"
	"net/http"
)

// SetPinnedConversations pins or unpins conversations.
func (a *API) SetPinnedConversations(ctx context.Context, pinned bool, threadIDs []string, threadType ThreadType) (string, error) {
	prefix := "u"
	if threadType == ThreadTypeGroup {
		prefix = "g"
	}
	convs := make([]string, len(threadIDs))
	for i, id := range threadIDs {
		convs[i] = prefix + id
	}
	action := 2
	if pinned {
		action = 1
	}
	return call[string](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("conversation")+"/api/pinconvers/updatev2", nil, true),
		map[string]any{"actionType": action, "conversations": convs})
}
