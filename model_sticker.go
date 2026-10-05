package zca

type StickerDetail struct {
	ID               int64   `json:"id"`
	CateID           int64   `json:"cateId"`
	Type             int64   `json:"type"`
	Text             string  `json:"text"`
	URI              string  `json:"uri"`
	Fkey             int64   `json:"fkey"`
	Status           int64   `json:"status"`
	StickerURL       string  `json:"stickerUrl"`
	StickerSpriteURL string  `json:"stickerSpriteUrl"`
	StickerWebpURL   *string `json:"stickerWebpUrl"`
	TotalFrames      int64   `json:"totalFrames"`
	Duration         int64   `json:"duration"`
	EffectID         int64   `json:"effectId"`
	Checksum         string  `json:"checksum"`
	Ext              int64   `json:"ext"`
	Source           int64   `json:"source"`
	Fss              any     `json:"fss"`
	FssInfo          any     `json:"fssInfo"`
	Version          int64   `json:"version"`
	ExtInfo          any     `json:"extInfo"`
}

type StickerBasic struct {
	Type      int64 `json:"type"`
	CateID    int64 `json:"cate_id"`
	StickerID int64 `json:"sticker_id"`
}
