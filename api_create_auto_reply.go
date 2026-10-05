package zca

import (
	"context"
	"net/http"
)

type CreateAutoReplyPayload struct {
	Content   string
	IsEnable  bool
	StartTime int64
	EndTime   int64
	Scope     AutoReplyScope
	// Only used for AutoReplyScopeSpecificFriends and AutoReplyScopeFriendsExcept.
	UIDs []string
}

type CreateAutoReplyResponse struct {
	Item    AutoReplyItem `json:"item"`
	Version int64         `json:"version"`
}

// CreateAutoReply creates an auto reply (zBusiness).
func (a *API) CreateAutoReply(ctx context.Context, payload CreateAutoReplyPayload) (*CreateAutoReplyResponse, error) {
	uids := []string{}
	if (payload.Scope == AutoReplyScopeSpecificFriends || payload.Scope == AutoReplyScopeFriendsExcept) && payload.UIDs != nil {
		uids = payload.UIDs
	}
	params := map[string]any{
		"cliLang":    a.Language,
		"enable":     payload.IsEnable,
		"content":    payload.Content,
		"startTime":  payload.StartTime,
		"endTime":    payload.EndTime,
		"recurrence": []string{"RRULE:FREQ=DAILY;"},
		"scope":      payload.Scope,
		"uids":       uids,
	}
	return call[*CreateAutoReplyResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("auto_reply")+"/api/autoreply/create", nil, true), params)
}
