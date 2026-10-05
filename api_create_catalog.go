package zca

import (
	"context"
	"net/http"
)

type CreateCatalogResponse struct {
	Item             CatalogItem `json:"item"`
	VersionLsCatalog int64       `json:"version_ls_catalog"`
	VersionCatalog   int64       `json:"version_catalog"`
}

// CreateCatalog creates a product catalog (zBusiness).
func (a *API) CreateCatalog(ctx context.Context, catalogName string) (*CreateCatalogResponse, error) {
	return call[*CreateCatalogResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("catalog")+"/api/prodcatalog/catalog/create", nil, true),
		map[string]any{"catalog_name": catalogName, "catalog_photo": ""})
}
