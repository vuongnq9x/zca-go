package zca

import (
	"context"
	"net/http"
	"strconv"
)

type ChangeGroupNameResponse struct {
	Status int64 `json:"status"`
}

// ChangeGroupName renames a group. An empty name becomes the current timestamp.
func (a *API) ChangeGroupName(ctx context.Context, name, groupID string) (*ChangeGroupNameResponse, error) {
	if name == "" {
		name = strconv.FormatInt(nowMs(), 10)
	}
	return call[*ChangeGroupNameResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("group")+"/api/group/updateinfo", nil, true),
		map[string]any{"grid": groupID, "gname": name, "imei": a.IMEI})
}
