package zca

import (
	"context"
	"net/http"
)

type UpdateCatalogPayload struct {
	CatalogID   string
	CatalogName string
}

type UpdateCatalogResponse struct {
	Item             CatalogItem `json:"item"`
	VersionLsCatalog int64       `json:"version_ls_catalog"`
	VersionCatalog   int64       `json:"version_catalog"`
}

// UpdateCatalog renames a catalog (zBusiness).
func (a *API) UpdateCatalog(ctx context.Context, payload UpdateCatalogPayload) (*UpdateCatalogResponse, error) {
	params := map[string]any{"catalog_id": payload.CatalogID, "catalog_name": payload.CatalogName, "catalog_photo": ""}
	return call[*UpdateCatalogResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("catalog")+"/api/prodcatalog/catalog/update", nil, true), params)
}
