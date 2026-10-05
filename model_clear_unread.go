package zca

type TClearUnread struct {
	IDTo      string `json:"idTo"`
	IsGroup   int64  `json:"isGroup"`
	LastMsgID string `json:"lastMsgId"`
	Type      int64  `json:"type"`
	Ts        int64  `json:"ts,omitempty"`
}

// ClearUnread is TS UserClearUnread | GroupClearUnread.
type ClearUnread struct {
	Type     ThreadType   `json:"type"`
	Data     TClearUnread `json:"data"`
	ThreadID string       `json:"threadId"`
}

func NewUserClearUnread(d TClearUnread) *ClearUnread {
	return &ClearUnread{Type: ThreadTypeUser, Data: d, ThreadID: d.IDTo}
}

func NewGroupClearUnread(d TClearUnread) *ClearUnread {
	return &ClearUnread{Type: ThreadTypeGroup, Data: d, ThreadID: d.IDTo}
}
