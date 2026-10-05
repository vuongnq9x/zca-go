package zca

import (
	"context"
	"net/http"
)

type GetListBankAccountResponse struct {
	HasMore bool          `json:"hasMore"`
	Total   int64         `json:"total"`
	MyBanks []BankAccount `json:"myBanks"`
}

// GetListBankAccount gets the user's bank accounts. page defaults to 0, limit 0 means 20.
// Error code 114 means invalid params.
func (a *API) GetListBankAccount(ctx context.Context, page, limit int) (*GetListBankAccountResponse, error) {
	if limit == 0 {
		limit = 20
	}
	return call[*GetListBankAccountResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("zimsg")+"/api/transfer/list", nil, true), map[string]any{"page": page, "limit": limit})
}
