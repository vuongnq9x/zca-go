package zca

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type GetGroupChatHistoryResponse struct {
	Error                int64      `json:"error"`
	LastMsgID            int64      `json:"lastMsgId"`
	MinMsgID             int64      `json:"minMsgId"`
	MaxMsgID             int64      `json:"maxMsgId"`
	MsgJumpID            int64      `json:"msgJumpId"`
	HasMore              bool       `json:"hasMore"`
	IsOld                bool       `json:"isOld"`
	IsFiltered           bool       `json:"isFiltered"`
	RootMsgID            int64      `json:"rootMsgId"`
	IsRootDel            bool       `json:"isRootDel"`
	TsJoinGroup          int64      `json:"tsJoinGroup"`
	IsFilteredByPhase    bool       `json:"isFilteredByPhase"`
	IsFilteredByTimeJoin bool       `json:"isFilteredByTimeJoin"`
	GroupMsgs            []*Message `json:"groupMsgs"`
}

// GetGroupChatHistory gets group chat history from cloud (newest first), paging via lastMsgId
// until count messages or no more. count 0 means 50. A leading 'g' on groupID is stripped.
func (a *API) GetGroupChatHistory(ctx context.Context, groupID string, count int) (*GetGroupChatHistoryResponse, error) {
	if count == 0 {
		count = 50
	}
	groupID = strings.TrimPrefix(groupID, "g")
	base := a.svc("group_cloud_message")
	if base == "" {
		return nil, newError("group_cloud_message service not available")
	}
	serviceURL := a.MakeURL(base+"/api/cm/getrecentv2", map[string]any{"nretry": 0}, true)

	// ids are numbers in TS but may arrive as strings; json.Number accepts both.
	type page struct {
		Error, LastMsgID, MinMsgID, MaxMsgID, MsgJumpID, RootMsgID, TsJoinGroup        json.Number
		HasMore, IsOld, IsFiltered, IsRootDel, IsFilteredByPhase, IsFilteredByTimeJoin bool
		GroupMsgs                                                                      []MessageData `json:"groupMsgs"`
	}
	var (
		last   page
		msgs   []MessageData
		seen   = map[string]bool{}
		cursor int64
	)
	for len(msgs) < count {
		raw, err := call[json.RawMessage](ctx, a.Session, http.MethodGet, serviceURL, map[string]any{
			"groupId": groupID, "globalMsgId": cursor, "count": min(50, count-len(msgs)),
			"msgIds": []any{}, "imei": a.IMEI, "src": 3,
		})
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
		if last, err = decodeData[page](raw); err != nil {
			return nil, err
		}
		if last.GroupMsgs == nil {
			break
		}
		for _, m := range last.GroupMsgs {
			if len(msgs) >= count {
				break
			}
			if seen[m.MsgID] {
				continue
			}
			seen[m.MsgID] = true
			msgs = append(msgs, m)
		}
		next, _ := last.LastMsgID.Int64()
		if !last.HasMore || next == 0 || next == cursor {
			break
		}
		cursor = next
	}

	num := func(n json.Number) int64 { v, _ := n.Int64(); return v }
	out := &GetGroupChatHistoryResponse{
		Error: num(last.Error), LastMsgID: num(last.LastMsgID), MinMsgID: num(last.MinMsgID), MaxMsgID: num(last.MaxMsgID),
		MsgJumpID: num(last.MsgJumpID), HasMore: last.HasMore, IsOld: last.IsOld, IsFiltered: last.IsFiltered,
		RootMsgID: num(last.RootMsgID), IsRootDel: last.IsRootDel, TsJoinGroup: num(last.TsJoinGroup),
		IsFilteredByPhase: last.IsFilteredByPhase, IsFilteredByTimeJoin: last.IsFilteredByTimeJoin,
	}
	for _, m := range msgs {
		out.GroupMsgs = append(out.GroupMsgs, NewGroupMessage(a.UID, m))
	}
	return out, nil
}
