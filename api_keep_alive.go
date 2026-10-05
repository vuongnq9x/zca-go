package zca

import (
	"context"
	"encoding/json"
	"net/http"
)

type KeepAliveResponse struct {
	ConfigVesion int64 `json:"config_vesion"`
}

// KeepAlive pings the chat service. Its response is not encrypted.
func (a *API) KeepAlive(ctx context.Context) (*KeepAliveResponse, error) {
	enc, err := a.EncodeAES(mustJSON(map[string]any{"imei": a.IMEI}))
	if err != nil {
		return nil, newError("Failed to encrypt params")
	}
	u := a.MakeURL(a.svc("chat")+"/keepalive", map[string]any{"params": enc}, true)
	resp, err := a.Request(ctx, http.MethodGet, u, nil, nil)
	if err != nil {
		return nil, err
	}
	raw, err := a.Resolve(resp, false)
	if err != nil {
		return nil, err
	}
	return decodeData[*KeepAliveResponse](json.RawMessage(raw))
}
