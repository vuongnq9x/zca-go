package zca

import (
	"context"
	"net/http"
)

type SearchGiphyGifResponse = GetGiphyGifTrendingResponse

// SearchGiphyGif searches giphy gifs (TS: searchGiphyGifs). limit 0 means 10 (limit appears to have no effect).
func (a *API) SearchGiphyGif(ctx context.Context, keyword string, offset, limit int) (*SearchGiphyGifResponse, error) {
	if limit == 0 {
		limit = 10
	}
	return call[*SearchGiphyGifResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("sticker")+"/api/message/giphy/search", nil, true),
		map[string]any{"keyword": keyword, "offset": offset, "limit": limit})
}
