package zca

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type AddUnreadMarkResponse struct {
	Data struct {
		UpdateID int64 `json:"updateId"`
	} `json:"data"`
	Status int64 `json:"status"`
}

// AddUnreadMark marks a conversation as unread.
func (a *API) AddUnreadMark(ctx context.Context, threadID string, threadType ThreadType) (*AddUnreadMarkResponse, error) {
	ts := nowMs()
	conv := []map[string]any{{"id": threadID, "cliMsgId": strconv.FormatInt(ts, 10), "fromUid": "0", "ts": ts}}
	inner := map[string]any{"convsUser": conv, "convsGroup": []any{}, "imei": a.IMEI}
	if threadType == ThreadTypeGroup {
		inner["convsUser"], inner["convsGroup"] = []any{}, conv
	}
	raw, err := call[struct {
		Data   json.RawMessage `json:"data"`
		Status int64           `json:"status"`
	}](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("conversation")+"/api/conv/addUnreadMark", nil, true),
		map[string]any{"param": mustJSON(inner)})
	if err != nil {
		return nil, err
	}
	out := &AddUnreadMarkResponse{Status: raw.Status}
	out.Data, err = decodeData[struct {
		UpdateID int64 `json:"updateId"`
	}](addReactionUnquote(raw.Data))
	if err != nil {
		return nil, err
	}
	return out, nil
}
