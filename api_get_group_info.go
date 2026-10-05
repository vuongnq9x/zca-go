package zca

import (
	"context"
	"net/http"
)

type GroupInfoPendingApprove struct {
	Time int64    `json:"time"`
	UIDs []string `json:"uids"`
}

// GroupInfoResponseGroup is TS GroupInfo & { memVerList, pendingApprove }.
type GroupInfoResponseGroup struct {
	GroupInfo
	MemVerList     []string                `json:"memVerList"`
	PendingApprove GroupInfoPendingApprove `json:"pendingApprove"`
}

type GroupInfoResponse struct {
	RemovedsGroup   []string                          `json:"removedsGroup"`
	UnchangedsGroup []string                          `json:"unchangedsGroup"`
	GridInfoMap     map[string]GroupInfoResponseGroup `json:"gridInfoMap"`
}

// GetGroupInfo gets information of one or more groups.
func (a *API) GetGroupInfo(ctx context.Context, groupIDs []string) (*GroupInfoResponse, error) {
	verMap := make(map[string]int, len(groupIDs))
	for _, id := range groupIDs {
		verMap[id] = 0
	}
	return call[*GroupInfoResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("group")+"/api/group/getmg-v2", nil, true), map[string]any{"gridVerMap": mustJSON(verMap)})
}
