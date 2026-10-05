package zca

type QuickMessageContent struct {
	Title  string  `json:"title"`
	Params *string `json:"params"`
}

type QuickMessageMediaItem struct {
	Type         int64  `json:"type"`
	PhotoID      int64  `json:"photoId"`
	Title        string `json:"title"`
	Width        int64  `json:"width"`
	Height       int64  `json:"height"`
	PreviewThumb string `json:"previewThumb"`
	RawURL       string `json:"rawUrl"`
	ThumbURL     string `json:"thumbUrl"`
	NormalURL    string `json:"normalUrl"`
	HdURL        string `json:"hdUrl"`
}

type QuickMessageMedia struct {
	Items []QuickMessageMediaItem `json:"items"`
}

type QuickMessage struct {
	ID           int64               `json:"id"`
	Keyword      string              `json:"keyword"`
	Type         int64               `json:"type"`
	CreatedTime  int64               `json:"createdTime"`
	LastModified int64               `json:"lastModified"`
	Message      QuickMessageContent `json:"message"`
	Media        *QuickMessageMedia  `json:"media"`
}
