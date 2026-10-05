package zca

import (
	"context"
	"net/http"
)

type GetAutoReplyListResponse struct {
	Item    []AutoReplyItem `json:"item"`
	Version int64           `json:"version"`
}

// GetAutoReplyList gets the auto reply list (zBusiness).
func (a *API) GetAutoReplyList(ctx context.Context) (*GetAutoReplyListResponse, error) {
	params := map[string]any{"version": 0, "cliLang": a.Language}
	return call[*GetAutoReplyListResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("auto_reply")+"/api/autoreply/list", nil, true), params)
}
