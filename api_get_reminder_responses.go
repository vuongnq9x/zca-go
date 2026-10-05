package zca

import (
	"context"
	"net/http"
)

type GetReminderResponsesResponse struct {
	RejectMember []string `json:"rejectMember"`
	AcceptMember []string `json:"acceptMember"`
}

// GetReminderResponses gets accept/reject responses of a reminder.
func (a *API) GetReminderResponses(ctx context.Context, reminderID string) (*GetReminderResponsesResponse, error) {
	return call[*GetReminderResponsesResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group_board")+"/api/board/topic/listResponseEvent", nil, true), map[string]any{"eventId": reminderID})
}
