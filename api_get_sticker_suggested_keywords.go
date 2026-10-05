package zca

import (
	"context"
	"net/http"
)

type GetStickerSuggestedKeywordsResponse struct {
	Keywords    []string `json:"keywords"`
	ExpiredTime int64    `json:"expired_time"`
	WordSearch  int64    `json:"word_search"`
}

// GetStickerSuggestedKeywords gets suggested sticker keywords.
func (a *API) GetStickerSuggestedKeywords(ctx context.Context) (*GetStickerSuggestedKeywordsResponse, error) {
	return call[*GetStickerSuggestedKeywordsResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("sticker")+"/api/message/sticker/suggest/keywords", nil, true), map[string]any{"imei": a.IMEI})
}
