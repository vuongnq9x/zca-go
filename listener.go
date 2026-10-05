package zca

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type CloseReason int

const (
	CloseReasonManualClosure       CloseReason = 1000
	CloseReasonAbnormalClosure     CloseReason = 1006
	CloseReasonDuplicateConnection CloseReason = 3000
	CloseReasonKickConnection      CloseReason = 3003
)

// Listener receives realtime events over Zalo's websocket (zca-js Listener).
// Set the On* handlers before Start; they run on the listener goroutine.
type Listener struct {
	OnConnected         func()
	OnDisconnected      func(code CloseReason, reason string)
	OnClosed            func(code CloseReason, reason string)
	OnError             func(error)
	OnMessage           func(*Message)
	OnTyping            func(*Typing)
	OnOldMessages       func([]*Message, ThreadType)
	OnSeenMessages      func([]*SeenMessage)
	OnDeliveredMessages func([]*DeliveredMessage)
	OnReaction          func(*Reaction)
	OnOldReactions      func(reactions []*Reaction, isGroup bool)
	OnUploadAttachment  func(UploadEventData)
	OnUndo              func(*Undo)
	OnFriendEvent       func(*FriendEvent)
	OnGroupEvent        func(*GroupEvent)
	OnCipherKey         func(key string)

	s           *Session
	urls        []string
	rotateCount int
	retries     map[int]*retryState

	mu        sync.Mutex
	conn      *websocket.Conn
	cipherKey string
	reqID     int
	stopPing  chan struct{}
}

type retryState struct {
	count, max int
	times      []int64
}

func newListener(s *Session, urls []string) *Listener {
	l := &Listener{s: s, urls: urls, retries: map[int]*retryState{}}
	for code, r := range s.Settings.Features.Socket.Retries {
		var c int
		if _, err := fmt.Sscan(code, &c); err != nil {
			continue
		}
		st := &retryState{max: r.Max}
		switch t := r.Times.(type) {
		case float64:
			st.times = []int64{int64(t)}
		case []any:
			for _, v := range t {
				if f, ok := v.(float64); ok {
					st.times = append(st.times, int64(f))
				}
			}
		}
		l.retries[c] = st
	}
	return l
}

func (l *Listener) wsURL() string {
	return l.s.MakeURL(l.urls[l.rotateCount], map[string]any{"t": nowMs()}, true)
}

// Start connects the websocket. With retryOnClose, it reconnects on the close codes
// configured by Zalo (rotating endpoints when asked).
func (l *Listener) Start(retryOnClose bool) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.conn != nil {
		return newError("Already started")
	}
	if len(l.urls) == 0 {
		return newError("No websocket endpoint")
	}
	u := l.wsURL()
	h := http.Header{}
	h.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	h.Set("Accept-Language", "en-US,en;q=0.9")
	h.Set("Cache-Control", "no-cache")
	h.Set("Origin", "https://chat.zalo.me")
	h.Set("Pragma", "no-cache")
	h.Set("User-Agent", l.s.UserAgent)
	h.Set("Cookie", l.s.jar.header("https://chat.zalo.me"))

	d := websocket.Dialer{EnableCompression: true, HandshakeTimeout: 30 * time.Second, Proxy: http.ProxyFromEnvironment}
	if t, ok := l.s.client.Transport.(*http.Transport); ok && t.Proxy != nil {
		d.Proxy = t.Proxy
	}
	conn, resp, err := d.Dial(u, h)
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
	if err != nil {
		return err
	}
	l.conn = conn
	l.stopPing = make(chan struct{})
	go l.readLoop(conn, retryOnClose)
	return nil
}

// Stop closes the connection with ManualClosure.
func (l *Listener) Stop() {
	l.mu.Lock()
	conn := l.conn
	l.mu.Unlock()
	if conn != nil {
		conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(int(CloseReasonManualClosure), ""), time.Now().Add(time.Second))
		conn.Close()
	}
}

func (l *Listener) reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.conn = nil
	l.cipherKey = ""
	if l.stopPing != nil {
		close(l.stopPing)
		l.stopPing = nil
	}
}

