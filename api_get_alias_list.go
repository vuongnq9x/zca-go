package zca

import (
	"context"
	"net/http"
)

type GetAliasListResponse struct {
	Items []struct {
		UserID string `json:"userId"`
		Alias  string `json:"alias"`
	} `json:"items"`
	UpdateTime string `json:"updateTime"`
}

// GetAliasList gets the alias list. count 0 means 100, page 0 means 1.
func (a *API) GetAliasList(ctx context.Context, count, page int) (*GetAliasListResponse, error) {
	if count == 0 {
		count = 100
	}
	if page == 0 {
		page = 1
	}
	params := map[string]any{"page": page, "count": count, "imei": a.IMEI}
	return call[*GetAliasListResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("alias")+"/api/alias/list", nil, true), params)
}
