package zca

import (
	"context"
	"net/http"
)

type UpdateBankAccountPayload struct {
	AccountID   int64
	BinBank     BinBankCard
	NumAccBank  string
	NameAccBank string
}

type UpdateBankAccountResponse = BankAccount

// UpdateBankAccount updates a bank account.
func (a *API) UpdateBankAccount(ctx context.Context, payload UpdateBankAccountPayload) (*UpdateBankAccountResponse, error) {
	params := map[string]any{
		"account_id": payload.AccountID, "bin": payload.BinBank, "bank_number": payload.NumAccBank, "language": a.Language,
	}
	if name := NormalizeHolderName(payload.NameAccBank); name != "" {
		params["holder_name"] = name
	}
	return call[*UpdateBankAccountResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("zimsg")+"/api/transfer/update", nil, true), params)
}
