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

type TenorSticker struct {
	ID  string `json:"id"`
	CID int64  `json:"cid"`
	EID int64  `json:"eid"`
}

type GifMedia struct {
	Width  int64  `json:"width"`
	Height int64  `json:"height"`
	URL    string `json:"url"`
}

type CategoryDetail struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Desc        string `json:"desc"`
	TotalImage  int64  `json:"totalImage"`
	ThumbURL    string `json:"thumbUrl"`
	IconURL     string `json:"iconUrl"`
	IconPreview string `json:"iconPreview"`
	Price       int64  `json:"price"`
	Group       int64  `json:"group"`
	Status      int64  `json:"status"`
	Version     int64  `json:"version"`
	ThumbImg    string `json:"thumbImg"`
	Source      any    `json:"source"`
	Type        int64  `json:"type"`
	SourceURL   string `json:"sourceUrl"`
	Permission  int64  `json:"permission"`
	ExpireTime  int64  `json:"expireTime"`
	IsHidden    int64  `json:"is_hidden"`
	Order       int64  `json:"order"`
}

type GiphyGif struct {
	ID       string   `json:"id"`
	Original GifMedia `json:"original"`
	Normal   GifMedia `json:"normal"`
	HD       GifMedia `json:"hd"`
	Preview  GifMedia `json:"preview"`
}
