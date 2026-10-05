package zca

import (
	"context"
	"net/http"
)

type GetCategoryDetailResponse = *CategoryDetail

// GetCategoryDetail gets info of a sticker category (not its stickers; see GetStickerCategoryDetail).
func (a *API) GetCategoryDetail(ctx context.Context, categoryID int64) (GetCategoryDetailResponse, error) {
	return call[GetCategoryDetailResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("sticker")+"/api/message/sticker/category/detail", nil, true), map[string]any{"cid": categoryID})
}
