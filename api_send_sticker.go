package zca

import (
	"context"
	"net/http"
)

type SendStickerPayload struct {
	ID     int64 `json:"id"`
	CateID int64 `json:"cateId"`
	Type   int64 `json:"type"`
}

type SendStickerResponse struct {
	MsgID StringOrNumber `json:"msgId"`
}

// SendSticker sends a sticker to a thread.
func (a *API) SendSticker(ctx context.Context, sticker SendStickerPayload, threadID string, threadType ThreadType) (*SendStickerResponse, error) {
	if threadID == "" {
		return nil, newError("Missing threadId")
	}
	if sticker.ID == 0 {
		return nil, newError("Missing sticker id")
	}
	// cateId 0 is accepted, like TS (it only rejects undefined/null).
	if sticker.Type == 0 {
		return nil, newError("Missing sticker type")
	}
	params := map[string]any{
		"stickerId": sticker.ID,
		"cateId":    sticker.CateID,
		"type":      sticker.Type,
		"clientId":  nowMs(),
		"imei":      a.IMEI,
		"zsource":   101,
	}
	u := a.svc("chat") + "/api/message/sticker"
	if threadType == ThreadTypeGroup {
		params["grid"] = threadID
		u = a.svc("group") + "/api/group/sticker"
	} else {
		params["toid"] = threadID
	}
	res, err := call[*SendStickerResponse](ctx, a.Session, http.MethodPost, a.MakeURL(u, map[string]any{"nretry": "0"}, true), params)
	if err == nil && res == nil {
		res = &SendStickerResponse{}
	}
	return res, err
}
