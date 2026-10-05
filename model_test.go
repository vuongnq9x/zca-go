package zca

import (
	"encoding/json"
	"testing"
)

func TestNewUserMessage(t *testing.T) {
	m := NewUserMessage("me", MessageData{UIDFrom: "0", IDTo: "peer"})
	if m.Type != ThreadTypeUser || m.ThreadID != "peer" || !m.IsSelf || m.Data.UIDFrom != "me" {
		t.Fatalf("self user message: %+v", m)
	}
	m = NewUserMessage("me", MessageData{UIDFrom: "peer", IDTo: "0"})
	if m.ThreadID != "peer" || m.IsSelf || m.Data.IDTo != "me" {
		t.Fatalf("incoming user message: %+v", m)
	}
}

func TestNewGroupMessage(t *testing.T) {
	m := NewGroupMessage("me", MessageData{UIDFrom: "0", IDTo: "g1"})
	if m.Type != ThreadTypeGroup || m.ThreadID != "g1" || !m.IsSelf || m.Data.UIDFrom != "me" {
		t.Fatalf("self group message: %+v", m)
	}
	m = NewGroupMessage("me", MessageData{UIDFrom: "u2", IDTo: "g1"})
	if m.ThreadID != "g1" || m.IsSelf {
		t.Fatalf("incoming group message: %+v", m)
	}
}

func TestQuoteOwnerIDNumber(t *testing.T) {
	var d MessageData
	if err := json.Unmarshal([]byte(`{"quote":{"ownerId":1234567890123456789}}`), &d); err != nil {
		t.Fatal(err)
	}
	if d.Quote.OwnerID != "1234567890123456789" {
		t.Fatalf("ownerId: %q", d.Quote.OwnerID)
	}
}

func TestInitializeEvents(t *testing.T) {
	fe := InitializeFriendEvent("me", json.RawMessage(`"123"`), getFriendEventType("block"))
	if fe.Type != FriendEventTypeBLOCK || fe.Data != "123" || fe.ThreadID != "123" || !fe.IsSelf {
		t.Fatalf("friend block: %+v", fe)
	}
	fe = InitializeFriendEvent("me", json.RawMessage(`"{\"topic\":{\"params\":\"{\\\"title\\\":\\\"x\\\"}\"},\"actorId\":\"me\",\"conversationId\":\"c\"}"`), FriendEventTypePIN_CREATE)
	if d := fe.Data.(*TFriendEventPinCreate); d.Topic.Params.Title != "x" || fe.ThreadID != "c" || !fe.IsSelf {
		t.Fatalf("friend pin_create: %+v", fe)
	}
	ge := InitializeGroupEvent("me", json.RawMessage(`{"group_id":"g1","creatorId":"me"}`), getGroupEventType("remind_topic"), "remind_topic")
	if ge.ThreadID != "g1" || !ge.IsSelf {
		t.Fatalf("group remind_topic: %+v", ge)
	}
	ge = InitializeGroupEvent("me", json.RawMessage(`{"groupId":"g2","sourceId":"x","updateMembers":[{"id":"me"}]}`), getGroupEventType("join"), "join")
	if ge.Type != GroupEventTypeJOIN || ge.ThreadID != "g2" || !ge.IsSelf {
		t.Fatalf("group join: %+v", ge)
	}
}
