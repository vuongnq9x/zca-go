package zca

import (
	"context"
	"net/http"
)

type SendTypingEventResponse struct {
	Status int64 `json:"status"`
}

// SendTypingEvent sends a typing event. destType (User threads only) 0 means DestTypeUser.
func (a *API) SendTypingEvent(ctx context.Context, threadID string, threadType ThreadType, destType DestType) (*SendTypingEventResponse, error) {
	if threadID == "" {
		return nil, newError("Missing threadId")
	}
	if destType == 0 {
		destType = DestTypeUser
	}
	params := map[string]any{"imei": a.IMEI}
	u := a.svc("chat") + "/api/message/typing"
	if threadType == ThreadTypeUser {
		params["toid"], params["destType"] = threadID, destType
	} else {
		params["grid"] = threadID
		u = a.svc("group") + "/api/group/typing"
	}
	return call[*SendTypingEventResponse](ctx, a.Session, http.MethodPost, a.MakeURL(u, nil, true), params)
}
