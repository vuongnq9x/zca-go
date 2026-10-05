package zca

import (
	"context"
	"net/http"
)

type GetGroupBlockedMemberPayload struct {
	Page  int // default 1 when 0
	Count int // default 50 when 0
}

type GetGroupBlockedMemberResponse struct {
	BlockedMembers []GroupCurrentMem `json:"blocked_members"`
	HasMore        int64             `json:"has_more"`
}

// GetGroupBlockedMember gets blocked members of a group.
func (a *API) GetGroupBlockedMember(ctx context.Context, payload GetGroupBlockedMemberPayload, groupID string) (*GetGroupBlockedMemberResponse, error) {
	if payload.Page == 0 {
		payload.Page = 1
	}
	if payload.Count == 0 {
		payload.Count = 50
	}
	params := map[string]any{"grid": groupID, "page": payload.Page, "count": payload.Count, "imei": a.IMEI}
	return call[*GetGroupBlockedMemberResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group")+"/api/group/blockedmems/list", nil, true), params)
}
