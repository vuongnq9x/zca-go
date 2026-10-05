package zca

import (
	"encoding/binary"
	"encoding/json"
	"testing"
)

func wsFrame(cmd, subCmd int, payload any) []byte {
	body, _ := json.Marshal(payload)
	f := []byte{1, 0, 0, byte(subCmd)}
	binary.LittleEndian.PutUint16(f[1:3], uint16(cmd))
	return append(f, body...)
}

func TestListenerHandle(t *testing.T) {
	s := &Session{UID: "999"}
	l := newListener(s, nil)
	var got []*Message
	var undos []*Undo
	l.OnMessage = func(m *Message) { got = append(got, m) }
	l.OnUndo = func(u *Undo) { undos = append(undos, u) }

	data := `{"data":{"msgs":[
		{"msgId":"1","uidFrom":"123","idTo":"0","content":"hi","ts":"1"},
		{"msgId":"2","uidFrom":"0","idTo":"123","content":"mine","ts":"2"},
		{"msgId":"3","uidFrom":"123","idTo":"0","content":{"deleteMsg":1,"globalMsgId":1},"ts":"3"}]}}`
	if err := l.handle(nil, wsFrame(501, 0, map[string]any{"data": data, "encrypt": 0})); err != nil {
		t.Fatal(err)
	}
	// self message dropped (SelfListen=false), undo routed separately
	if len(got) != 1 || got[0].ThreadID != "123" || got[0].Data.IDTo != "999" || got[0].IsSelf {
		t.Fatalf("messages = %+v", got)
	}
	if len(undos) != 1 {
		t.Fatalf("undos = %d", len(undos))
	}

	s.Options.SelfListen = true
	got = nil
	l.handle(nil, wsFrame(501, 0, map[string]any{"data": data, "encrypt": 0}))
	if len(got) != 2 || !got[1].IsSelf || got[1].ThreadID != "123" {
		t.Fatalf("selfListen messages = %+v", got)
	}

	var typing *Typing
	l.OnTyping = func(ty *Typing) { typing = ty }
	td := `{"data":{"actions":[{"act_type":"typing","act":"typing","data":"\"uid\":\"5\",\"ts\":\"1\",\"isPC\":0"}]}}`
	if err := l.handle(nil, wsFrame(602, 0, map[string]any{"data": td, "encrypt": 0})); err != nil || typing == nil || typing.Data.UID != "5" {
		t.Fatalf("typing = %+v, %v", typing, err)
	}

	// cmd 551 (E2EE 1-1 push) routes like 501.
	got = nil
	if err := l.handle(nil, wsFrame(551, 0, map[string]any{"data": data, "encrypt": 0})); err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].Type != ThreadTypeUser || got[0].Data.Content != "hi" {
		t.Fatalf("551 messages = %+v", got)
	}
	if err := l.handle(nil, wsFrame(621, 0, map[string]any{})); err != nil {
		t.Fatalf("unhandled cmd: %v", err)
	}

	var cleared [][]*ClearUnread
	l.OnUnreadCleared = func(c []*ClearUnread) { cleared = append(cleared, c) }
	cu := `{"data":{"clearUnreads":[{"idTo":"123","isGroup":0,"lastMsgId":"77","type":0},{"idTo":"-1","isGroup":0,"lastMsgId":"0","type":2}]}}`
	l.handle(nil, wsFrame(504, 0, map[string]any{"data": cu, "encrypt": 0}))
	l.handle(nil, wsFrame(524, 0, map[string]any{"data": cu, "encrypt": 0}))
	l.handle(nil, wsFrame(504, 0, map[string]any{"data": `{"data":{"clearUnreads":[]}}`, "encrypt": 0}))
	if len(cleared) != 2 || len(cleared[0]) != 1 || cleared[0][0].Type != ThreadTypeUser || cleared[0][0].ThreadID != "123" ||
		cleared[0][0].Data.LastMsgID != "77" || cleared[1][0].Type != ThreadTypeGroup {
		t.Fatalf("cleared = %+v", cleared)
	}
}

func TestListenerRetry(t *testing.T) {
	s := &Session{}
	s.Settings.Features.Socket.Retries = map[string]struct {
		Max   int `json:"max"`
		Times any `json:"times"`
	}{"internal": {Max: 2, Times: []any{1000.0, 2000.0}}}
	l := newListener(s, []string{"a", "b"})
	for _, want := range []int64{1000, 2000} {
		if d, ok := l.canRetry(CloseReasonAbnormalClosure); !ok || d != want {
			t.Fatalf("canRetry = %d %v, want %d", d, ok, want)
		}
	}
	if _, ok := l.canRetry(CloseReasonAbnormalClosure); ok {
		t.Fatal("retry past max")
	}
	if _, ok := l.canRetry(CloseReasonManualClosure); ok {
		t.Fatal("manual closure retried")
	}
	if !l.shouldRotate(CloseReasonAbnormalClosure) {
		t.Fatal("1006 should rotate")
	}
	l.resetRetryCount()
	if _, ok := l.canRetry(CloseReasonAbnormalClosure); !ok {
		t.Fatal("reset did not restore retries")
	}
}
