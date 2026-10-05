package zca

type BoardType int

const (
	BoardTypeNote          BoardType = 1
	BoardTypePinnedMessage BoardType = 2
	BoardTypePoll          BoardType = 3
)

type PollDetail struct {
	Creator           string        `json:"creator"`
	Question          string        `json:"question"`
	Options           []PollOptions `json:"options"`
	Joined            bool          `json:"joined"`
	Closed            bool          `json:"closed"`
	PollID            int64         `json:"poll_id"`
	AllowMultiChoices bool          `json:"allow_multi_choices"`
	AllowAddNewOption bool          `json:"allow_add_new_option"`
	IsAnonymous       bool          `json:"is_anonymous"`
	PollType          int64         `json:"poll_type"`
	CreatedTime       int64         `json:"created_time"`
	UpdatedTime       int64         `json:"updated_time"`
	ExpiredTime       int64         `json:"expired_time"`
	IsHideVotePreview bool          `json:"is_hide_vote_preview"`
	NumVote           int64         `json:"num_vote"`
}

type PollOptions struct {
	Content  string   `json:"content"`
	Votes    int64    `json:"votes"`
	Voted    bool     `json:"voted"`
	Voters   []string `json:"voters"`
	OptionID int64    `json:"option_id"`
}

type NoteDetailParams struct {
	Title string `json:"title"`
	Extra string `json:"extra,omitempty"`
}

type NoteDetail struct {
	ID         string           `json:"id"`
	Type       int64            `json:"type"`
	Color      int64            `json:"color"`
	Emoji      string           `json:"emoji"`
	StartTime  int64            `json:"startTime"`
	Duration   int64            `json:"duration"`
	Params     NoteDetailParams `json:"params"`
	CreatorID  string           `json:"creatorId"`
	EditorID   string           `json:"editorId"`
	CreateTime int64            `json:"createTime"`
	EditTime   int64            `json:"editTime"`
	Repeat     int64            `json:"repeat"`
}

type PinnedMessageDetail struct {
	ID         string         `json:"id"`
	Type       int64          `json:"type"`
	Color      int64          `json:"color"`
	Emoji      string         `json:"emoji"`
	StartTime  int64          `json:"startTime"`
	Duration   int64          `json:"duration"`
	Params     map[string]any `json:"params"`
	CreatorID  string         `json:"creatorId"`
	EditorID   string         `json:"editorId"`
	CreateTime int64          `json:"createTime"`
	EditTime   int64          `json:"editTime"`
	Repeat     int64          `json:"repeat"`
}
