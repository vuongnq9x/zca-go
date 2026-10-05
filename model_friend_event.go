package zca

import (
	"bytes"
	"encoding/json"
)

type FriendEventType int

const (
	FriendEventTypeADD FriendEventType = iota
	FriendEventTypeREMOVE

	FriendEventTypeREQUEST
	FriendEventTypeUNDO_REQUEST
	FriendEventTypeREJECT_REQUEST

	FriendEventTypeSEEN_FRIEND_REQUEST

	FriendEventTypeBLOCK
	FriendEventTypeUNBLOCK
	FriendEventTypeBLOCK_CALL
	FriendEventTypeUNBLOCK_CALL

	FriendEventTypePIN_UNPIN
	FriendEventTypePIN_CREATE

	FriendEventTypeUNKNOWN
)

type TFriendEventRejectUndo struct {
	ToUID   string `json:"toUid"`
	FromUID string `json:"fromUid"`
}

type TFriendEventRequest struct {
	ToUID   string `json:"toUid"`
	FromUID string `json:"fromUid"`
	Src     int64  `json:"src"`
	Message string `json:"message"`
}

type TFriendEventPinCreateTopicParams struct {
	SenderUID   string `json:"senderUid"`
	SenderName  string `json:"senderName"`
	ClientMsgID string `json:"client_msg_id"`
	GlobalMsgID string `json:"global_msg_id"`
	MsgType     int64  `json:"msg_type"`
	Title       string `json:"title"`
}

// UnmarshalJSON also accepts params sent as a JSON-encoded string (listen.ts parses it).
func (p *TFriendEventPinCreateTopicParams) UnmarshalJSON(b []byte) error {
	type plain TFriendEventPinCreateTopicParams
	b = unwrapJSONString(b)
	return json.Unmarshal(b, (*plain)(p))
}

type TFriendEventPinTopic struct {
	TopicID   string `json:"topicId"`
	TopicType int64  `json:"topicType"`
}

type TFriendEventPinCreateTopic struct {
	Type       int64                            `json:"type"`
	Color      int64                            `json:"color"`
	Emoji      string                           `json:"emoji"`
	StartTime  int64                            `json:"startTime"`
	Duration   int64                            `json:"duration"`
	Params     TFriendEventPinCreateTopicParams `json:"params"`
	ID         string                           `json:"id"`
	CreatorID  string                           `json:"creatorId"`
	CreateTime int64                            `json:"createTime"`
	EditorID   string                           `json:"editorId"`
	EditTime   int64                            `json:"editTime"`
	Repeat     int64                            `json:"repeat"`
	Action     int64                            `json:"action"`
}

type TFriendEventPinCreate struct {
	OldTopic       *TFriendEventPinTopic      `json:"oldTopic,omitempty"`
	Topic          TFriendEventPinCreateTopic `json:"topic"`
	ActorID        string                     `json:"actorId"`
	OldVersion     int64                      `json:"oldVersion"`
	Version        int64                      `json:"version"`
	ConversationID string                     `json:"conversationId"`
}

type TFriendEventPinUnpin struct {
	Topic          TFriendEventPinTopic `json:"topic"`
	ActorID        string               `json:"actorId"`
	OldVersion     int64                `json:"oldVersion"`
	Version        int64                `json:"version"`
	ConversationID string               `json:"conversationId"`
}

// FriendEvent is the TS FriendEvent union. Data holds, by Type:
//   - ADD, REMOVE, BLOCK, UNBLOCK, BLOCK_CALL, UNBLOCK_CALL: string (uid)
//   - REJECT_REQUEST, UNDO_REQUEST: *TFriendEventRejectUndo
//   - REQUEST: *TFriendEventRequest
//   - SEEN_FRIEND_REQUEST: []string
//   - PIN_CREATE: *TFriendEventPinCreate
//   - PIN_UNPIN: *TFriendEventPinUnpin
//   - UNKNOWN: string (JSON text)
type FriendEvent struct {
	Type     FriendEventType `json:"type"`
	Data     any             `json:"data"`
	ThreadID string          `json:"threadId"`
	IsSelf   bool            `json:"isSelf"`
}

