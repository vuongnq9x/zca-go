package zca

import (
	"context"
	"net/http"
	"strings"
)

type UserInfoResponse struct {
	UnchangedProfiles map[string]UnchangedProfileInfo `json:"unchanged_profiles"`
	PhonebookVersion  int64                           `json:"phonebook_version"`
	ChangedProfiles   map[string]User                 `json:"changed_profiles"`
}

// GetUserInfo gets profiles of one or more users. avatarSize 0 means AvatarSizeSmall.
func (a *API) GetUserInfo(ctx context.Context, userIDs []string, avatarSize AvatarSize) (*UserInfoResponse, error) {
	if len(userIDs) == 0 {
		return nil, newError("Missing user id")
	}
	if avatarSize == 0 {
		avatarSize = AvatarSizeSmall
	}
	ids := make([]string, len(userIDs))
	for i, id := range userIDs {
		if !strings.Contains(id, "_") {
			id += "_0"
		}
		ids[i] = id
	}
	params := map[string]any{
		"phonebook_version":   a.ExtraVer["phonebook"],
		"friend_pversion_map": ids,
		"avatar_size":         avatarSize,
		"language":            a.Language,
		"show_online_status":  1,
		"imei":                a.IMEI,
	}
	return call[*UserInfoResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("profile")+"/api/social/friend/getprofiles/v2", nil, true), params)
}
