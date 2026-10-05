package zca

import (
	"context"
	"net/http"
)

type VotePollResponse struct {
	Options []PollOptions `json:"options"`
}

// VotePoll votes on a poll; empty optionIDs removes the vote.
func (a *API) VotePoll(ctx context.Context, pollID int64, optionIDs []int64) (*VotePollResponse, error) {
	if optionIDs == nil {
		optionIDs = []int64{}
	}
	return call[*VotePollResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group")+"/api/poll/vote", nil, true),
		map[string]any{"poll_id": pollID, "option_ids": optionIDs, "imei": a.IMEI})
}
