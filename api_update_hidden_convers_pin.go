package zca

import (
	"context"
	"net/http"
	"regexp"
)

var updateHiddenConversPinRegex = regexp.MustCompile(`^\d{4}$`)

// UpdateHiddenConversPin sets the 4-digit hidden conversation pin.
func (a *API) UpdateHiddenConversPin(ctx context.Context, pin string) (string, error) {
	if !updateHiddenConversPinRegex.MatchString(pin) {
		return "", newError("Pin must be a 4-digit number between 0000-9999")
	}
	params := map[string]any{"new_pin": EncryptPin(pin), "imei": a.IMEI}
	return call[string](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("conversation")+"/api/hiddenconvers/update-pin", nil, true), params)
}
