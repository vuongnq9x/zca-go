package zca

import (
	"context"
	"net/http"
)

type GiphyGifDictionary struct {
	ID         string   `json:"id"`
	EnKeywords []string `json:"enKeywords"`
	ViKeywords []string `json:"viKeywords"`
}

type GetGiphyGifDictionariesResponse struct {
	DataList []GiphyGifDictionary `json:"data_list"`
	HasMore  int64                `json:"has_more"`
}

// GetGiphyGifDictionaries gets giphy gif dictionaries. limit 0 means 10 (limit appears to have no effect).
func (a *API) GetGiphyGifDictionaries(ctx context.Context, offset, limit int) (*GetGiphyGifDictionariesResponse, error) {
	if limit == 0 {
		limit = 10
	}
	return call[*GetGiphyGifDictionariesResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("sticker")+"/api/message/giphy/dicts", nil, true),
		map[string]any{"offset": offset, "limit": limit, "imei": a.IMEI})
}
