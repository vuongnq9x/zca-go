package zca

type AutoReplyItem struct {
	ID           int64          `json:"id"`
	Weight       int64          `json:"weight"`
	Enable       bool           `json:"enable"`
	ModifiedTime int64          `json:"modifiedTime"`
	StartTime    int64          `json:"startTime"`
	EndTime      int64          `json:"endTime"`
	Content      string         `json:"content"`
	Scope        AutoReplyScope `json:"scope"`
	UIDs         []string       `json:"uids"`
	OwnerID      int64          `json:"ownerId"`
	Recurrence   []string       `json:"recurrence"`
	CreatedTime  int64          `json:"createdTime"`
}

type AutoReplyScope int

const (
	AutoReplyScopeEveryone        AutoReplyScope = 0
	AutoReplyScopeStranger        AutoReplyScope = 1
	AutoReplyScopeSpecificFriends AutoReplyScope = 2
	AutoReplyScopeFriendsExcept   AutoReplyScope = 3
)
