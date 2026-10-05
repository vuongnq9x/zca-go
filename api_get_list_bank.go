package zca

import (
	"context"
	"net/http"
)

type GetListBankResponse struct {
	Banks []BankInfo `json:"banks"`
}

// GetListBank gets the list of supported banks.
func (a *API) GetListBank(ctx context.Context) (*GetListBankResponse, error) {
	return call[*GetListBankResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("zimsg")+"/api/transfer/conf", nil, true), map[string]any{})
}
