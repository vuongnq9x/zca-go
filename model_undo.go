package zca

type TUndoContent struct {
	GlobalMsgID int64 `json:"globalMsgId"`
	CliMsgID    int64 `json:"cliMsgId"`
	DeleteMsg   int64 `json:"deleteMsg"`
	SrcID       int64 `json:"srcId"`
	DestID      int64 `json:"destId"`
}

type TUndo struct {
	ActionID  string       `json:"actionId"`
	MsgID     string       `json:"msgId"`
	CliMsgID  string       `json:"cliMsgId"`
	MsgType   string       `json:"msgType"`
	UIDFrom   string       `json:"uidFrom"`
	IDTo      string       `json:"idTo"`
	DName     string       `json:"dName"`
	TS        string       `json:"ts"`
	Status    int64        `json:"status"`
	Content   TUndoContent `json:"content"`
	Notify    string       `json:"notify"`
	TTL       int64        `json:"ttl"`
	UserID    string       `json:"userId"`
	Uin       string       `json:"uin"`
	Cmd       int64        `json:"cmd"`
	St        int64        `json:"st"`
	At        int64        `json:"at"`
	RealMsgID string       `json:"realMsgId"`
}

// Undo is TS class Undo; isGroup is expressed as Type == ThreadTypeGroup.
type Undo struct {
	Type     ThreadType `json:"type"`
	Data     TUndo      `json:"data"`
	ThreadID string     `json:"threadId"`
	IsSelf   bool       `json:"isSelf"`
}

func NewUndo(uid string, data TUndo, isGroup bool) *Undo {
	u := &Undo{Type: ThreadTypeUser, ThreadID: data.UIDFrom, IsSelf: data.UIDFrom == "0"}
	if isGroup {
		u.Type = ThreadTypeGroup
	}
	if isGroup || data.UIDFrom == "0" {
		u.ThreadID = data.IDTo
	}
	if data.IDTo == "0" {
		data.IDTo = uid
	}
	if data.UIDFrom == "0" {
		data.UIDFrom = uid
	}
	u.Data = data
	return u
}
