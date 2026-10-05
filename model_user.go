package zca

type User struct {
	UserID         string           `json:"userId"`
	Username       string           `json:"username"`
	DisplayName    string           `json:"displayName"`
	ZaloName       string           `json:"zaloName"`
	Avatar         string           `json:"avatar"`
	Bgavatar       string           `json:"bgavatar"`
	Cover          string           `json:"cover"`
	Gender         Gender           `json:"gender"`
	Dob            int64            `json:"dob"`
	Sdob           string           `json:"sdob"`
	Status         string           `json:"status"`
	PhoneNumber    string           `json:"phoneNumber"`
	IsFr           int64            `json:"isFr"`
	IsBlocked      int64            `json:"isBlocked"`
	LastActionTime int64            `json:"lastActionTime"`
	LastUpdateTime int64            `json:"lastUpdateTime"`
	IsActive       int64            `json:"isActive"`
	Key            int64            `json:"key"`
	Type           int64            `json:"type"`
	IsActivePC     int64            `json:"isActivePC"`
	IsActiveWeb    int64            `json:"isActiveWeb"`
	IsValid        int64            `json:"isValid"`
	UserKey        string           `json:"userKey"`
	AccountStatus  int64            `json:"accountStatus"`
	OaInfo         any              `json:"oaInfo"`
	UserMode       int64            `json:"user_mode"`
	GlobalID       string           `json:"globalId"`
	BizPkg         ZBusinessPackage `json:"bizPkg"`
	CreatedTs      int64            `json:"createdTs"`
	OaStatus       any              `json:"oa_status"`
}

type UserBasic struct {
	Avatar      string           `json:"avatar"`
	Cover       string           `json:"cover"`
	Status      string           `json:"status"`
	Gender      Gender           `json:"gender"`
	Dob         int64            `json:"dob"`
	Sdob        string           `json:"sdob"`
	GlobalID    string           `json:"globalId"`
	BizPkg      ZBusinessPackage `json:"bizPkg"`
	UID         string           `json:"uid"`
	ZaloName    string           `json:"zalo_name"`
	DisplayName string           `json:"display_name"`
}

type UserSetting struct {
	AddFriendViaContact      int64 `json:"add_friend_via_contact"`
	DisplayOnRecommendFriend int64 `json:"display_on_recommend_friend"`
	AddFriendViaGroup        int64 `json:"add_friend_via_group"`
	AddFriendViaQR           int64 `json:"add_friend_via_qr"`
	QuickMessageStatus       int64 `json:"quick_message_status"`
	ShowOnlineStatus         bool  `json:"show_online_status"`
	AcceptStrangerCall       int64 `json:"accept_stranger_call"`
	ArchivedChatStatus       int64 `json:"archived_chat_status"`
	ReceiveMessage           int64 `json:"receive_message"`
	AddFriendViaPhone        int64 `json:"add_friend_via_phone"`
	DisplaySeenStatus        int64 `json:"display_seen_status"`
	ViewBirthday             int64 `json:"view_birthday"`
	Setting2FAStatus         int64 `json:"setting_2FA_status"`
}

type UnchangedProfileInfoOaStatus struct {
	Blocked int64 `json:"blocked"`
}

type UnchangedProfileInfo struct {
	OaStatus       UnchangedProfileInfoOaStatus `json:"oa_status"`
	IsFr           int64                        `json:"isFr"`
	IsBlocked      bool                         `json:"isBlocked"`
	LastActionTime int64                        `json:"lastActionTime"`
}