// InitializeFriendEvent ports initializeFriendEvent. data is the raw control content.data:
// a JSON string holding JSON is unwrapped first, as listen.ts does.
// ponytail: decode errors are ignored (best-effort, like the TS casts).
func InitializeFriendEvent(uid string, data json.RawMessage, typ FriendEventType) *FriendEvent {
	inner := unwrapJSONString(data)
	if isJSONNumber(inner) {
		// listen.ts passes the original content.data when the parsed value is a number.
		inner = bytes.TrimSpace(data)
	}
	switch typ {
	case FriendEventTypeADD, FriendEventTypeREMOVE, FriendEventTypeBLOCK,
		FriendEventTypeUNBLOCK, FriendEventTypeBLOCK_CALL, FriendEventTypeUNBLOCK_CALL:
		s := jsonScalarString(inner)
		return &FriendEvent{Type: typ, Data: s, ThreadID: s, IsSelf: typ != FriendEventTypeADD && typ != FriendEventTypeREMOVE}
	case FriendEventTypeREJECT_REQUEST, FriendEventTypeUNDO_REQUEST:
		var d TFriendEventRejectUndo
		_ = json.Unmarshal(inner, &d)
		return &FriendEvent{Type: typ, Data: &d, ThreadID: d.ToUID, IsSelf: d.FromUID == uid}
	case FriendEventTypeREQUEST:
		var d TFriendEventRequest
		_ = json.Unmarshal(inner, &d)
		return &FriendEvent{Type: typ, Data: &d, ThreadID: d.ToUID, IsSelf: d.FromUID == uid}
	case FriendEventTypeSEEN_FRIEND_REQUEST:
		var d []string
		_ = json.Unmarshal(inner, &d)
		return &FriendEvent{Type: typ, Data: d, ThreadID: uid, IsSelf: true}
	case FriendEventTypePIN_CREATE:
		var d TFriendEventPinCreate
		_ = json.Unmarshal(inner, &d)
		return &FriendEvent{Type: typ, Data: &d, ThreadID: d.ConversationID, IsSelf: d.ActorID == uid}
	case FriendEventTypePIN_UNPIN:
		var d TFriendEventPinUnpin
		_ = json.Unmarshal(inner, &d)
		return &FriendEvent{Type: typ, Data: &d, ThreadID: d.ConversationID, IsSelf: d.ActorID == uid}
	default:
		return &FriendEvent{Type: FriendEventTypeUNKNOWN, Data: string(inner)}
	}
}

func getFriendEventType(act string) FriendEventType {
	switch act {
	case "add":
		return FriendEventTypeADD
	case "remove":
		return FriendEventTypeREMOVE
	case "block":
		return FriendEventTypeBLOCK
	case "unblock":
		return FriendEventTypeUNBLOCK
	case "block_call":
		return FriendEventTypeBLOCK_CALL
	case "unblock_call":
		return FriendEventTypeUNBLOCK_CALL
	case "req_v2":
		return FriendEventTypeREQUEST
	case "reject":
		return FriendEventTypeREJECT_REQUEST
	case "undo_req":
		return FriendEventTypeUNDO_REQUEST
	case "seen_fr_req":
		return FriendEventTypeSEEN_FRIEND_REQUEST
	case "pin_unpin":
		return FriendEventTypePIN_UNPIN
	case "pin_create":
		return FriendEventTypePIN_CREATE
	}
	return FriendEventTypeUNKNOWN
}

// unwrapJSONString returns the inner JSON when b is a JSON string whose content is valid JSON
// (mirrors `typeof x == "string" ? JSON.parse(x) : x`); otherwise b unchanged.
func unwrapJSONString(b []byte) []byte {
	b = bytes.TrimSpace(b)
	if len(b) == 0 || b[0] != '"' {
		return b
	}
	var s string
	if json.Unmarshal(b, &s) == nil && json.Valid([]byte(s)) {
		return bytes.TrimSpace([]byte(s))
	}
	return b
}

func isJSONNumber(b []byte) bool {
	return len(b) > 0 && (b[0] == '-' || (b[0] >= '0' && b[0] <= '9'))
}

// jsonScalarString returns a JSON string's value, or the raw text for other scalars (numbers).
func jsonScalarString(b []byte) string {
	var s string
	if json.Unmarshal(b, &s) == nil {
		return s
	}
	return string(b)
}
