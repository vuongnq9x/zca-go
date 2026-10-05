package zca

import (
	"context"
	"net/http"
	"strconv"
)

type SendCardOptions struct {
	UserID      string
	PhoneNumber string // optional
	TTL         int64  // milliseconds
}

type SendCardResponse struct {
	MsgID StringOrNumber `json:"msgId"`
}

// SendCard sends a user's contact card to a user or group.
func (a *API) SendCard(ctx context.Context, opts SendCardOptions, threadID string, threadType ThreadType) (*SendCardResponse, error) {
	// ponytail: inline mget-qr call; switch to a.GetQR once that API lands.
	qr, err := call[map[string]string](ctx, a.Session, http.MethodPost, a.MakeURL(a.svc("friend")+"/api/friend/mget-qr", nil, true),
		map[string]any{"fids": []string{opts.UserID}})
	if err != nil {
		return nil, err
	}
	msgInfo := map[string]any{"contactUid": opts.UserID}
	if u, ok := qr[opts.UserID]; ok {
		msgInfo["qrCodeUrl"] = u
	}
	if opts.PhoneNumber != "" {
		msgInfo["phone"] = opts.PhoneNumber
	}
	params := map[string]any{
		"ttl": opts.TTL, "msgType": 6, "clientId": strconv.FormatInt(nowMs(), 10), "msgInfo": mustJSON(msgInfo),
	}
	path := "/api/message/forward"
	if threadType == ThreadTypeGroup {
		params["visibility"], params["grid"] = 0, threadID
		path = "/api/group/forward"
	} else {
		params["toId"], params["imei"] = threadID, a.IMEI
	}
	res, err := call[*SendCardResponse](ctx, a.Session, http.MethodPost, a.MakeURL(a.svc("file")+path, nil, true), params)
	if err == nil && res == nil {
		res = &SendCardResponse{}
	}
	return res, err
}
