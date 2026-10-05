package zca

import (
	"context"
	"net/http"
)

type EnableGroupLinkResponse struct {
	Link           string `json:"link"`
	ExpirationDate int64  `json:"expiration_date"`
	Enabled        int64  `json:"enabled"`
}

// EnableGroupLink enables and creates a new group invite link.
func (a *API) EnableGroupLink(ctx context.Context, groupID string) (*EnableGroupLinkResponse, error) {
	return call[*EnableGroupLinkResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group")+"/api/group/link/new", nil, true),
		map[string]any{"grid": groupID, "imei": a.IMEI})
}
