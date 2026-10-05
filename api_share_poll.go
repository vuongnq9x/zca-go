package zca

import (
	"context"
	"net/http"
)

// SharePoll shares a poll.
func (a *API) SharePoll(ctx context.Context, pollID int64) (string, error) {
	return call[string](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("group")+"/api/poll/share", nil, true), map[string]any{"poll_id": pollID, "imei": a.IMEI})
}
