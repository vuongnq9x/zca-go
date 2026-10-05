package zca

import (
	"context"
	"net/http"
)

// GetExtraGroupInfoResponse maps group id to its extra info.
type GetExtraGroupInfoResponse = map[string]ExtraInfo

// GetExtraGroupInfo gets extra info of groups.
func (a *API) GetExtraGroupInfo(ctx context.Context, groupIDs []string) (GetExtraGroupInfoResponse, error) {
	if len(groupIDs) == 0 {
		return nil, newError("Missing group id")
	}
	return call[GetExtraGroupInfoResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("group")+"/api/group/userctx/multiget", nil, true),
		map[string]any{"grids": groupIDs, "imei": a.IMEI})
}
