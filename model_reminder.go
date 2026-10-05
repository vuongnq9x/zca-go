package zca

type ReminderRepeatMode int

const (
	ReminderRepeatModeNone    ReminderRepeatMode = 0
	ReminderRepeatModeDaily   ReminderRepeatMode = 1
	ReminderRepeatModeWeekly  ReminderRepeatMode = 2
	ReminderRepeatModeMonthly ReminderRepeatMode = 3
)

type ReminderParams struct {
	Title    string `json:"title"`
	SetTitle bool   `json:"setTitle"`
}

type ReminderUser struct {
	CreatorUID string             `json:"creatorUid"`
	ToUID      string             `json:"toUid"`
	Emoji      string             `json:"emoji"`
	Color      int64              `json:"color"`
	ReminderID string             `json:"reminderId"`
	CreateTime int64              `json:"createTime"`
	Repeat     ReminderRepeatMode `json:"repeat"`
	StartTime  int64              `json:"startTime"`
	EditTime   int64              `json:"editTime"`
	EndTime    int64              `json:"endTime"`
	Params     ReminderParams     `json:"params"`
	Type       int64              `json:"type"`
}

type ReminderResponseMem struct {
	RejectMember int64 `json:"rejectMember"`
	MyResp       int64 `json:"myResp"`
	AcceptMember int64 `json:"acceptMember"`
}

type ReminderRepeatInfo struct {
	ListTS []any `json:"list_ts"`
}

type ReminderGroup struct {
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
