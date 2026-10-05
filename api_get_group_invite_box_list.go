package zca

import (
	"context"
	"net/http"
)

type GetGroupInviteBoxListPayload struct {
	MPage      int // default 1 when 0
	Page       int // default 0
	InvPerPage int // default 12 when 0
	MCount     int // default 10 when 0
}

type GetGroupInviteBoxListResponse struct {
	Invitations []struct {
		GroupInfo     GroupInfo       `json:"groupInfo"`
		InviterInfo   GroupCurrentMem `json:"inviterInfo"`
		GrCreatorInfo GroupCurrentMem `json:"grCreatorInfo"`
		// ExpiredTs is the expiry timestamp (max 7 days).
		ExpiredTs string `json:"expiredTs"`
		Type      int64  `json:"type"`
	} `json:"invitations"`
	Total   int64 `json:"total"`
	HasMore bool  `json:"hasMore"`
}

// GetGroupInviteBoxList gets the list of pending group invitations.
func (a *API) GetGroupInviteBoxList(ctx context.Context, payload GetGroupInviteBoxListPayload) (*GetGroupInviteBoxListResponse, error) {
	if payload.MPage == 0 {
		payload.MPage = 1
	}
	if payload.InvPerPage == 0 {
		payload.InvPerPage = 12
	}
	if payload.MCount == 0 {
		payload.MCount = 10
	}
	params := map[string]any{
		"mpage":              payload.MPage,
		"page":               payload.Page,
		"invPerPage":         payload.InvPerPage,
		"mcount":             payload.MCount,
		"lastGroupId":        nil,
		"avatar_size":        120,
		"member_avatar_size": 120,
	}
	return call[*GetGroupInviteBoxListResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group")+"/api/group/inv-box/list", nil, true), params)
}
