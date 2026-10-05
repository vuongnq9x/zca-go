package zca

import (
	"context"
	"net/http"
	"strconv"
)

type SendVideoOptions struct {
	Msg          string // optional caption
	VideoURL     string
	ThumbnailURL string
	Duration     int64 // milliseconds
	Width        int64 // 0 means 1280
	Height       int64 // 0 means 720
	TTL          int64 // milliseconds
}

type SendVideoResponse struct {
	MsgID StringOrNumber `json:"msgId"`
}

// mediaContentLength HEADs u without Zalo headers/cookies (TS raw request) and returns its
// Content-Length, 0 if the response is not OK or has none.
func (a *API) mediaContentLength(ctx context.Context, u string) (int64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, u, nil)
	if err != nil {
		return 0, err
	}
	c := *a.client
	c.CheckRedirect = nil // fetch follows redirects
	resp, err := c.Do(req)
	if err != nil {
		return 0, err
	}
	resp.Body.Close()
	if !isHTTPOK(resp) {
		return 0, nil
	}
	return max(resp.ContentLength, 0), nil
}

// mediaForward posts encrypted params to file/api/{message,group}/forward.
func (a *API) mediaForward(ctx context.Context, params map[string]any, threadType ThreadType) (*SendVideoResponse, error) {
	path := "/api/message/forward"
	if threadType == ThreadTypeGroup {
		path = "/api/group/forward"
	}
	res, err := call[*SendVideoResponse](ctx, a.Session, http.MethodPost, a.MakeURL(a.svc("file")+path, nil, true), params)
	if err == nil && res == nil {
		res = &SendVideoResponse{}
	}
	return res, err
}

// SendVideo sends a video by URL to a user or group.
func (a *API) SendVideo(ctx context.Context, opts SendVideoOptions, threadID string, threadType ThreadType) (*SendVideoResponse, error) {
	if threadType != ThreadTypeUser && threadType != ThreadTypeGroup {
		return nil, newError("Thread type is invalid")
	}
	fileSize, err := a.mediaContentLength(ctx, opts.VideoURL)
	if err != nil {
		return nil, newError("Unable to get video content: " + err.Error())
	}
	width, height := opts.Width, opts.Height
	if width == 0 {
		width = 1280
	}
	if height == 0 {
		height = 720
	}
	params := map[string]any{
		"clientId": strconv.FormatInt(nowMs(), 10),
		"ttl":      opts.TTL,
		"zsource":  704,
		"msgType":  5,
		"msgInfo": mustJSON(map[string]any{
			"videoUrl": opts.VideoURL,
			"thumbUrl": opts.ThumbnailURL,
			"duration": opts.Duration,
			"width":    width,
			"height":   height,
			"fileSize": fileSize,
			"properties": map[string]any{
				"color": -1, "size": -1, "type": 1003, "subType": 0,
				"ext": map[string]any{"sSrcType": -1, "sSrcStr": "", "msg_warning_type": 0},
			},
			"title": opts.Msg,
		}),
		"imei": a.IMEI,
	}
	if threadType == ThreadTypeUser {
		params["toId"] = threadID
	} else {
		params["grid"], params["visibility"] = threadID, 0
	}
	return a.mediaForward(ctx, params, threadType)
}
