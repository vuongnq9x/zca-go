package zca

import (
	"context"
	"net/http"
)

type CreateBankAccountPayload struct {
	BinBank     BinBankCard
	NumAccBank  string
	NameAccBank string
}

type CreateBankAccountResponse = BankAccount

// CreateBankAccount adds a bank account. Zalo codes: -263 already exists, 810 internal error/invalid input.
func (a *API) CreateBankAccount(ctx context.Context, payload CreateBankAccountPayload) (*CreateBankAccountResponse, error) {
	params := map[string]any{
		"bin":         payload.BinBank,
		"bank_number": payload.NumAccBank,
		"language":    a.Language,
	}
	if name := NormalizeHolderName(payload.NameAccBank); name != "" {
		params["holder_name"] = name
	}
	return call[*CreateBankAccountResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("zimsg")+"/api/transfer/create", nil, true), params)
}
