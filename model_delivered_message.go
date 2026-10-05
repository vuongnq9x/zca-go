package zca

import "slices"

// TDeliveredMessage is TS TDeliveredMessage plus TGroupDeliveredMessage.groupId.
type TDeliveredMessage struct {
	MsgID         string   `json:"msgId"`
	Seen          int64    `json:"seen"`
	DeliveredUIDs []string `json:"deliveredUids"`
	SeenUIDs      []string `json:"seenUids"`
	RealMsgID     string   `json:"realMsgId"`
	MSTs          int64    `json:"mSTs"`
	GroupID       string   `json:"groupId,omitempty"`
}

// DeliveredMessage is TS UserDeliveredMessage | GroupDeliveredMessage.
type DeliveredMessage struct {
	Type     ThreadType        `json:"type"`
	Data     TDeliveredMessage `json:"data"`
	ThreadID string            `json:"threadId"`
	IsSelf   bool              `json:"isSelf"`
}

func NewUserDeliveredMessage(data TDeliveredMessage) *DeliveredMessage {
	m := &DeliveredMessage{Type: ThreadTypeUser, Data: data}
	if len(data.DeliveredUIDs) > 0 {
		m.ThreadID = data.DeliveredUIDs[0]
	}
	return m
}

func NewGroupDeliveredMessage(uid string, data TDeliveredMessage) *DeliveredMessage {
	return &DeliveredMessage{Type: ThreadTypeGroup, Data: data, ThreadID: data.GroupID, IsSelf: slices.Contains(data.DeliveredUIDs, uid)}
}
