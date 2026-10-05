package zca

import (
	"context"
	"net/http"
)

// JoinGroupLink joins a group via invite link. Zalo may return code 240 (approval required) or 178 (already a member).
func (a *API) JoinGroupLink(ctx context.Context, link string) (string, error) {
	return call[string](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group")+"/api/group/link/join", nil, true), map[string]any{"link": link, "clientLang": a.Language})
}
