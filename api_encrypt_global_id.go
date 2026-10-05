package zca

import (
	"context"
	"net/http"
)

// EncryptGlobalIDResponse maps user id to global id.
type EncryptGlobalIDResponse = map[string]string

// EncryptGlobalID encrypts user ids into global ids.
func (a *API) EncryptGlobalID(ctx context.Context, userIDs []string) (EncryptGlobalIDResponse, error) {
	if len(userIDs) == 0 {
		return nil, newError("Missing global id")
	}
	// noiseUids is a JSON-encoded string, not an array.
	return call[EncryptGlobalIDResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("profile")+"/api/gid/encrypt", nil, true), map[string]any{"noiseUids": mustJSON(userIDs)})
}
