package zca

import (
	"context"
	"net/http"
)

type AddPollOptionsOption struct {
	Voted   bool   `json:"voted"`
	Content string `json:"content"`
}

type AddPollOptionsPayload struct {
	PollID         int64
	Options        []AddPollOptionsOption
	VotedOptionIDs []int64
}

type AddPollOptionsResponse struct {
	Options []PollOptions `json:"options"`
}

// AddPollOptions adds new options to a poll.
func (a *API) AddPollOptions(ctx context.Context, payload AddPollOptionsPayload) (*AddPollOptionsResponse, error) {
	params := map[string]any{
		"poll_id":          payload.PollID,
		"new_options":      mustJSON(payload.Options),
		"voted_option_ids": payload.VotedOptionIDs,
	}
	return call[*AddPollOptionsResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group")+"/api/poll/option/add", nil, true), params)
}
