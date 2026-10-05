package zca

import (
	"context"
	"net/http"
)

type GetProductCatalogListPayload struct {
	CatalogID      string
	Limit          int64 // default 100 when 0
	VersionCatalog int64
	LastProductID  string // default -1 when ""
	Page           int64
}

type GetProductCatalogListResponse struct {
	Items   []ProductCatalogItem `json:"items"`
	Version int64                `json:"version"`
	HasMore int64                `json:"has_more"`
}

// GetProductCatalogList lists products of a catalog (zBusiness).
func (a *API) GetProductCatalogList(ctx context.Context, payload GetProductCatalogListPayload) (*GetProductCatalogListResponse, error) {
	limit := payload.Limit
	if limit == 0 {
		limit = 100
	}
	var lastProductID any = payload.LastProductID
	if payload.LastProductID == "" {
		lastProductID = -1
	}
	params := map[string]any{
		"catalog_id":      payload.CatalogID,
		"limit":           limit,
		"version_catalog": payload.VersionCatalog,
		"last_product_id": lastProductID,
		"page":            payload.Page,
	}
	return call[*GetProductCatalogListResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("catalog")+"/api/prodcatalog/product/list", nil, true), params)
}
