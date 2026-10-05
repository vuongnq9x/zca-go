package zca

import (
	"context"
	"net/http"
)

type DeleteProductCatalogPayload struct {
	ProductIDs []string
	CatalogID  string
}

type DeleteProductCatalogResponse struct {
	Item             []int64 `json:"item"`
	VersionLsCatalog int64   `json:"version_ls_catalog"`
	VersionCatalog   int64   `json:"version_catalog"`
}

// DeleteProductCatalog deletes product(s) from a catalog (zBusiness).
func (a *API) DeleteProductCatalog(ctx context.Context, payload DeleteProductCatalogPayload) (*DeleteProductCatalogResponse, error) {
	return call[*DeleteProductCatalogResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("catalog")+"/api/prodcatalog/product/mdelete", nil, true),
		map[string]any{"product_ids": payload.ProductIDs, "catalog_id": payload.CatalogID})
}
