package zca

import (
	"context"
	"encoding/json"
	"net/http"
)

type ListBoardOptions struct {
	Page  int // default 1 when 0
	Count int // default 20 when 0
}

type BoardItem struct {
	BoardType BoardType `json:"boardType"`
	// Data is *PollDetail, *NoteDetail or *PinnedMessageDetail by BoardType
	// (map[string]any for unknown types).
	Data any `json:"data"`
}

type GetListBoardResponse struct {
	Items []BoardItem `json:"items"`
	Count int64       `json:"count"`
}

// GetListBoard gets board items (notes, pinned messages, polls) of a group.
func (a *API) GetListBoard(ctx context.Context, options ListBoardOptions, groupID string) (*GetListBoardResponse, error) {
	if options.Page == 0 {
		options.Page = 1
	}
	if options.Count == 0 {
		options.Count = 20
	}
	params := map[string]any{
		"group_id":   groupID,
		"board_type": 0,
		"page":       options.Page,
		"count":      options.Count,
		"last_id":    0,
		"last_type":  0,
		"imei":       a.IMEI,
	}
	raw, err := call[json.RawMessage](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group_board")+"/api/board/list", nil, true), params)
	if err != nil {
		return nil, err
	}
	return getListBoardDecode(raw)
}

// getListBoardDecode parses string params of non-poll items and decodes each item by board type.
func getListBoardDecode(raw json.RawMessage) (*GetListBoardResponse, error) {
	data, err := decodeData[*struct {
		Items []struct {
			BoardType BoardType                  `json:"boardType"`
			Data      map[string]json.RawMessage `json:"data"`
		} `json:"items"`
		Count int64 `json:"count"`
	}](raw)
	if err != nil || data == nil {
		return nil, err
	}
	out := &GetListBoardResponse{Count: data.Count, Items: make([]BoardItem, len(data.Items))}
	for i, it := range data.Items {
		if p := it.Data["params"]; it.BoardType != BoardTypePoll && len(p) > 0 && p[0] == '"' {
			var s string
			if err := json.Unmarshal(p, &s); err != nil {
				return nil, err
			}
			it.Data["params"] = json.RawMessage(s)
		}
		var v any
		switch it.BoardType {
		case BoardTypePoll:
			v = &PollDetail{}
		case BoardTypeNote:
			v = &NoteDetail{}
		case BoardTypePinnedMessage:
			v = &PinnedMessageDetail{}
		default:
			v = &map[string]any{}
		}
		b, err := json.Marshal(it.Data)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(b, v); err != nil {
			return nil, err
		}
		if m, ok := v.(*map[string]any); ok {
			v = *m
		}
		out.Items[i] = BoardItem{BoardType: it.BoardType, Data: v}
	}
	return out, nil
}
