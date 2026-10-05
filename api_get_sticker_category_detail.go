package zca

import (
	"context"
	"net/http"
)

type GetStickerCategoryDetailResponse = []StickerDetail

// GetStickerCategoryDetail gets stickers of a sticker category.
func (a *API) GetStickerCategoryDetail(ctx context.Context, cateID int64) (GetStickerCategoryDetailResponse, error) {
	return call[GetStickerCategoryDetailResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("sticker")+"/api/message/sticker/category/sticker_detail", nil, true), map[string]any{"cid": cateID})
}
