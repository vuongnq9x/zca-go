package zca

import (
	"context"
	"net/http"
	"strconv"
)

type CreateGroupResponse struct {
	GroupType     int64          `json:"groupType"`
	SucessMembers []string       `json:"sucessMembers"`
	GroupID       string         `json:"groupId"`
	ErrorMembers  []string       `json:"errorMembers"`
	ErrorData     map[string]any `json:"error_data"`
}

type CreateGroupOptions struct {
	// Group name; empty means the current timestamp.
	Name string
	// Member IDs to add to the group.
	Members []string
	// Optional avatar (file path via Path, or in-memory Data).
	AvatarSource *AttachmentSource
}

// CreateGroup creates a new group. Avatar upload failures are logged, not returned.
func (a *API) CreateGroup(ctx context.Context, options CreateGroupOptions) (*CreateGroupResponse, error) {
	if len(options.Members) == 0 {
		return nil, newError("Group must have at least one member")
	}
	types := make([]int, len(options.Members))
	for i := range types {
		types[i] = -1
	}
	params := map[string]any{
		"clientId":     nowMs(),
		"gname":        strconv.FormatInt(nowMs(), 10),
		"gdesc":        nil,
		"members":      options.Members,
		"membersTypes": types,
		"nameChanged":  0,
		"createLink":   1,
		"clientLang":   a.Language,
		"imei":         a.IMEI,
		"zsource":      601,
	}
	if options.Name != "" {
		params["gname"], params["nameChanged"] = options.Name, 1
	}
	enc, err := a.EncodeAES(mustJSON(params))
	if err != nil {
		return nil, newError("Failed to encrypt message")
	}
	// The params go in the query string even though the method is POST.
	u := a.MakeURL(a.svc("group")+"/api/group/create/v2", map[string]any{"params": enc}, true)
	resp, err := a.Request(ctx, http.MethodPost, u, nil, nil)
	if err != nil {
		return nil, err
	}
	raw, err := a.Resolve(resp, true)
	if err != nil {
		return nil, err
	}
	data, err := decodeData[*CreateGroupResponse](raw)
	if err != nil {
		return nil, err
	}
	if options.AvatarSource != nil && data != nil {
		if _, err := a.ChangeGroupAvatar(ctx, *options.AvatarSource, data.GroupID); err != nil {
			a.log().Error("createGroup: change avatar failed", "err", err)
		}
	}
	return data, nil
}
