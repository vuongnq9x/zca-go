package zca

import (
	"context"
	"encoding/json"
	"net/http"
)

type GetGroupChatHistoryResponse struct {
	LastActionID      string     `json:"lastActionId"`
	LastActionIDOther string     `json:"lastActionIdOther"`
	More              int64      `json:"more"`
	GroupMsgs         []*Message `json:"groupMsgs"`
}

// GetGroupChatHistory gets recent messages of a group. count 0 means 50.
func (a *API) GetGroupChatHistory(ctx context.Context, groupID string, count int) (*GetGroupChatHistoryResponse, error) {
	if count == 0 {
		count = 50
	}
	raw, err := call[json.RawMessage](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group")+"/api/group/history", nil, true), map[string]any{"grid": groupID, "count": count})
	if err != nil {
		return nil, err
	}
	// data may arrive as a JSON-encoded string.
	if len(raw) > 0 && raw[0] == '"' {
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return nil, err
		}
		raw = json.RawMessage(s)
	}
	data, err := decodeData[struct {
		LastActionID      string        `json:"lastActionId"`
		LastActionIDOther string        `json:"lastActionIdOther"`
		More              int64         `json:"more"`
		GroupMsgs         []MessageData `json:"groupMsgs"`
	}](raw)
	if err != nil {
		return nil, err
	}
	out := &GetGroupChatHistoryResponse{LastActionID: data.LastActionID, LastActionIDOther: data.LastActionIDOther, More: data.More}
	for _, m := range data.GroupMsgs {
		out.GroupMsgs = append(out.GroupMsgs, NewGroupMessage(a.UID, m))
	}
	return out, nil
}
