package zca

import (
	"context"
	"net/http"
)

type GiphyCategory struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	EnName   string   `json:"enName"`
	Icon     string   `json:"icon"`
	Type     int64    `json:"type"`
	Keywords []string `json:"keywords"`
}

type GetGiphyCategoriesResponse struct {
	DataList []GiphyCategory `json:"data_list"`
	HasMore  int64           `json:"has_more"`
}

// GetGiphyCategories gets giphy categories. limit 0 means 10 (limit appears to have no effect).
func (a *API) GetGiphyCategories(ctx context.Context, offset, limit int) (*GetGiphyCategoriesResponse, error) {
	if limit == 0 {
		limit = 10
	}
	return call[*GetGiphyCategoriesResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("sticker")+"/api/message/giphy/cates", nil, true),
		map[string]any{"offset": offset, "limit": limit})
}
