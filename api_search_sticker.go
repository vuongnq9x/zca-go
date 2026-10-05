package zca

import (
	"context"
	"net/http"
)

type SearchStickerResponse []StickerBasic

// SearchSticker searches stickers by keyword. limit 0 means 50.
func (a *API) SearchSticker(ctx context.Context, keyword string, limit int) (SearchStickerResponse, error) {
	if limit == 0 {
		limit = 50
	}
	params := map[string]any{"keyword": keyword, "limit": limit, "srcType": 0, "imei": a.IMEI}
	return call[SearchStickerResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("sticker")+"/api/message/sticker/search", nil, true), params)
}
