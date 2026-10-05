package zca

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type SendLinkOptions struct {
	Msg  string // optional; the link is appended if missing
	Link string
	TTL  int64 // milliseconds
}

type SendLinkResponse struct {
	MsgID StringOrNumber `json:"msgId"`
}

type sendLinkParsed struct {
	Data struct {
		Thumb string          `json:"thumb"`
		Title string          `json:"title"`
		Desc  string          `json:"desc"`
		Src   string          `json:"src"`
		Href  string          `json:"href"`
		Media json.RawMessage `json:"media"`
	} `json:"data"`
}

// SendLink sends a link with its parsed preview.
func (a *API) SendLink(ctx context.Context, opts SendLinkOptions, threadID string, threadType ThreadType) (*SendLinkResponse, error) {
	// ponytail: inline parselink call; switch to a.ParseLink once that API lands.
	parsed, err := call[*sendLinkParsed](ctx, a.Session, http.MethodGet, a.MakeURL(a.svc("file")+"/api/message/parselink", nil, true),
		map[string]any{"link": opts.Link, "version": 1, "imei": a.IMEI})
	if err != nil {
		return nil, err
	}
	if parsed == nil {
		parsed = &sendLinkParsed{}
	}
	msg := opts.Link
	if strings.TrimSpace(opts.Msg) != "" {
		msg = opts.Msg
		if !strings.Contains(opts.Msg, opts.Link) {
			msg += " " + opts.Link
		}
	}
	d := parsed.Data
	params := map[string]any{
		"msg": msg, "href": d.Href, "src": d.Src, "title": d.Title, "desc": d.Desc, "thumb": d.Thumb,
		"type": 2, "ttl": opts.TTL, "clientId": nowMs(),
	}
	if len(d.Media) > 0 { // JSON.stringify(undefined) drops the key
		params["media"] = string(d.Media)
	}
	u := a.svc("chat") + "/api/message/link"
	if threadType == ThreadTypeGroup {
		params["grid"], params["imei"] = threadID, a.IMEI
		u = a.svc("group") + "/api/group/sendlink"
	} else {
		params["toId"], params["mentionInfo"] = threadID, ""
	}
	res, err := call[*SendLinkResponse](ctx, a.Session, http.MethodPost, a.MakeURL(u, map[string]any{"nretry": 0}, true), params)
	if err == nil && res == nil {
		res = &SendLinkResponse{}
	}
	return res, err
}
