package zca

import (
	"context"
	"net/http"
	"time"
)

// MuteDuration is a mute duration in seconds or one of the predefined values.
type MuteDuration int64

const (
	MuteDurationOneHour   MuteDuration = 3600
	MuteDurationFourHours MuteDuration = 14400
	MuteDurationForever   MuteDuration = -1
	// MuteDurationUntil8AM stands for TS "until8AM" (mute until the next local 8:00).
	MuteDurationUntil8AM MuteDuration = -2
)

type MuteAction int

const (
	MuteActionMute   MuteAction = 1
	MuteActionUnmute MuteAction = 3
)

// SetMuteParams: zero Duration means MuteDurationForever, zero Action means MuteActionMute.
type SetMuteParams struct {
	Duration MuteDuration
	Action   MuteAction
}

// SetMute mutes or unmutes a thread.
func (a *API) SetMute(ctx context.Context, p SetMuteParams, threadID string, threadType ThreadType) (string, error) {
	if p.Duration == 0 {
		p.Duration = MuteDurationForever
	}
	if p.Action == 0 {
		p.Action = MuteActionMute
	}
	now := time.Now()
	duration := int64(p.Duration)
	switch {
	case p.Action == MuteActionUnmute, p.Duration == MuteDurationForever:
		duration = -1
	case p.Duration == MuteDurationUntil8AM:
		next := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, time.Local)
		if now.Hour() >= 8 {
			next = next.AddDate(0, 0, 1)
		}
		duration = int64(next.Sub(now) / time.Second)
	}
	muteType := 2
	if threadType == ThreadTypeUser {
		muteType = 1
	}
	params := map[string]any{
		"toid": threadID, "duration": duration, "action": p.Action,
		"startTime": now.Unix(), "muteType": muteType, "imei": a.IMEI,
	}
	return call[string](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("profile")+"/api/social/profile/setmute", nil, true), params)
}
