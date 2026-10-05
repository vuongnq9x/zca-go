package zca

import (
	"context"
	"encoding/json"
	"net/http"
)

type EditReminderOptions struct {
	Title   string
	TopicID string
	Emoji   string
	// Unix ms; default now when 0.
	StartTime int64
	Repeat    ReminderRepeatMode
}

type EditReminderUser = ReminderUser

// EditReminderGroup is ReminderGroup without responseMem.
type EditReminderGroup = ReminderGroup

// EditReminderResponse holds User for ThreadTypeUser, Group for ThreadTypeGroup.
type EditReminderResponse = CreateReminderResponse

// EditReminder edits an existing reminder.
func (a *API) EditReminder(ctx context.Context, options EditReminderOptions, threadID string, threadType ThreadType) (*EditReminderResponse, error) {
	if options.StartTime == 0 {
		options.StartTime = nowMs()
	}
	var params map[string]any
	var base string
	if threadType == ThreadTypeUser {
		base = a.svc("group_board") + "/api/board/oneone/update"
		params = map[string]any{
			"objectData": mustJSON(map[string]any{
				"toUid":      threadID,
				"type":       0,
				"color":      -16777216,
				"emoji":      options.Emoji,
				"startTime":  options.StartTime,
				"duration":   -1,
				"params":     map[string]string{"title": options.Title},
				"needPin":    false,
				"reminderId": options.TopicID,
				"repeat":     options.Repeat,
			}),
		}
	} else {
		base = a.svc("group_board") + "/api/board/topic/updatev2"
		params = map[string]any{
			"grid":      threadID,
			"type":      0,
			"color":     -16777216,
			"emoji":     options.Emoji,
			"startTime": options.StartTime,
			"duration":  -1,
			"params":    mustJSON(map[string]string{"title": options.Title}),
			"topicId":   options.TopicID,
			"repeat":    options.Repeat,
			"imei":      a.IMEI,
			"pinAct":    2,
		}
	}
	raw, err := call[json.RawMessage](ctx, a.Session, http.MethodPost, a.MakeURL(base, nil, true), params)
	if err != nil {
		return nil, err
	}
	return createReminderDecode(raw, threadType)
}
