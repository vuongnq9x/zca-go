package zca

import (
	"context"
	"net/http"
)

type UpdateGroupSettingsOptions struct {
	BlockName        bool // disallow members to change group name and avatar
	SignAdminMsg     bool // highlight messages from owner/admins
	SetTopicOnly     bool // don't pin messages, notes and polls to the top
	EnableMsgHistory bool // allow new members to read most recent messages
	JoinAppr         bool // membership approval
	LockCreatePost   bool // disallow members to create notes & reminders
	LockCreatePoll   bool // disallow members to create polls
	LockSendMsg      bool // disallow members to send messages
	LockViewMember   bool // disallow members to view full member list (community only)
}

// UpdateGroupSettings updates group settings. Zalo error 166 means insufficient permission.
func (a *API) UpdateGroupSettings(ctx context.Context, o UpdateGroupSettingsOptions, groupID string) (string, error) {
	b := func(v bool) int {
		if v {
			return 1
		}
		return 0
	}
	params := map[string]any{
		"blockName": b(o.BlockName), "signAdminMsg": b(o.SignAdminMsg), "setTopicOnly": b(o.SetTopicOnly),
		"enableMsgHistory": b(o.EnableMsgHistory), "joinAppr": b(o.JoinAppr), "lockCreatePost": b(o.LockCreatePost),
		"lockCreatePoll": b(o.LockCreatePoll), "lockSendMsg": b(o.LockSendMsg), "lockViewMember": b(o.LockViewMember),
		"bannFeature": 0, "dirtyMedia": 0, "banDuration": 0, "blocked_members": []string{},
		"grid": groupID, "imei": a.IMEI,
	}
	return call[string](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group")+"/api/group/setting/update", nil, true), params)
}
