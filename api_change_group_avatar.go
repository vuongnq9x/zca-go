package zca

import (
	"cmp"
	"context"
	"time"
)

// ChangeGroupAvatar changes a group's avatar. Path sources need Options.ImageMetadataGetter.
func (a *API) ChangeGroupAvatar(ctx context.Context, src AttachmentSource, groupID string) (string, error) {
	meta, err := a.mediaImageMetadata(src)
	if err != nil {
		return "", err
	}
	params := map[string]any{
		"grid":         groupID,
		"avatarSize":   120,
		"clientId":     "g" + groupID + time.Now().Format("15:04 02/01/2006"),
		"imei":         a.IMEI,
		"originWidth":  cmp.Or(meta.Width, 1080),
		"originHeight": cmp.Or(meta.Height, 1080),
	}
	return a.mediaUploadAvatar(ctx, src, a.svc("file")+"/api/group/upavatar", params)
}
