package zca

import (
	"context"
	"net/http"
)

type GetStickerCategoriesResponse struct {
	DataURL   string  `json:"dataUrl"`
	CS        string  `json:"cs"`
	ExCateIDs []int64 `json:"exCateIds"`
	Data      []any   `json:"data"`
}

// GetStickerCategories gets sticker categories.
func (a *API) GetStickerCategories(ctx context.Context) (*GetStickerCategoriesResponse, error) {
	return call[*GetStickerCategoriesResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("sticker")+"/api/message/sticker/category/list/v2", nil, true), map[string]any{})
}
