package zca

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type UpdateLabelsPayload struct {
	LabelData []LabelData
	Version   int64
}

type UpdateLabelsResponse struct {
	LabelData      []LabelData `json:"labelData"`
	Version        int64       `json:"version"`
	LastUpdateTime int64       `json:"lastUpdateTime"`
}

// UpdateLabels replaces the conversation labels.
func (a *API) UpdateLabels(ctx context.Context, payload UpdateLabelsPayload) (*UpdateLabelsResponse, error) {
	params := map[string]any{"labelData": mustJSON(payload.LabelData), "version": payload.Version, "imei": a.IMEI}
	raw, err := call[*struct {
		LabelData      string `json:"labelData"`
		Version        int64  `json:"version"`
		LastUpdateTime int64  `json:"lastUpdateTime"`
	}](ctx, a.Session, http.MethodPost, a.MakeURL(a.svc("label")+"/api/convlabel/update", nil, true), params)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, newError("Failed to parse response data")
	}
	res := &UpdateLabelsResponse{Version: raw.Version, LastUpdateTime: raw.LastUpdateTime}
	if err := json.Unmarshal([]byte(raw.LabelData), &res.LabelData); err != nil {
		return nil, fmt.Errorf("zca: decode labelData: %w", err)
	}
	return res, nil
}
