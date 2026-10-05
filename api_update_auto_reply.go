package zca

import (
	"context"
	"net/http"
)

type UpdateAutoReplyPayload struct {
	ID        int64
	Content   string
	IsEnable  bool
	StartTime int64
	EndTime   int64
	Scope     AutoReplyScope
	UIDs      []string // used only with AutoReplyScopeSpecificFriends / AutoReplyScopeFriendsExcept
}

type UpdateAutoReplyResponse struct {
	Item    AutoReplyItem `json:"item"`
	Version int64         `json:"version"`
}

// UpdateAutoReply updates an auto reply (zBusiness).
func (a *API) UpdateAutoReply(ctx context.Context, payload UpdateAutoReplyPayload) (*UpdateAutoReplyResponse, error) {
	uids := []string{}
	if payload.Scope == AutoReplyScopeSpecificFriends || payload.Scope == AutoReplyScopeFriendsExcept {
		uids = append(uids, payload.UIDs...)
	}
	params := map[string]any{
		"cliLang": a.Language, "id": payload.ID, "enable": payload.IsEnable, "content": payload.Content,
		"startTime": payload.StartTime, "endTime": payload.EndTime,
		"recurrence": []string{"RRULE:FREQ=DAILY;"}, "scope": payload.Scope, "uids": uids,
	}
	return call[*UpdateAutoReplyResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("auto_reply")+"/api/autoreply/update", nil, true), params)
}
