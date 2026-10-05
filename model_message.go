package zca

import (
	"bytes"
	"encoding/json"
)

// StringOrNumber is a string that also accepts a bare JSON number (digits kept exactly).
// Used where TS does String(x) on a field Zalo may send as a number.
type StringOrNumber string

func (s *StringOrNumber) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) > 0 && b[0] == '"' {
		var v string
		if err := json.Unmarshal(b, &v); err != nil {
			return err
		}
		*s = StringOrNumber(v)
		return nil
	}
	if string(b) == "null" {
		return nil
	}
	*s = StringOrNumber(b)
	return nil
}

type TAttachmentContent struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Href        string `json:"href"`
	Thumb       string `json:"thumb"`
	Childnumber int64  `json:"childnumber"`
	Action      string `json:"action"`
	Params      string `json:"params"`
	Type        string `json:"type"`
}

type TMessagePropertyExt struct {
	Color   int64  `json:"color"`
	Size    int64  `json:"size"`
	Type    int64  `json:"type"`
	SubType int64  `json:"subType"`
	Ext     string `json:"ext"`
}

type TMessageParamsExt struct {
	CountUnread  int64 `json:"countUnread"`
	ContainType  int64 `json:"containType"`
	PlatformType int64 `json:"platformType"`
}

// MessageData is TS TMessage plus TGroupMessage.mentions.
type MessageData struct {
	ActionID string `json:"actionId"`
	MsgID    string `json:"msgId"`
	CliMsgID string `json:"cliMsgId"`
	MsgType  string `json:"msgType"`
	UIDFrom  string `json:"uidFrom"`
	IDTo     string `json:"idTo"`
	DName    string `json:"dName"`
	TS       string `json:"ts"`
	Status   int64  `json:"status"`
	// Content is a string or an object (map[string]any; see TAttachmentContent).
	Content           any                  `json:"content"`
	Notify            string               `json:"notify"`
	TTL               int64                `json:"ttl"`
	UserID            string               `json:"userId"`
	Uin               string               `json:"uin"`
	TopOut            string               `json:"topOut"`
	TopOutTimeOut     string               `json:"topOutTimeOut"`
	TopOutImprTimeOut string               `json:"topOutImprTimeOut"`
	PropertyExt       *TMessagePropertyExt `json:"propertyExt,omitempty"`
	ParamsExt         TMessageParamsExt    `json:"paramsExt"`
	Cmd               int64                `json:"cmd"`
	St                int64                `json:"st"`
	At                int64                `json:"at"`
	RealMsgID         string               `json:"realMsgId"`
	Quote             *TQuote              `json:"quote,omitempty"`
	Mentions          []TMention           `json:"mentions,omitempty"`
}

type TQuote struct {
	OwnerID     StringOrNumber `json:"ownerId"`
	CliMsgID    int64          `json:"cliMsgId"`
	GlobalMsgID int64          `json:"globalMsgId"`
	CliMsgType  int64          `json:"cliMsgType"`
	TS          int64          `json:"ts"`
	Msg         string         `json:"msg"`
	Attach      string         `json:"attach"`
	FromD       string         `json:"fromD"`
	TTL         int64          `json:"ttl"`
}

type TMention struct {
	UID  string `json:"uid"`
	Pos  int64  `json:"pos"`
	Len  int64  `json:"len"`
	Type int64  `json:"type"` // 0 | 1
}

// Message is TS UserMessage | GroupMessage, discriminated by Type.
type Message struct {
	Type     ThreadType  `json:"type"`
	Data     MessageData `json:"data"`
	ThreadID string      `json:"threadId"`
	// IsSelf is true if the message is sent by the logged in account.
	IsSelf bool `json:"isSelf"`
}

func NewUserMessage(uid string, data MessageData) *Message {
	m := &Message{Type: ThreadTypeUser, ThreadID: data.UIDFrom, IsSelf: data.UIDFrom == "0"}
	if data.UIDFrom == "0" {
		m.ThreadID = data.IDTo
	}
	if data.IDTo == "0" {
		data.IDTo = uid
	}
	if data.UIDFrom == "0" {
		data.UIDFrom = uid
	}
	m.Data = data
	return m
}

func NewGroupMessage(uid string, data MessageData) *Message {
	m := &Message{Type: ThreadTypeGroup, ThreadID: data.IDTo, IsSelf: data.UIDFrom == "0"}
	if data.UIDFrom == "0" {
		data.UIDFrom = uid
	}
	m.Data = data
	return m
}
