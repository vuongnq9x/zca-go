package zca

import (
	"context"
	"net/http"
)

type CreatePollOptions struct {
	Question string
	Options  []string
	// Poll expiration time in milliseconds (0 = no expiration).
	ExpiredTime       int64
	AllowMultiChoices bool
	AllowAddNewOption bool
	// Hides voting results until the user has voted.
	HideVotePreview bool
	// Hides poll voters (anonymous poll).
	IsAnonymous bool
}

type CreatePollResponse = PollDetail

// CreatePoll creates a poll in a group.
func (a *API) CreatePoll(ctx context.Context, options CreatePollOptions, groupID string) (*CreatePollResponse, error) {
	params := map[string]any{
		"group_id":             groupID,
		"question":             options.Question,
		"options":              options.Options,
		"expired_time":         options.ExpiredTime,
		"pinAct":               false,
		"allow_multi_choices":  options.AllowMultiChoices,
		"allow_add_new_option": options.AllowAddNewOption,
		"is_hide_vote_preview": options.HideVotePreview,
		"is_anonymous":         options.IsAnonymous,
		"poll_type":            0,
		"src":                  1,
		"imei":                 a.IMEI,
	}
	return call[*CreatePollResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("group")+"/api/poll/create", nil, true), params)
}
