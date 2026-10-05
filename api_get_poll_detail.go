package zca

import (
	"context"
	"net/http"
)

type PollDetailResponse = PollDetail

// GetPollDetail gets poll detail.
func (a *API) GetPollDetail(ctx context.Context, pollID int64) (*PollDetailResponse, error) {
	if pollID == 0 {
		return nil, newError("Missing poll id")
	}
	params := map[string]any{"poll_id": pollID, "imei": a.IMEI}
	return call[*PollDetailResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("group")+"/api/poll/detail", nil, true), params)
}
