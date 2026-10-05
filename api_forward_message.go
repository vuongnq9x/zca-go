package zca

import (
	"context"
	"net/http"
	"strconv"
)

type ForwardMessageReference struct {
	ID         string `json:"id"`
	TS         int64  `json:"ts"`
	LogSrcType int64  `json:"logSrcType"`
	FwLvl      int64  `json:"fwLvl"`
}

type ForwardMessagePayload struct {
	Message   string
	TTL       int64 // milliseconds
	Reference *ForwardMessageReference
}

type ForwardMessageSuccess struct {
	ClientID StringOrNumber `json:"clientId"`
	MsgID    StringOrNumber `json:"msgId"`
}

type ForwardMessageFail struct {
	ClientID  StringOrNumber `json:"clientId"`
	ErrorCode StringOrNumber `json:"error_code"`
}

type ForwardMessageResponse struct {
	Success []ForwardMessageSuccess `json:"success"`
	Fail    []ForwardMessageFail    `json:"fail"`
}

// ForwardMessage forwards a text message to multiple threads of the same type.
func (a *API) ForwardMessage(ctx context.Context, payload ForwardMessagePayload, threadIDs []string, threadType ThreadType) (*ForwardMessageResponse, error) {
	if payload.Message == "" {
		return nil, newError("Missing message content")
	}
	if len(threadIDs) == 0 {
		return nil, newError("Missing thread IDs")
	}
	clientID := strconv.FormatInt(nowMs(), 10)

	msgInfo := map[string]any{"message": payload.Message}
	var decorLog any // JSON null without a reference
	if r := payload.Reference; r != nil {
		msgInfo["reference"] = mustJSON(map[string]any{"type": 3, "data": mustJSON(r)})
		msg := map[string]any{"st": 1, "ts": r.TS, "id": r.ID}
		decorLog = map[string]any{"fw": map[string]any{"pmsg": msg, "rmsg": msg, "fwLvl": r.FwLvl}}
	}

	params := map[string]any{
		"ttl":      payload.TTL,
		"msgType":  "1",
		"totalIds": len(threadIDs),
		"msgInfo":  mustJSON(msgInfo),
		"decorLog": mustJSON(decorLog),
	}
	targets := make([]map[string]any, len(threadIDs))
	path := "/api/message/mforward"
	for i, id := range threadIDs {
		targets[i] = map[string]any{"clientId": clientID, "ttl": payload.TTL}
		if threadType == ThreadTypeUser {
			targets[i]["toUid"] = id
		} else {
			targets[i]["grid"] = id
		}
	}
	if threadType == ThreadTypeUser {
		params["toIds"], params["imei"] = targets, a.IMEI
	} else {
		params["grids"] = targets
		path = "/api/group/mforward"
	}
	res, err := call[*ForwardMessageResponse](ctx, a.Session, http.MethodPost, a.MakeURL(a.svc("file")+path, nil, true), params)
	if err == nil && res == nil {
		res = &ForwardMessageResponse{}
	}
	return res, err
}
