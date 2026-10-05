package zca

import (
	"context"
	"net/http"
)

type GetTenorStickerMappingResponse struct {
	TenorStickerMap map[string]TenorSticker `json:"tenor_sticker_map"`
	ExpiredTime     int64                   `json:"expired_time"`
}

// GetTenorStickerMapping gets the tenor sticker mapping.
func (a *API) GetTenorStickerMapping(ctx context.Context) (*GetTenorStickerMappingResponse, error) {
	return call[*GetTenorStickerMappingResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("sticker")+"/api/message/sticker/tenor/mapping", nil, true), map[string]any{"imei": a.IMEI})
}
