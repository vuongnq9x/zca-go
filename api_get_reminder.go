package zca

import (
	"context"
	"net/http"
)

type GetReminderResponse = ReminderGroup

// GetReminder gets reminder details from a group.
func (a *API) GetReminder(ctx context.Context, reminderID string) (*GetReminderResponse, error) {
	params := map[string]any{"eventId": reminderID, "imei": a.IMEI}
	return call[*GetReminderResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group_board")+"/api/board/topic/getReminder", nil, true), params)
}
