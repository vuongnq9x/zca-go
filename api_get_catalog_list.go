package zca

import (
	"context"
	"net/http"
)

type GetCatalogListPayload struct {
	Limit         int   // default 20 when 0
	LastProductID int64 // default -1 when 0
	Page          int   // default 0
}

type GetCatalogListResponse struct {
	Items   []CatalogItem `json:"items"`
	Version int64         `json:"version"`
	HasMore int64         `json:"has_more"`
}

// GetCatalogList gets the catalog list (zBusiness).
func (a *API) GetCatalogList(ctx context.Context, payload GetCatalogListPayload) (*GetCatalogListResponse, error) {
	if payload.Limit == 0 {
		payload.Limit = 20
	}
	if payload.LastProductID == 0 {
		payload.LastProductID = -1
	}
	params := map[string]any{
		"version_list_catalog": 0,
		"limit":                payload.Limit,
		"last_product_id":      payload.LastProductID,
		"page":                 payload.Page,
	}
	return call[*GetCatalogListResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("catalog")+"/api/prodcatalog/catalog/list", nil, true), params)
}
