package zca

import (
	"context"
	"net/http"
)

// RemoveReminder removes a reminder in a user or group thread. Result is "" or a number.
func (a *API) RemoveReminder(ctx context.Context, reminderID, threadID string, threadType ThreadType) (any, error) {
	path := "/api/board/oneone/remove"
	params := map[string]any{"uid": threadID, "reminderId": reminderID}
	if threadType == ThreadTypeGroup {
		path = "/api/board/topic/remove"
		params = map[string]any{"grid": threadID, "topicId": reminderID, "imei": a.IMEI}
	}
	return call[any](ctx, a.Session, http.MethodPost, a.MakeURL(a.svc("group_board")+path, nil, true), params)
}
