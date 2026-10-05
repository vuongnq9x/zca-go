package zca

// TTyping is TS TTyping plus TGroupTyping.gid.
type TTyping struct {
	UID  string `json:"uid"`
	TS   string `json:"ts"`
	IsPC int64  `json:"isPC"` // 0 | 1
	GID  string `json:"gid,omitempty"`
}

// Typing is TS UserTyping | GroupTyping. IsSelf is always false.
type Typing struct {
	Type     ThreadType `json:"type"`
	Data     TTyping    `json:"data"`
	ThreadID string     `json:"threadId"`
	IsSelf   bool       `json:"isSelf"`
}

func NewUserTyping(data TTyping) *Typing {
	return &Typing{Type: ThreadTypeUser, Data: data, ThreadID: data.UID}
}

func NewGroupTyping(data TTyping) *Typing {
	return &Typing{Type: ThreadTypeGroup, Data: data, ThreadID: data.GID}
}
