package zca

import (
	"context"
	"net/http"
)

// DeleteCatalog deletes a catalog (zBusiness).
func (a *API) DeleteCatalog(ctx context.Context, catalogID string) (string, error) {
	return call[string](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("catalog")+"/api/prodcatalog/catalog/delete", nil, true), map[string]any{"catalog_id": catalogID})
}
