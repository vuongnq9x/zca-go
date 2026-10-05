package zca

import (
	"context"
	"net/http"
)

type MuteEntriesInfo struct {
	ID          string `json:"id"`
	Duration    int64  `json:"duration"`
	StartTime   int64  `json:"startTime"`
	SystemTime  int64  `json:"systemTime"`
	CurrentTime int64  `json:"currentTime"`
	MuteMode    int64  `json:"muteMode"`
}

type GetMuteResponse struct {
	ChatEntries      []MuteEntriesInfo `json:"chatEntries"`
	GroupChatEntries []MuteEntriesInfo `json:"groupChatEntries"`
}

// GetMute gets muted conversations.
func (a *API) GetMute(ctx context.Context) (*GetMuteResponse, error) {
	return call[*GetMuteResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("profile")+"/api/social/profile/getmute", nil, true), map[string]any{"imei": a.IMEI})
}
