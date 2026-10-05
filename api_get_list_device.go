package zca

import (
	"context"
	"net/http"
)

type GetListDeviceResponse struct {
	Devices struct {
		MasterID     string `json:"masterId"`
		EncIdentity  string `json:"encIdentity"`
		LastUpdateTs int64  `json:"lastUpdateTs"`
		EncSignature string `json:"encSignature"`
		Companions   []any  `json:"companions"`
	} `json:"devices"`
}

// GetListDevice gets linked devices.
func (a *API) GetListDevice(ctx context.Context) (*GetListDeviceResponse, error) {
	return call[*GetListDeviceResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("aext")+"/api/devices/linked", nil, true), map[string]any{"imei": a.IMEI})
}
