package zca

type LabelData struct {
	ID            int64    `json:"id"`
	Text          string   `json:"text"`
	TextKey       string   `json:"textKey"`
	Conversations []string `json:"conversations"`
	Color         string   `json:"color"`
	Offset        int64    `json:"offset"`
	Emoji         string   `json:"emoji"`
	CreateTime    int64    `json:"createTime"`
}
