package zca

import (
	"context"
	"net/http"
)

type GetGroupLinkInfoPayload struct {
	Link       string
	MemberPage int // default 1 when 0
}

type GetGroupLinkInfoResponse struct {
	GroupID       string            `json:"groupId"`
	Name          string            `json:"name"`
	Desc          string            `json:"desc"`
	Type          GroupType         `json:"type"`
	CreatorID     string            `json:"creatorId"`
	Avt           string            `json:"avt"`
	FullAvt       string            `json:"fullAvt"`
	AdminIDs      []string          `json:"adminIds"`
	CurrentMems   []GroupCurrentMem `json:"currentMems"`
	Admins        []any             `json:"admins"`
	HasMoreMember int64             `json:"hasMoreMember"`
	SubType       int64             `json:"subType"`
	TotalMember   int64             `json:"totalMember"`
	Setting       GroupSetting      `json:"setting"`
	GlobalID      string            `json:"globalId"`
}

// GetGroupLinkInfo gets group info from an invite link.
func (a *API) GetGroupLinkInfo(ctx context.Context, payload GetGroupLinkInfoPayload) (*GetGroupLinkInfoResponse, error) {
	if payload.MemberPage == 0 {
		payload.MemberPage = 1
	}
	params := map[string]any{
		"link":               payload.Link,
		"avatar_size":        120,
		"member_avatar_size": 120,
		"mpage":              payload.MemberPage,
	}
	return call[*GetGroupLinkInfoResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group")+"/api/group/link/ginfo", nil, true), params)
}
