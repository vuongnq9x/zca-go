package zca

import (
	"context"
	"net/http"
)

// LockPoll locks a poll, preventing further votes.
func (a *API) LockPoll(ctx context.Context, pollID int64) (string, error) {
	return call[string](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("group")+"/api/poll/end", nil, true), map[string]any{"poll_id": pollID, "imei": a.IMEI})
}
