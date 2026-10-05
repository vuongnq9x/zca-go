package zca

import (
	"context"
	"net/http"
)

type SendSeenEventResponse struct {
	Status int64 `json:"status"`
}

// SendSeenEventMessageParams has the same shape as SendDeliveredEventMessageParams.
type SendSeenEventMessageParams = SendDeliveredEventMessageParams

// SendSeenEvent sends a message seen event (1-50 messages of the same thread).
func (a *API) SendSeenEvent(ctx context.Context, messages []SendSeenEventMessageParams, threadType ThreadType) (*SendSeenEventResponse, error) {
	if len(messages) == 0 || len(messages) > 50 { // MAX_MESSAGES_PER_SEND
		return nil, newError("messages must contain between 1 and 50 messages.")
	}
	isGroup := threadType == ThreadTypeGroup
	threadOf := func(m SendSeenEventMessageParams) string {
		if isGroup {
			return m.IDTo
		}
		return m.UIDFrom
	}
	threadID := threadOf(messages[0])
	data := make([]sendDeliveredEventRequestMsg, len(messages))
	for i, m := range messages {
		if threadOf(m) != threadID {
			return nil, newError("All messages must belong to the same thread.")
		}
		data[i] = sendDeliveredEventMsg(a.UID, m)
	}
	msgInfos := map[string]any{"data": data}
	params := map[string]any{}
	u := a.svc("chat") + "/api/message/seenv2"
	if isGroup {
		msgInfos["grid"] = threadID
		params["imei"] = a.IMEI
		u = a.svc("group") + "/api/group/seenv2"
	} else {
		msgInfos["senderId"] = threadID
	}
	params["msgInfos"] = mustJSON(msgInfos)
	return call[*SendSeenEventResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(u, map[string]any{"nretry": 0}, true), params)
}
