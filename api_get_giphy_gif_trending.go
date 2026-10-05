package zca

import (
	"context"
	"net/http"
)

type GetGiphyGifTrendingResponse struct {
	DataList []GiphyGif `json:"data_list"`
	HasMore  int64      `json:"has_more"`
}

// GetGiphyGifTrending gets trending giphy gifs. limit 0 means 10 (limit appears to have no effect).
func (a *API) GetGiphyGifTrending(ctx context.Context, offset, limit int) (*GetGiphyGifTrendingResponse, error) {
	if limit == 0 {
		limit = 10
	}
	return call[*GetGiphyGifTrendingResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("sticker")+"/api/message/giphy/trending", nil, true),
		map[string]any{"offset": offset, "limit": limit, "imei": a.IMEI})
}
