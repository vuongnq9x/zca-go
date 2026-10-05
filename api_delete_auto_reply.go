package zca

import (
	"context"
	"net/http"
)

type DeleteAutoReplyResponse struct {
	Item    int64 `json:"item"`
	Version int64 `json:"version"`
}

// DeleteAutoReply deletes an auto reply (zBusiness).
func (a *API) DeleteAutoReply(ctx context.Context, id int64) (*DeleteAutoReplyResponse, error) {
	return call[*DeleteAutoReplyResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("auto_reply")+"/api/autoreply/delete", nil, true),
		map[string]any{"cliLang": a.Language, "id": id})
}
