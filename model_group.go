package zca

type GroupSetting struct {
	BlockName        int64 `json:"blockName"`
	SignAdminMsg     int64 `json:"signAdminMsg"`
	AddMemberOnly    int64 `json:"addMemberOnly"`
	SetTopicOnly     int64 `json:"setTopicOnly"`
	EnableMsgHistory int64 `json:"enableMsgHistory"`
	JoinAppr         int64 `json:"joinAppr"`
	LockCreatePost   int64 `json:"lockCreatePost"`
	LockCreatePoll   int64 `json:"lockCreatePoll"`
	LockSendMsg      int64 `json:"lockSendMsg"`
	LockViewMember   int64 `json:"lockViewMember"`
	BannFeature      int64 `json:"bannFeature"`
	DirtyMedia       int64 `json:"dirtyMedia"`
	BanDuration      int64 `json:"banDuration"`
}

type GroupTopicType int

const (
	GroupTopicTypeNote    GroupTopicType = 0
	GroupTopicTypeMessage GroupTopicType = 2
	GroupTopicTypePoll    GroupTopicType = 3
)

type GroupTopicNoteParams struct {
	ClientMsgID string `json:"client_msg_id"`
	GlobalMsgID string `json:"global_msg_id"`
	Title       string `json:"title"`
}

// GroupTopicMessageParams merges the TS text/voice/image/video/file/gif message params,
// discriminated by MsgType (1 text, 31 voice, 32 image, 44 video, 46 file, 49 gif).
type GroupTopicMessageParams struct {
	SenderUID   string               `json:"senderUid"`
	SenderName  string               `json:"senderName"`
	ClientMsgID string               `json:"client_msg_id"`
	GlobalMsgID string               `json:"global_msg_id"`
	MsgType     int64                `json:"msg_type"`
	Title       string               `json:"title"`
	Thumb       string               `json:"thumb,omitempty"`
	Extra       *GroupTopicFileExtra `json:"extra,omitempty"`
}

type GroupTopicFileExtra struct {
	FileSize    string `json:"fileSize"`
	Checksum    string `json:"checksum"`
	ChecksumSha any    `json:"checksumSha"`
	FileExt     string `json:"fileExt"`
	Fdata       string `json:"fdata"`
	FType       int64  `json:"fType"`
}

type GroupTopicPollParams struct {
	PollID int64  `json:"pollId"`
	Title  string `json:"title"`
}

type GroupTopic struct {
	Type      GroupTopicType `json:"type"`
	Color     int64          `json:"color"`
	Emoji     string         `json:"emoji"`
	StartTime int64          `json:"startTime"`
	Duration  int64          `json:"duration"`
	// Params is note/message/poll params or any other object; decoded generically.
	Params     any    `json:"params"`
	ID         string `json:"id"`
	CreatorID  string `json:"creatorId"`
	CreateTime int64  `json:"createTime"`
	EditorID   string `json:"editorId"`
	EditTime   int64  `json:"editTime"`
	Repeat     int64  `json:"repeat"`
	Action     int64  `json:"action"`
}

type GroupType int

const (
	GroupTypeGroup     GroupType = 1
	GroupTypeCommunity GroupType = 2
)

type GroupCurrentMem struct {
	ID            string `json:"id"`
	DName         string `json:"dName"`
	ZaloName      string `json:"zaloName"`
	Avatar        string `json:"avatar"`
	Avatar25      string `json:"avatar_25"`
	AccountStatus int64  `json:"accountStatus"`
	Type          int64  `json:"type"`
}

type GroupInfoExtraInfo struct {
	EnableMediaStore int64 `json:"enable_media_store"`
}

type GroupInfo struct {
	GroupID       string            `json:"groupId"`
	Name          string            `json:"name"`
	Desc          string            `json:"desc"`
	Type          GroupType         `json:"type"`
	CreatorID     string            `json:"creatorId"`
	Version       string            `json:"version"`
	Avt           string            `json:"avt"`
	FullAvt       string            `json:"fullAvt"`
	MemberIDs     []string          `json:"memberIds"`
	AdminIDs      []string          `json:"adminIds"`
	CurrentMems   []GroupCurrentMem `json:"currentMems"`
	UpdateMems    []any             `json:"updateMems"`
	Admins        []any             `json:"admins"`
	HasMoreMember int64             `json:"hasMoreMember"`
	SubType       int64             `json:"subType"`
	TotalMember   int64             `json:"totalMember"`
	MaxMember     int64             `json:"maxMember"`
	Setting       GroupSetting      `json:"setting"`
	CreatedTime   int64             `json:"createdTime"`
	Visibility    int64             `json:"visibility"`
	GlobalID      string            `json:"globalId"`
	// E2ee: 1 true, 0 false.
	E2ee      int64              `json:"e2ee"`
	ExtraInfo GroupInfoExtraInfo `json:"extraInfo"`
}
