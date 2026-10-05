package zca

import (
	"context"
	"encoding/json"
	"net/http"
)

type GetLabelsResponse struct {
	LabelData      []LabelData `json:"labelData"`
	Version        int64       `json:"version"`
	LastUpdateTime int64       `json:"lastUpdateTime"`
}

// GetLabels gets all conversation labels.
func (a *API) GetLabels(ctx context.Context) (*GetLabelsResponse, error) {
	raw, err := call[*struct {
		LabelData      string `json:"labelData"`
		Version        int64  `json:"version"`
		LastUpdateTime int64  `json:"lastUpdateTime"`
	}](ctx, a.Session, http.MethodGet, a.MakeURL(a.svc("label")+"/api/convlabel/get", nil, true), map[string]any{"imei": a.IMEI})
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, newError("Failed to parse response data")
	}
	out := &GetLabelsResponse{Version: raw.Version, LastUpdateTime: raw.LastUpdateTime}
	if err := json.Unmarshal([]byte(raw.LabelData), &out.LabelData); err != nil {
		return nil, err
	}
	return out, nil
}
