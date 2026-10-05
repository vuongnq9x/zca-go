package zca

import (
	"context"
	"encoding/json"
	"net/http"
)

type GetGroupInviteBoxInfoPayload struct {
	GroupID string
	MPage   int // default 1 when 0
	MCount  int // default 10 when 0
}

// GetGroupInviteBoxInfoGroupInfo is TS GroupInfo & { topic?: Omit<GroupTopic, "action"> }.
type GetGroupInviteBoxInfoGroupInfo struct {
	GroupInfo
	// Topic.Params is decoded JSON (map[string]any) with params.extra also decoded when it was a string.
	Topic *GroupTopic `json:"topic,omitempty"`
}

type GetGroupInviteBoxInfoResponse struct {
	GroupInfo     GetGroupInviteBoxInfoGroupInfo `json:"groupInfo"`
	InviterInfo   GroupCurrentMem                `json:"inviterInfo"`
	GrCreatorInfo GroupCurrentMem                `json:"grCreatorInfo"`
	ExpiredTs     string                         `json:"expiredTs"`
	Type          int64                          `json:"type"`
}

// GetGroupInviteBoxInfo gets info of a group invitation.
func (a *API) GetGroupInviteBoxInfo(ctx context.Context, payload GetGroupInviteBoxInfoPayload) (*GetGroupInviteBoxInfoResponse, error) {
	if payload.MCount == 0 {
		payload.MCount = 10
	}
	if payload.MPage == 0 {
		payload.MPage = 1
	}
	params := map[string]any{"grId": payload.GroupID, "mcount": payload.MCount, "mpage": payload.MPage}
	data, err := call[*GetGroupInviteBoxInfoResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group")+"/api/group/inv-box/inv-info", nil, true), params)
	if err != nil || data == nil || data.GroupInfo.Topic == nil {
		return data, err
	}
	return data, getGroupInviteBoxInfoParseTopic(data.GroupInfo.Topic)
}

// getGroupInviteBoxInfoParseTopic decodes topic.params (and params.extra) when they are JSON strings.
func getGroupInviteBoxInfoParseTopic(t *GroupTopic) error {
	s, ok := t.Params.(string)
	if !ok {
		return nil
	}
	var params any
	if err := json.Unmarshal([]byte(s), &params); err != nil {
		return err
	}
	if m, ok := params.(map[string]any); ok {
		if extra, ok := m["extra"].(string); ok {
			var v any
			if err := json.Unmarshal([]byte(extra), &v); err != nil {
				return err
			}
			m["extra"] = v
		}
	}
	t.Params = params
	return nil
}
