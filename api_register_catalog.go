package zca

import (
	"context"
	"net/http"
)

// RegisterCatalog enables or disables the product catalog.
func (a *API) RegisterCatalog(ctx context.Context, enable bool) (string, error) {
	e := 0
	if enable {
		e = 1
	}
	return call[string](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("catalog")+"/api/prodcatalog/catalog/register", nil, true), map[string]any{"enable": e})
}
