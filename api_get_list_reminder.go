package zca

import (
	"context"
	"encoding/json"
	"net/http"
)

type ListReminderOptions struct {
	Page  int // default 1 when 0
	Count int // default 20 when 0
}

// ReminderListItem is TS ReminderListUser & ReminderListGroup (fields of both, flattened).
type ReminderListItem struct {
	CreatorUID  string              `json:"creatorUid"`
	ToUID       string              `json:"toUid"`
	ReminderID  string              `json:"reminderId"`
	EndTime     int64               `json:"endTime"`
	EditorID    string              `json:"editorId"`
	Emoji       string              `json:"emoji"`
	Color       int64               `json:"color"`
	GroupID     string              `json:"groupId"`
	CreatorID   string              `json:"creatorId"`
	EditTime    int64               `json:"editTime"`
	EventType   int64               `json:"eventType"`
	ResponseMem ReminderResponseMem `json:"responseMem"`
	Params      ReminderParams      `json:"params"`
	Type        int64               `json:"type"`
	Duration    int64               `json:"duration"`
	RepeatInfo  *ReminderRepeatInfo `json:"repeatInfo"`
	RepeatData  []any               `json:"repeatData"`
	CreateTime  int64               `json:"createTime"`
	Repeat      ReminderRepeatMode  `json:"repeat"`
	StartTime   int64               `json:"startTime"`
	ID          string              `json:"id"`
}

// GetListReminder gets reminders of a user or group thread.
func (a *API) GetListReminder(ctx context.Context, options ListReminderOptions, threadID string, threadType ThreadType) ([]ReminderListItem, error) {
	if options.Page == 0 {
		options.Page = 1
	}
	if options.Count == 0 {
		options.Count = 20
	}
	obj := map[string]any{
		"board_type": 1,
		"page":       options.Page,
		"count":      options.Count,
		"last_id":    0,
		"last_type":  0,
	}
	params := map[string]any{}
	path := "/api/board/oneone/list"
	if threadType == ThreadTypeGroup {
		obj["group_id"] = threadID
		params["imei"] = a.IMEI
		path = "/api/board/listReminder"
	} else {
		obj["uid"] = threadID
	}
	params["objectData"] = mustJSON(obj)
	s, err := call[string](ctx, a.Session, http.MethodGet, a.MakeURL(a.svc("group_board")+path, nil, true), params)
	if err != nil {
		return nil, err
	}
	var out []ReminderListItem
	if err := json.Unmarshal([]byte(s), &out); err != nil {
		return nil, err
	}
	return out, nil
}
