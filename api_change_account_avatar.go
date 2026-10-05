package zca

import (
	"cmp"
	"context"
	"time"
)

// ChangeAccountAvatar changes the logged-in account's avatar. Path sources need
// Options.ImageMetadataGetter; in-memory sources use src.Metadata.
func (a *API) ChangeAccountAvatar(ctx context.Context, src AttachmentSource) (string, error) {
	meta, err := a.mediaImageMetadata(src)
	if err != nil {
		return "", err
	}
	w, h := cmp.Or(meta.Width, 1080), cmp.Or(meta.Height, 1080)
	params := map[string]any{
		"avatarSize": 120,
		"clientId":   a.UID + time.Now().Format("15:04 02/01/2006"),
		"language":   a.Language,
		"metaData": mustJSON(map[string]any{
			"origin":    map[string]any{"width": w, "height": h},
			"processed": map[string]any{"width": w, "height": h, "size": meta.TotalSize},
		}),
	}
	return a.mediaUploadAvatar(ctx, src, a.svc("file")+"/api/profile/upavatar", params)
}

// mediaUploadAvatar posts the image as fileContent (blob, image/jpeg) with encrypted params.
func (a *API) mediaUploadAvatar(ctx context.Context, src AttachmentSource, serviceURL string, params map[string]any) (string, error) {
	buf, err := mediaSourceBytes(src)
	if err != nil {
		return "", err
	}
	body, header := mediaMultipart("fileContent", "blob", "image/jpeg", buf)
	enc, err := a.EncodeAES(mustJSON(params))
	if err != nil {
		return "", newError("Failed to encrypt params")
	}
	return mediaPost[string](ctx, a.Session, a.MakeURL(serviceURL, map[string]any{"params": enc}, true), body, header)
}
