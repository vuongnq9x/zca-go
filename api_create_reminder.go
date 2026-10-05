package zca

import (
	"context"
	"encoding/json"
	"net/http"
)

type CreateReminderOptions struct {
	Title string
	// Default "⏰" when empty.
	Emoji string
	// Unix ms; default now when 0.
	StartTime int64
	// Default ReminderRepeatModeNone.
	Repeat ReminderRepeatMode
}

type CreateReminderUser = ReminderUser

// CreateReminderGroup is ReminderGroup without responseMem.
type CreateReminderGroup = ReminderGroup

// CreateReminderResponse holds User for ThreadTypeUser, Group for ThreadTypeGroup.
type CreateReminderResponse struct {
	User  *CreateReminderUser
	Group *CreateReminderGroup
}

// CreateReminder creates a reminder in a user or group thread.
func (a *API) CreateReminder(ctx context.Context, options CreateReminderOptions, threadID string, threadType ThreadType) (*CreateReminderResponse, error) {
	if options.Emoji == "" {
		options.Emoji = "⏰"
	}
	if options.StartTime == 0 {
		options.StartTime = nowMs()
	}
	var params map[string]any
	var base string
	if threadType == ThreadTypeUser {
		base = a.svc("group_board") + "/api/board/oneone/create"
		params = map[string]any{
			"objectData": mustJSON(map[string]any{
				"toUid":      threadID,
				"type":       0,
				"color":      -16245706,
				"emoji":      options.Emoji,
				"startTime":  options.StartTime,
				"duration":   -1,
				"params":     map[string]string{"title": options.Title},
				"needPin":    false,
				"repeat":     options.Repeat,
				"creatorUid": a.UID,
				"src":        1,
			}),
			"imei": a.IMEI,
		}
	} else {
		base = a.svc("group_board") + "/api/board/topic/createv2"
		params = map[string]any{
			"grid":      threadID,
			"type":      0,
			"color":     -16245706,
			"emoji":     options.Emoji,
			"startTime": options.StartTime,
			"duration":  -1,
			"params":    mustJSON(map[string]string{"title": options.Title}),
			"repeat":    options.Repeat,
			"src":       1,
			"imei":      a.IMEI,
			"pinAct":    0,
		}
	}
	raw, err := call[json.RawMessage](ctx, a.Session, http.MethodPost, a.MakeURL(base, nil, true), params)
	if err != nil {
		return nil, err
	}
	return createReminderDecode(raw, threadType)
}

// createReminderDecode decodes a reminder response by thread type (shared with EditReminder).
func createReminderDecode(raw json.RawMessage, threadType ThreadType) (*CreateReminderResponse, error) {
	var out CreateReminderResponse
	var err error
	if threadType == ThreadTypeUser {
		out.User, err = decodeData[*ReminderUser](raw)
	} else {
		out.Group, err = decodeData[*ReminderGroup](raw)
	}
	if err != nil {
		return nil, err
	}
	return &out, nil
}
