package zca

import (
	"context"
	"net/http"
)

type StickerSuggestions struct {
	SuggSticker []StickerBasic `json:"sugg_sticker"`
	SuggGuggy   []StickerBasic `json:"sugg_guggy"`
	SuggGif     []StickerBasic `json:"sugg_gif"`
}

// GetStickers searches stickers by keyword and returns sticker IDs.
func (a *API) GetStickers(ctx context.Context, keyword string) ([]int64, error) {
	if keyword == "" {
		return nil, newError("Missing keyword")
	}
	params := map[string]any{"keyword": keyword, "gif": 1, "guggy": 0, "imei": a.IMEI}
	s, err := call[*StickerSuggestions](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("sticker")+"/api/message/sticker/suggest/stickers", nil, true), params)
	if err != nil {
		return nil, err
	}
	ids := []int64{}
	if s != nil {
		// ponytail: TS also ignores sugg_guggy / sugg_gif.
		for _, st := range s.SuggSticker {
			ids = append(ids, st.StickerID)
		}
	}
	return ids, nil
}
