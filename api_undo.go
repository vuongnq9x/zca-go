package zca

import (
	"context"
	"net/http"
)

type UndoPayload struct {
	MsgID    string
	CliMsgID string
}

type UndoResponse struct {
	Status int64 `json:"status"`
}

// Undo recalls a message.
func (a *API) Undo(ctx context.Context, payload UndoPayload, threadID string, threadType ThreadType) (*UndoResponse, error) {
	params := map[string]any{"msgId": payload.MsgID, "clientId": nowMs(), "cliMsgIdUndo": payload.CliMsgID}
	u := a.svc("chat") + "/api/message/undo"
	if threadType == ThreadTypeGroup {
		params["grid"], params["visibility"], params["imei"] = threadID, 0, a.IMEI
		u = a.svc("group") + "/api/group/undomsg"
	} else {
		params["toid"] = threadID
	}
	return call[*UndoResponse](ctx, a.Session, http.MethodPost, a.MakeURL(u, nil, true), params)
}
