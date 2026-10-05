package zca

import (
	"context"
	"net/http"
)

type UpdateSettingsType string

const (
	UpdateSettingsTypeViewBirthday             UpdateSettingsType = "view_birthday"
	UpdateSettingsTypeShowOnlineStatus         UpdateSettingsType = "show_online_status"
	UpdateSettingsTypeDisplaySeenStatus        UpdateSettingsType = "display_seen_status"
	UpdateSettingsTypeReceiveMessage           UpdateSettingsType = "receive_message"
	UpdateSettingsTypeAcceptCall               UpdateSettingsType = "accept_stranger_call"
	UpdateSettingsTypeAddFriendViaPhone        UpdateSettingsType = "add_friend_via_phone"
	UpdateSettingsTypeAddFriendViaQR           UpdateSettingsType = "add_friend_via_qr"
	UpdateSettingsTypeAddFriendViaGroup        UpdateSettingsType = "add_friend_via_group"
	UpdateSettingsTypeAddFriendViaContact      UpdateSettingsType = "add_friend_via_contact"
	UpdateSettingsTypeDisplayOnRecommendFriend UpdateSettingsType = "display_on_recommend_friend"
	UpdateSettingsTypeArchivedChat             UpdateSettingsType = "archivedChatStatus"
	UpdateSettingsTypeQuickMessage             UpdateSettingsType = "quickMessageStatus"
)

// UpdateSettings sets an account setting.
//
// Values: ViewBirthday 0 hide, 1 full date, 2 day/month; ReceiveMessage 1 everyone, 2 friends only;
// AcceptCall 2 friends only, 3 everyone, 4 friends and contacted people; all others 0 off/hide, 1 on/show.
func (a *API) UpdateSettings(ctx context.Context, settingType UpdateSettingsType, value int) (string, error) {
	return call[string](ctx, a.Session, http.MethodGet,
		a.MakeURL("https://wpa.chat.zalo.me/api/setting/update", nil, true), map[string]any{string(settingType): value})
}
