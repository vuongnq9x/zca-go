package zca

import (
	"context"
	"strconv"
)

type SendVoiceOptions struct {
	VoiceURL string
	TTL      int64 // milliseconds
}

type SendVoiceResponse = SendVideoResponse // {msgId}

// SendVoice sends a voice message by URL to a user or group.
func (a *API) SendVoice(ctx context.Context, opts SendVoiceOptions, threadID string, threadType ThreadType) (*SendVoiceResponse, error) {
	if threadType != ThreadTypeUser && threadType != ThreadTypeGroup {
		return nil, newError("Thread type is invalid")
	}
	fileSize, err := a.mediaContentLength(ctx, opts.VoiceURL)
	if err != nil {
		return nil, newError("Unable to get voice content: " + err.Error())
	}
	params := map[string]any{
		"ttl":      opts.TTL,
		"zsource":  -1,
		"msgType":  3,
		"clientId": strconv.FormatInt(nowMs(), 10),
		"msgInfo":  mustJSON(map[string]any{"voiceUrl": opts.VoiceURL, "m4aUrl": opts.VoiceURL, "fileSize": fileSize}),
		"imei":     a.IMEI,
	}
	if threadType == ThreadTypeUser {
		params["toId"] = threadID
	} else {
		params["grid"], params["visibility"] = threadID, 0
	}
	return a.mediaForward(ctx, params, threadType)
}
