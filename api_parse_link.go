package zca

import (
	"context"
	"net/http"
)

type ParseLinkErrorMaps map[string]int64

type ParseLinkResponse struct {
	Data struct {
		Thumb string `json:"thumb"`
		Title string `json:"title"`
		Desc  string `json:"desc"`
		Src   string `json:"src"`
		Href  string `json:"href"`
		Media struct {
			Type       int64  `json:"type"`
			Count      int64  `json:"count"`
			MediaTitle string `json:"mediaTitle"`
			Artist     string `json:"artist"`
			StreamURL  string `json:"streamUrl"`
			StreamIcon string `json:"stream_icon"`
		} `json:"media"`
		StreamIcon string `json:"stream_icon"`
	} `json:"data"`
	ErrorMaps ParseLinkErrorMaps `json:"error_maps"`
}

// ParseLink parses a link preview.
func (a *API) ParseLink(ctx context.Context, link string) (*ParseLinkResponse, error) {
	params := map[string]any{"link": link, "version": 1, "imei": a.IMEI}
	return call[*ParseLinkResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("file")+"/api/message/parselink", nil, true), params)
}
