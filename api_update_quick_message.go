package zca

import (
	"context"
	"net/http"
)

type UpdateQuickMessagePayload struct {
	Keyword string
	Title   string
	Media   *AttachmentSource // uploaded with UploadProductPhoto
}

type UpdateQuickMessageResponse struct {
	Item    QuickMessage `json:"item"`
	Version int64        `json:"version"`
}

// UpdateQuickMessage updates a quick message. Zalo error 212 means itemID does not exist.
func (a *API) UpdateQuickMessage(ctx context.Context, payload UpdateQuickMessagePayload, itemID int64) (*UpdateQuickMessageResponse, error) {
	params := map[string]any{
		"itemId":  itemID,
		"keyword": payload.Keyword,
		"message": map[string]any{"title": payload.Title, "params": ""},
		"type":    0,
	}
	if payload.Media != nil {
		params["type"] = 1
		up, err := a.UploadProductPhoto(ctx, UploadProductPhotoPayload{File: *payload.Media})
		if err != nil {
			return nil, err
		}
		or := func(x, y string) string {
			if x != "" {
				return x
			}
			return y
		}
		params["media"] = map[string]any{"items": []map[string]any{{
			"type": 0, "photoId": up.PhotoID, "title": "", "width": "", "height": "",
			"previewThumb": up.ThumbURL, "rawUrl": or(up.NormalURL, up.HdURL), "thumbUrl": up.ThumbURL,
			"normalUrl": or(up.NormalURL, up.HdURL), "hdUrl": or(up.HdURL, up.NormalURL),
		}}}
	}
	return call[*UpdateQuickMessageResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("quick_message")+"/api/quickmessage/update", nil, true), params)
}