func (l *Listener) readLoop(conn *websocket.Conn, retryOnClose bool) {
	if l.OnConnected != nil {
		l.OnConnected()
	}
	for {
		mt, data, err := conn.ReadMessage()
		if err != nil {
			code, reason := CloseReasonAbnormalClosure, ""
			var ce *websocket.CloseError
			if errors.As(err, &ce) {
				code, reason = CloseReason(ce.Code), ce.Text
			}
			conn.Close()
			l.onClose(code, reason, retryOnClose)
			return
		}
		if mt != websocket.BinaryMessage {
			continue
		}
		if err := l.handle(conn, data); err != nil {
			l.emitError(err)
		}
	}
}

func (l *Listener) onClose(code CloseReason, reason string, retryOnClose bool) {
	l.reset()
	if l.OnDisconnected != nil {
		l.OnDisconnected(code, reason)
	}
	if retryOnClose {
		if delay, ok := l.canRetry(code); ok {
			if l.shouldRotate(code) {
				l.rotateCount++
				l.s.log().Debug("Rotating websocket endpoint", "index", l.rotateCount)
			}
			time.AfterFunc(time.Duration(delay)*time.Millisecond, func() {
				if err := l.Start(true); err != nil {
					l.emitError(err)
					if l.OnClosed != nil {
						l.OnClosed(CloseReasonAbnormalClosure, err.Error())
					}
				}
			})
			return
		}
	}
	if l.OnClosed != nil {
		l.OnClosed(code, reason)
	}
}

func (l *Listener) canRetry(code CloseReason) (int64, bool) {
	sock := l.s.Settings.Features.Socket
	r := l.retries[int(code)]
	if !slices.Contains(sock.CloseAndRetryCodes, int(code)) || r == nil || r.count >= r.max || len(r.times) == 0 {
		return 0, false
	}
	r.count++
	delay := r.times[len(r.times)-1]
	if r.count-1 < len(r.times) {
		delay = r.times[r.count-1]
	}
	l.s.log().Debug("Retry websocket", "code", code, "delayMs", delay, "count", r.count, "max", r.max)
	return delay, true
}

func (l *Listener) shouldRotate(code CloseReason) bool {
	return slices.Contains(l.s.Settings.Features.Socket.RotateErrorCodes, int(code)) && l.rotateCount < len(l.urls)-1
}

func (l *Listener) emitError(err error) {
	if l.OnError != nil {
		l.OnError(err)
	}
}

// SendWs sends a frame: [version, cmd (uint16 LE), subCmd, json data].
func (l *Listener) SendWs(version, cmd, subCmd int, data map[string]any, requireID bool) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.conn == nil {
		return nil
	}
	if requireID {
		data["req_id"] = fmt.Sprintf("req_%d", l.reqID)
		l.reqID++
	}
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	frame := make([]byte, 4, 4+len(body))
	frame[0] = byte(version)
	binary.LittleEndian.PutUint16(frame[1:3], uint16(cmd))
	frame[3] = byte(subCmd)
	return l.conn.WriteMessage(websocket.BinaryMessage, append(frame, body...))
}

// RequestOldMessages asks for message history; results arrive on OnOldMessages.
func (l *Listener) RequestOldMessages(t ThreadType, lastMsgID *string) error {
	cmd := 510
	if t == ThreadTypeGroup {
		cmd = 511
	}
	return l.SendWs(1, cmd, 1, map[string]any{"first": true, "lastId": lastMsgID, "preIds": []any{}}, true)
}

// RequestOldReactions asks for reaction history; results arrive on OnOldReactions.
func (l *Listener) RequestOldReactions(t ThreadType, lastMsgID *string) error {
	cmd := 610
	if t == ThreadTypeGroup {
		cmd = 611
	}
	return l.SendWs(1, cmd, 1, map[string]any{"first": true, "lastId": lastMsgID, "preIds": []any{}}, true)
}

type wsEnvelope struct {
	Key     *string `json:"key"`
	Data    string  `json:"data"`
	Encrypt int     `json:"encrypt"`
}

// decode decrypts an event and unmarshals its "data" field into out.
func (l *Listener) decode(env *wsEnvelope, out any) error {
	l.mu.Lock()
	key := l.cipherKey
	l.mu.Unlock()
	plain, err := decodeEventData(env.Data, env.Encrypt, key)
	if err != nil {
		return err
	}
	return json.Unmarshal(plain, &struct {
		Data any `json:"data"`
	}{out})
}

func (l *Listener) selfOK(isSelf bool) bool { return !isSelf || l.s.Options.SelfListen }

