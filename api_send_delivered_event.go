package zca

import (
	"context"
	"net/http"
	"strings"
)

// SendDeliveredEventResponse is "" or {status: number}.
type SendDeliveredEventResponse = any

type SendDeliveredEventMessageParams struct {
	MsgID    string
	CliMsgID string
	UIDFrom  string
	IDTo     string
	MsgType  string
	St       int64
	At       int64
	Cmd      int64
	Ts       string // TS string | number
}

type sendDeliveredEventRequestMsg struct {
	Cmi string `json:"cmi"`
	Gmi string `json:"gmi"`
	Si  string `json:"si"`
	Di  string `json:"di"`
	Mt  string `json:"mt"`
	St  int64  `json:"st"`
	At  int64  `json:"at"`
	Cmd int64  `json:"cmd"`
	Ts  int64  `json:"ts"`
}

// sendDeliveredEventMsg mirrors the TS mapping: st/at/cmd are always defined in Go so they
// become 0; ts is 0 when parseInt(ts) is a number, else -1.
func sendDeliveredEventMsg(uid string, m SendDeliveredEventMessageParams) sendDeliveredEventRequestMsg {
	di := m.IDTo
	if di == uid {
		di = "0"
	}
	ts := int64(-1)
	s := strings.TrimLeft(strings.TrimSpace(m.Ts), "+-")
	if s != "" && s[0] >= '0' && s[0] <= '9' {
		ts = 0
	}
	return sendDeliveredEventRequestMsg{Cmi: m.CliMsgID, Gmi: m.MsgID, Si: m.UIDFrom, Di: di, Mt: m.MsgType, Ts: ts}
}

// SendDeliveredEvent sends a message delivered event (1-50 messages).
func (a *API) SendDeliveredEvent(ctx context.Context, isSeen bool, messages []SendDeliveredEventMessageParams, threadType ThreadType) (SendDeliveredEventResponse, error) {
	if len(messages) == 0 || len(messages) > 50 { // MAX_MESSAGES_PER_SEND
		return nil, newError("messages must contain between 1 and 50 messages.")
	}
	idTo := messages[0].IDTo
	data := make([]sendDeliveredEventRequestMsg, len(messages))
	for i, m := range messages {
		if threadType == ThreadTypeGroup && m.IDTo != idTo {
			return nil, newError("All messages must have the same idTo for Group thread")
		}
		data[i] = sendDeliveredEventMsg(a.UID, m)
	}
	seen := 0
	if isSeen {
		seen = 1
	}
	msgInfos := map[string]any{"seen": seen, "data": data}
	params := map[string]any{}
	u := a.svc("chat") + "/api/message/deliveredv2"
	if threadType == ThreadTypeGroup {
		msgInfos["grid"] = idTo
		params["imei"] = a.IMEI
		u = a.svc("group") + "/api/group/deliveredv2"
	}
	params["msgInfos"] = mustJSON(msgInfos)
	return call[SendDeliveredEventResponse](ctx, a.Session, http.MethodPost, a.MakeURL(u, nil, true), params)
}
