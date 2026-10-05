package zca

import (
	"context"
	"net/http"
)

// DecryptGlobalIDResponse maps global id to user id.
type DecryptGlobalIDResponse = map[string]string

// DecryptGlobalID decrypts global ids into user ids.
func (a *API) DecryptGlobalID(ctx context.Context, globalIDs []string) (DecryptGlobalIDResponse, error) {
	if len(globalIDs) == 0 {
		return nil, newError("Missing global id")
	}
	// globalUids is a JSON-encoded string, not an array.
	return call[DecryptGlobalIDResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("profile")+"/api/gid/decrypt", nil, true), map[string]any{"globalUids": mustJSON(globalIDs)})
}
