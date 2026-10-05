package zca

import "slices"

// TSeenMessage merges TS TUserSeenMessage (idTo, realMsgId) and TGroupSeenMessage (groupId, seenUids).
type TSeenMessage struct {
	MsgID     string   `json:"msgId"`
	IDTo      string   `json:"idTo,omitempty"`
	RealMsgID string   `json:"realMsgId,omitempty"`
	GroupID   string   `json:"groupId,omitempty"`
	SeenUIDs  []string `json:"seenUids,omitempty"`
}

// SeenMessage is TS UserSeenMessage | GroupSeenMessage.
type SeenMessage struct {
	Type     ThreadType   `json:"type"`
	Data     TSeenMessage `json:"data"`
	ThreadID string       `json:"threadId"`
	IsSelf   bool         `json:"isSelf"`
}

func NewUserSeenMessage(data TSeenMessage) *SeenMessage {
	return &SeenMessage{Type: ThreadTypeUser, Data: data, ThreadID: data.IDTo}
}

func NewGroupSeenMessage(uid string, data TSeenMessage) *SeenMessage {
	return &SeenMessage{Type: ThreadTypeGroup, Data: data, ThreadID: data.GroupID, IsSelf: slices.Contains(data.SeenUIDs, uid)}
}
