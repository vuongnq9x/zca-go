package zca

import (
	"context"
	"net/http"
	"strings"
)

// FindUser finds a user by phone number. avatarSize 0 means AvatarSizeLarge.
func (a *API) FindUser(ctx context.Context, phoneNumber string, avatarSize AvatarSize) (*UserBasic, error) {
	if phoneNumber == "" {
		return nil, newError("Missing phoneNumber")
	}
	if avatarSize == 0 {
		avatarSize = AvatarSizeLarge
	}
	if strings.HasPrefix(phoneNumber, "0") && a.Language == "vi" {
		phoneNumber = "84" + phoneNumber[1:]
	}
	params := map[string]any{
		"phone":       phoneNumber,
		"avatar_size": avatarSize,
		"language":    a.Language,
		"imei":        a.IMEI,
		"reqSrc":      40,
	}
	enc, err := a.EncodeAES(mustJSON(params))
	if err != nil {
		return nil, newError("Failed to encrypt message")
	}
	// ponytail: zca-js appends "params" twice (once on the URL, once via makeURL); kept as-is.
	base := a.MakeURL(a.svc("friend")+"/api/friend/profile/get", map[string]any{"params": enc}, true)
	resp, err := a.Request(ctx, http.MethodGet, a.MakeURL(base, map[string]any{"params": enc}, true), nil, nil)
	if err != nil {
		return nil, err
	}
	raw, err := a.Resolve(resp, true)
	if err != nil {
		return nil, err
	}
	return decodeData[*UserBasic](raw)
}
