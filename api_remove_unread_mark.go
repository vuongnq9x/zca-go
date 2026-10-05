package zca

import (
	"context"
	"encoding/json"
	"net/http"
)

type RemoveUnreadMarkResponse struct {
	Data struct {
		UpdateID int64 `json:"updateId"`
	} `json:"data"`
	Status int64 `json:"status"`
}

// RemoveUnreadMark removes the unread mark from a conversation.
func (a *API) RemoveUnreadMark(ctx context.Context, threadID string, threadType ThreadType) (*RemoveUnreadMarkResponse, error) {
	self, other := "convsUser", "convsGroup"
	if threadType == ThreadTypeGroup {
		self, other = other, self
	}
	param := map[string]any{
		self:           []string{threadID},
		other:          []string{},
		self + "Data":  []map[string]any{{"id": threadID, "ts": nowMs()}},
		other + "Data": []any{},
	}
	raw, err := call[json.RawMessage](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("conversation")+"/api/conv/removeUnreadMark", nil, true), map[string]any{"param": mustJSON(param)})
	if err != nil {
		return nil, err
	}
	out := &RemoveUnreadMarkResponse{}
	return out, getUnreadMarkDecode(raw, out)
}
