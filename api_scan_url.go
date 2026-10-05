package zca

import (
	"context"
	"net/http"
)

type ScanURLResponse struct {
	IsSafe bool `json:"isSafe"`
}

// ScanURL checks whether a URL is safe. Zalo error 114 means invalid params.
func (a *API) ScanURL(ctx context.Context, u string) (*ScanURLResponse, error) {
	return call[*ScanURLResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("file")+"/api/message/scanurl", nil, true), map[string]any{"url": u})
}
