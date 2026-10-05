package zca

import (
	"context"
	"net/http"
)

type DeleteBankAccountPayload struct {
	AccountID int64
	IsDefault bool
}

type DeleteBankAccountResponse struct {
	HasMore bool          `json:"hasMore"`
	Total   int64         `json:"total"`
	MyBanks []BankAccount `json:"myBanks"`
}

// DeleteBankAccount deletes a bank account. Zalo codes: 114 invalid params,
// -265 account not exists, 810 internal error (e.g. deleting the default account).
func (a *API) DeleteBankAccount(ctx context.Context, payload DeleteBankAccountPayload) (*DeleteBankAccountResponse, error) {
	return call[*DeleteBankAccountResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("zimsg")+"/api/transfer/delete", nil, true),
		map[string]any{"account_id": payload.AccountID, "is_default": payload.IsDefault, "language": a.Language})
}
