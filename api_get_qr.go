package zca

import (
	"context"
	"net/http"
)

// GetQRResponse maps user ID -> QR code.
type GetQRResponse map[string]string

// GetQR gets QR codes for users.
func (a *API) GetQR(ctx context.Context, userIDs []string) (GetQRResponse, error) {
	return call[GetQRResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("friend")+"/api/friend/mget-qr", nil, true), map[string]any{"fids": userIDs})
}
