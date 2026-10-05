package zca

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

type SendBankCardPayload struct {
	BinBank     BinBankCard
	NumAccBank  string
	NameAccBank string // optional, upper-cased; "---" when empty
}

// SendBankCard sends a bank account card to a thread.
func (a *API) SendBankCard(ctx context.Context, payload SendBankCardPayload, threadID string, threadType ThreadType) (string, error) {
	name := strings.ToUpper(payload.NameAccBank)
	if name == "" {
		name = "---"
	}
	destType := 0
	if threadType == ThreadTypeGroup {
		destType = 1
	}
	now := nowMs()
	params := map[string]any{
		"binBank":     payload.BinBank,
		"numAccBank":  payload.NumAccBank,
		"nameAccBank": name,
		"cliMsgId":    strconv.FormatInt(now, 10),
		"tsMsg":       now,
		"destUid":     threadID,
		"destType":    destType,
	}
	return call[string](ctx, a.Session, http.MethodPost, a.MakeURL(a.svc("zimsg")+"/api/transfer/card", nil, true), params)
}
