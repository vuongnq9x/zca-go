package zca

import (
	"context"
	"net/http"
)

type AddQuickMessagePayload struct {
	Keyword string
	Title   string
	Media   *AttachmentSource
}

type AddQuickMessageResponse struct {
	Item    QuickMessage `json:"item"`
	Version int64        `json:"version"`
}

// AddQuickMessage adds a quick message. Zalo may return code 821 when the limit is reached.
func (a *API) AddQuickMessage(ctx context.Context, payload AddQuickMessagePayload) (*AddQuickMessageResponse, error) {
	typ := 0
	if payload.Media != nil {
		typ = 1
	}
	params := map[string]any{
		"keyword": payload.Keyword,
		"message": map[string]any{"title": payload.Title, "params": ""},
		"type":    typ,
		"imei":    a.IMEI,
	}
	if payload.Media != nil {
		up, err := a.UploadProductPhoto(ctx, UploadProductPhotoPayload{File: *payload.Media})
		if err != nil {
			return nil, err
		}
		normal, hd := up.NormalURL, up.HdURL
		if normal == "" {
			normal = up.HdURL
		}
		if hd == "" {
			hd = up.NormalURL
		}
		params["media"] = map[string]any{"items": []map[string]any{{
			"type":         0,
			"photoId":      up.PhotoID,
			"title":        "",
			"width":        "",
			"height":       "",
			"previewThumb": up.ThumbURL,
			"rawUrl":       normal,
			"thumbUrl":     up.ThumbURL,
			"normalUrl":    normal,
			"hdUrl":        hd,
		}}}
	}
	return call[*AddQuickMessageResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("quick_message")+"/api/quickmessage/create", nil, true), params)
}
