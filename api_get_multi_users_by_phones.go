package zca

import (
	"context"
	"net/http"
	"strings"
)

// GetMultiUsersByPhonesResponse maps phone number -> user.
type GetMultiUsersByPhonesResponse map[string]UserBasic

// GetMultiUsersByPhones gets users by phone numbers. avatarSize 0 means AvatarSizeLarge.
func (a *API) GetMultiUsersByPhones(ctx context.Context, phoneNumbers []string, avatarSize AvatarSize) (GetMultiUsersByPhonesResponse, error) {
	if len(phoneNumbers) == 0 {
		return nil, newError("Missing phoneNumbers")
	}
	if avatarSize == 0 {
		avatarSize = AvatarSizeLarge
	}
	phones := make([]string, len(phoneNumbers))
	for i, p := range phoneNumbers {
		if strings.HasPrefix(p, "0") && a.Language == "vi" {
			p = "84" + p[1:]
		}
		phones[i] = p
	}
	params := map[string]any{"phones": phones, "avatar_size": avatarSize, "language": a.Language}
	return call[GetMultiUsersByPhonesResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("friend")+"/api/friend/profile/multiget", nil, true), params)
}
