package zca

import (
	"context"
	"net/http"
)

type DeleteMessageResponse struct {
	Status int64 `json:"status"`
}

type DeleteMessageDestination struct {
	Data struct {
		CliMsgID string
		MsgID    string
		UIDFrom  string
	}
	ThreadID string
	Type     ThreadType
}

// DeleteMessage deletes a message. Deleting for everyone only works for others'
// messages in groups; use Undo for your own messages.
func (a *API) DeleteMessage(ctx context.Context, dest DeleteMessageDestination, onlyMe bool) (*DeleteMessageResponse, error) {
	isGroup := dest.Type == ThreadTypeGroup
	if a.UID == dest.Data.UIDFrom && !onlyMe {
		return nil, newError("To delete your message for everyone, use undo api instead")
	}
	if !isGroup && !onlyMe {
		return nil, newError("Can't delete message for everyone in a private chat")
	}
	only := 0
	if onlyMe {
		only = 1
	}
	params := map[string]any{
		"cliMsgId": nowMs(),
		"msgs": []map[string]any{{
			"cliMsgId":    dest.Data.CliMsgID,
			"globalMsgId": dest.Data.MsgID,
			"ownerId":     dest.Data.UIDFrom,
			"destId":      dest.ThreadID,
		}},
		"onlyMe": only,
	}
	u := a.svc("chat") + "/api/message/delete"
	if isGroup {
		params["grid"] = dest.ThreadID
		u = a.svc("group") + "/api/group/deletemsg"
	} else {
		params["toid"] = dest.ThreadID
		params["imei"] = a.IMEI
	}
	return call[*DeleteMessageResponse](ctx, a.Session, http.MethodPost, a.MakeURL(u, nil, true), params)
}
