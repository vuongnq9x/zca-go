package zca

import (
	"context"
	"net/http"
)

type GetGroupLinkDetailResponse struct {
	Link           string `json:"link,omitempty"`
	ExpirationDate int64  `json:"expiration_date,omitempty"`
	// Enabled: 1 enabled, 0 disabled.
	Enabled int64 `json:"enabled"`
}

// GetGroupLinkDetail gets the invite link detail of a group.
func (a *API) GetGroupLinkDetail(ctx context.Context, groupID string) (*GetGroupLinkDetailResponse, error) {
	return call[*GetGroupLinkDetailResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group")+"/api/group/link/detail", nil, true), map[string]any{"grid": groupID, "imei": a.IMEI})
}