func (l *Listener) handle(conn *websocket.Conn, data []byte) error {
	if len(data) < 4 {
		return newError("Invalid header")
	}
	version, cmd, subCmd := int(data[0]), int(binary.LittleEndian.Uint16(data[1:3])), int(data[3])
	payload := data[4:]
	if len(payload) == 0 {
		return nil
	}
	var env wsEnvelope
	if err := json.Unmarshal(payload, &env); err != nil {
		return err
	}
	uid := l.s.UID

	switch {
	case version == 1 && cmd == 1 && subCmd == 1 && env.Key != nil:
		l.mu.Lock()
		l.cipherKey = *env.Key
		stop := l.stopPing
		l.mu.Unlock()
		if l.OnCipherKey != nil {
			l.OnCipherKey(*env.Key)
		}
		if iv := l.s.Settings.Features.Socket.PingInterval; iv > 0 && stop != nil {
			go func() {
				t := time.NewTicker(time.Duration(iv) * time.Millisecond)
				defer t.Stop()
				for {
					select {
					case <-stop:
						return
					case <-t.C:
						l.SendWs(1, 2, 1, map[string]any{"eventId": nowMs()}, false)
					}
				}
			}()
		}

	case version == 1 && (cmd == 501 || cmd == 521) && subCmd == 0:
		isGroup := cmd == 521
		var d struct {
			Msgs      []json.RawMessage `json:"msgs"`
			GroupMsgs []json.RawMessage `json:"groupMsgs"`
		}
		if err := l.decode(&env, &d); err != nil {
			return err
		}
		msgs := d.Msgs
		if isGroup {
			msgs = d.GroupMsgs
		}
		for _, raw := range msgs {
			var probe struct {
				Content json.RawMessage `json:"content"`
			}
			json.Unmarshal(raw, &probe)
			var obj map[string]json.RawMessage
			if json.Unmarshal(probe.Content, &obj) == nil && obj["deleteMsg"] != nil {
				var u TUndo
				if err := json.Unmarshal(raw, &u); err != nil {
					l.emitError(err)
					continue
				}
				undo := NewUndo(uid, u, isGroup)
				if l.selfOK(undo.IsSelf) && l.OnUndo != nil {
					l.OnUndo(undo)
				}
				continue
			}
			var md MessageData
			if err := json.Unmarshal(raw, &md); err != nil {
				l.emitError(err)
				continue
			}
			m := NewUserMessage(uid, md)
			if isGroup {
				m = NewGroupMessage(uid, md)
			}
			if l.selfOK(m.IsSelf) && l.OnMessage != nil {
				l.OnMessage(m)
			}
		}

	case version == 1 && cmd == 601 && subCmd == 0:
		var d struct {
			Controls []struct {
				Content struct {
					ActType string          `json:"act_type"`
					Act     string          `json:"act"`
					Data    json.RawMessage `json:"data"`
					FileID  json.RawMessage `json:"fileId"`
				} `json:"content"`
			} `json:"controls"`
		}
		if err := l.decode(&env, &d); err != nil {
			return err
		}
		for _, c := range d.Controls {
			ct := c.Content
			switch ct.ActType {
			case "file_done":
				var fd struct {
					URL string `json:"url"`
				}
				json.Unmarshal(unwrapJSONString(ct.Data), &fd)
				ev := UploadEventData{FileURL: fd.URL, FileID: jsonScalarString(ct.FileID)}
				l.s.notifyUpload(ev)
				if l.OnUploadAttachment != nil {
					l.OnUploadAttachment(ev)
				}
			case "group":
				// Zalo sends both join and join_reject when an admin approves; ignore join_reject like zca-js.
				if ct.Act == "join_reject" {
					continue
				}
				ge := InitializeGroupEvent(uid, ct.Data, getGroupEventType(ct.Act), ct.Act)
				if l.selfOK(ge.IsSelf) && l.OnGroupEvent != nil {
					l.OnGroupEvent(ge)
				}
			case "fr":
				// Zalo sends both req and req_v2; ignore req like zca-js.
				if ct.Act == "req" {
					continue
				}
				fe := InitializeFriendEvent(uid, ct.Data, getFriendEventType(ct.Act))
				if l.selfOK(fe.IsSelf) && l.OnFriendEvent != nil {
					l.OnFriendEvent(fe)
				}
			}
		}

	case cmd == 612:
		var d struct {
			Reacts      []TReaction `json:"reacts"`
			ReactGroups []TReaction `json:"reactGroups"`
		}
		if err := l.decode(&env, &d); err != nil {
			return err
		}
		for i, list := range [][]TReaction{d.Reacts, d.ReactGroups} {
			for _, r := range list {
				ro := NewReaction(uid, r, i == 1)
				if l.selfOK(ro.IsSelf) && l.OnReaction != nil {
					l.OnReaction(ro)
				}
			}
		}

	case cmd == 610 || cmd == 611:
		isGroup := cmd == 611
		var d struct {
			Reacts      []TReaction `json:"reacts"`
			ReactGroups []TReaction `json:"reactGroups"`
		}
		if err := l.decode(&env, &d); err != nil {
			return err
		}
		list := d.Reacts
		if isGroup {
			list = d.ReactGroups
		}
		out := make([]*Reaction, 0, len(list))
		for _, r := range list {
			out = append(out, NewReaction(uid, r, isGroup))
		}
		if l.OnOldReactions != nil {
			l.OnOldReactions(out, isGroup)
		}

	case (cmd == 510 || cmd == 511) && subCmd == 1:
		isGroup := cmd == 511
		var d struct {
			Msgs      []MessageData `json:"msgs"`
			GroupMsgs []MessageData `json:"groupMsgs"`
		}
		if err := l.decode(&env, &d); err != nil {
			return err
		}
		var out []*Message
		if isGroup {
			for _, m := range d.GroupMsgs {
				out = append(out, NewGroupMessage(uid, m))
			}
		} else {
			for _, m := range d.Msgs {
				out = append(out, NewUserMessage(uid, m))
			}
		}
		if l.OnOldMessages != nil {
			t := ThreadTypeUser
			if isGroup {
				t = ThreadTypeGroup
			}
			l.OnOldMessages(out, t)
		}

	case cmd == 602 && subCmd == 0:
		var d struct {
			Actions []struct {
				ActType string `json:"act_type"`
				Act     string `json:"act"`
				Data    string `json:"data"`
			} `json:"actions"`
		}
		if err := l.decode(&env, &d); err != nil {
			return err
		}
		for _, a := range d.Actions {
			if a.ActType != "typing" {
				continue
			}
			var td TTyping
			if err := json.Unmarshal([]byte("{"+a.Data+"}"), &td); err != nil {
				l.emitError(err)
				continue
			}
			var t *Typing
			switch a.Act {
			case "typing":
				t = NewUserTyping(td)
			case "gtyping":
				t = NewGroupTyping(td)
			}
			if t != nil && l.OnTyping != nil {
				l.OnTyping(t)
			}
		}

	case (cmd == 502 || cmd == 522) && subCmd == 0:
		isGroup := cmd == 522
		var d struct {
			Delivereds []TDeliveredMessage `json:"delivereds"`
			Seens      []TSeenMessage      `json:"seens"`
			GroupSeens []TSeenMessage      `json:"groupSeens"`
		}
		if err := l.decode(&env, &d); err != nil {
			return err
		}
		var delivered []*DeliveredMessage
		for _, m := range d.Delivereds {
			dm := NewUserDeliveredMessage(m)
			if isGroup {
				dm = NewGroupDeliveredMessage(uid, m)
			}
			if !isGroup || l.selfOK(dm.IsSelf) {
				delivered = append(delivered, dm)
			}
		}
		if len(d.Delivereds) > 0 && l.OnDeliveredMessages != nil {
			l.OnDeliveredMessages(delivered)
		}
		seenList := d.Seens
		if isGroup {
			seenList = d.GroupSeens
		}
		var seen []*SeenMessage
		for _, m := range seenList {
			sm := NewUserSeenMessage(m)
			if isGroup {
				sm = NewGroupSeenMessage(uid, m)
			}
			if !isGroup || l.selfOK(sm.IsSelf) {
				seen = append(seen, sm)
			}
		}
		if len(seenList) > 0 && l.OnSeenMessages != nil {
			l.OnSeenMessages(seen)
		}

	case version == 1 && cmd == 3000 && subCmd == 0:
		l.s.log().Error("Another connection is opened, closing this one")
		conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(int(CloseReasonDuplicateConnection), ""), time.Now().Add(time.Second))
	}
	return nil
}
