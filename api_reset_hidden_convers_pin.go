package zca

import (
	"context"
	"net/http"
)

// ResetHiddenConversPin resets the hidden conversations PIN.
func (a *API) ResetHiddenConversPin(ctx context.Context) (string, error) {
	return call[string](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("conversation")+"/api/hiddenconvers/reset", nil, true), map[string]any{})
}
