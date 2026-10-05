package zca

import (
	"context"
	"encoding/json"
	"net/http"
)

type GetFriendOnlinesStatus struct {
	UserID string `json:"userId"`
	Status string `json:"status"`
}

type GetFriendOnlinesResponse struct {
	Predefine   []string                 `json:"predefine"`
	OwnerStatus string                   `json:"ownerStatus"`
	Onlines     []GetFriendOnlinesStatus `json:"onlines"`
}

// GetFriendOnlines gets online friends.
func (a *API) GetFriendOnlines(ctx context.Context) (*GetFriendOnlinesResponse, error) {
	data, err := call[*GetFriendOnlinesResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("profile")+"/api/social/friend/onlines", nil, true), map[string]any{"imei": a.IMEI})
	if err != nil || data == nil {
		return data, err
	}
	return data, getFriendOnlinesUnwrap(data)
}

// getFriendOnlinesUnwrap replaces each status (a JSON string like {"status":"..."}) with its inner status.
func getFriendOnlinesUnwrap(data *GetFriendOnlinesResponse) error {
	for i := range data.Onlines {
		var parsed any
		if err := json.Unmarshal([]byte(data.Onlines[i].Status), &parsed); err != nil {
			return err
		}
		if m, ok := parsed.(map[string]any); ok {
			if s, ok := m["status"].(string); ok {
				data.Onlines[i].Status = s
			}
		}
	}
	return nil
}
