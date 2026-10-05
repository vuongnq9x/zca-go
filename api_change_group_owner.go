package zca

import (
	"context"
	"net/http"
)

type ChangeGroupOwnerResponse struct {
	Time int64 `json:"time"`
}

// ChangeGroupOwner transfers group ownership. You lose admin rights afterwards.
func (a *API) ChangeGroupOwner(ctx context.Context, memberID, groupID string) (*ChangeGroupOwnerResponse, error) {
	return call[*ChangeGroupOwnerResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group")+"/api/group/change-owner", nil, true),
		map[string]any{"grid": groupID, "newAdminId": memberID, "imei": a.IMEI, "language": a.Language})
}
