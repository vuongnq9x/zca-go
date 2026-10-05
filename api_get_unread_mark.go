package zca

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type UnreadMark struct {
	ID       int64 `json:"id"`
	CliMsgID int64 `json:"cliMsgId"`
	FromUID  int64 `json:"fromUid"`
	Ts       int64 `json:"ts"`
}

type GetUnreadMarkResponse struct {
	Data struct {
		ConvsGroup []UnreadMark `json:"convsGroup"`
		ConvsUser  []UnreadMark `json:"convsUser"`
	} `json:"data"`
	Status int64 `json:"status"`
}

// GetUnreadMark gets conversations marked as unread.
func (a *API) GetUnreadMark(ctx context.Context) (*GetUnreadMarkResponse, error) {
	raw, err := call[json.RawMessage](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("conversation")+"/api/conv/getUnreadMark", nil, true), map[string]any{})
	if err != nil {
		return nil, err
	}
	out := &GetUnreadMarkResponse{}
	return out, getUnreadMarkDecode(raw, out)
}

// getUnreadMarkDecode decodes {data, status} where data may be a JSON-encoded string.
func getUnreadMarkDecode(raw json.RawMessage, out any) error {
	var env struct {
		Data   json.RawMessage `json:"data"`
		Status int64           `json:"status"`
	}
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}
	if err := json.Unmarshal(raw, &env); err != nil {
		return fmt.Errorf("zca: decode response: %w", err)
	}
	if len(env.Data) > 0 && env.Data[0] == '"' {
		var s string
		if err := json.Unmarshal(env.Data, &s); err != nil {
			return fmt.Errorf("zca: decode response: %w", err)
		}
		env.Data = json.RawMessage(s)
	}
	b, err := json.Marshal(env)
	if err == nil {
		err = json.Unmarshal(b, out)
	}
	if err != nil {
		return fmt.Errorf("zca: decode response: %w", err)
	}
	return nil
}
