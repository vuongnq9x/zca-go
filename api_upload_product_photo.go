package zca

import (
	"context"
	"fmt"
)

type UploadProductPhotoPayload struct {
	File AttachmentSource
}

type UploadProductPhotoResponse struct {
	NormalURL    string         `json:"normalUrl"`
	PhotoID      StringOrNumber `json:"photoId"`
	Finished     int64          `json:"finished"`
	HdURL        string         `json:"hdUrl"`
	ThumbURL     string         `json:"thumbUrl"`
	ClientFileID int64          `json:"clientFileId"`
	ChunkID      int64          `json:"chunkId"`
}

// UploadProductPhoto uploads a photo for quick messages, product catalogs or custom storage.
// Path sources need Options.ImageMetadataGetter.
func (a *API) UploadProductPhoto(ctx context.Context, payload UploadProductPhotoPayload) (*UploadProductPhotoResponse, error) {
	src := payload.File
	meta, err := a.mediaImageMetadata(src)
	if err != nil {
		return nil, err
	}
	buf, err := mediaSourceBytes(src)
	if err != nil {
		return nil, err
	}
	body, header := mediaMultipart("chunkContent", "undefined", "application/octet-stream", buf)
	now := nowMs()
	params := map[string]any{
		"totalChunk": 1,
		"fileName":   fmt.Sprintf("Base64_Img_Picker_%d.jpg", now),
		"clientId":   now,
		"totalSize":  meta.TotalSize,
		"imei":       a.IMEI,
		"chunkId":    1,
		"toid":       a.LoginInfo["send2me_id"],
		"featureId":  1,
	}
	enc, err := a.EncodeAES(mustJSON(params))
	if err != nil {
		return nil, newError("Failed to encrypt params")
	}
	res, err := mediaPost[*UploadProductPhotoResponse](ctx, a.Session,
		a.MakeURL(a.svc("file")+"/api/product/upload/photo", map[string]any{"params": enc}, true), body, header)
	if err == nil && res == nil {
		res = &UploadProductPhotoResponse{}
	}
	return res, err
}
