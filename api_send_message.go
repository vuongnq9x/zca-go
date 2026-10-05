package zca

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"unicode/utf16"
)

type SendMessageResult struct {
	MsgID StringOrNumber `json:"msgId"`
}

type SendMessageResponse struct {
	Message    *SendMessageResult  `json:"message"`
	Attachment []SendMessageResult `json:"attachment"`
}

// SendMessageQuote is the quoted message; build it with QuoteFromMessage.
type SendMessageQuote struct {
	Content     any                  `json:"content"` // string or object (map[string]any)
	MsgType     string               `json:"msgType"`
	PropertyExt *TMessagePropertyExt `json:"propertyExt"`
	UIDFrom     string               `json:"uidFrom"`
	MsgID       string               `json:"msgId"`
	CliMsgID    string               `json:"cliMsgId"`
	TS          string               `json:"ts"`
	TTL         int64                `json:"ttl"`
}

// QuoteFromMessage builds a SendMessageQuote from a received message.
func QuoteFromMessage(m *Message) *SendMessageQuote {
	d := m.Data
	return &SendMessageQuote{
		Content: d.Content, MsgType: d.MsgType, PropertyExt: d.PropertyExt,
		UIDFrom: d.UIDFrom, MsgID: d.MsgID, CliMsgID: d.CliMsgID, TS: d.TS, TTL: d.TTL,
	}
}

type TextStyle string

const (
	TextStyleBold          TextStyle = "b"
	TextStyleItalic        TextStyle = "i"
	TextStyleUnderline     TextStyle = "u"
	TextStyleStrikeThrough TextStyle = "s"
	TextStyleRed           TextStyle = "c_db342e"
	TextStyleOrange        TextStyle = "c_f27806"
	TextStyleYellow        TextStyle = "c_f7b503"
	TextStyleGreen         TextStyle = "c_15a85f"
	TextStyleSmall         TextStyle = "f_13"
	TextStyleBig           TextStyle = "f_18"
	TextStyleUnorderedList TextStyle = "lst_1"
	TextStyleOrderedList   TextStyle = "lst_2"
	TextStyleIndent        TextStyle = "ind_$"
)

type Style struct {
	Start int64     `json:"start"`
	Len   int64     `json:"len"`
	St    TextStyle `json:"st"`
	// IndentSize is the number of indent levels for TextStyleIndent (0 means 1).
	IndentSize int `json:"-"`
}

type Urgency int

const (
	UrgencyDefault Urgency = iota
	UrgencyImportant
	UrgencyUrgent
)

type Mention struct {
	Pos int64  // mention position (UTF-16 units, like JS)
	UID string // id of the mentioned user, "-1" mentions everyone
	Len int64  // length of the mention
}

type MessageContent struct {
	Msg         string
	Styles      []Style
	Urgency     Urgency
	Quote       *SendMessageQuote
	Mentions    []Mention
	Attachments []AttachmentSource
	TTL         int64 // milliseconds
}

// clientMessageType is zca-js getClientMessageType.
func clientMessageType(msgType string) int {
	switch msgType {
	case "webchat":
		return 1
	case "chat.voice":
		return 31
	case "chat.photo":
		return 32
	case "chat.sticker":
		return 36
	case "chat.doodle":
		return 37
	case "chat.recommended", "chat.link":
		return 38
	case "chat.video.msg":
		return 44
	case "share.file":
		return 46
	case "chat.gif":
		return 49
	case "chat.location.new":
		return 43
	}
	return 1
}

// sendMessageQuoteContent returns the quote content as a string or an object map.
func sendMessageQuoteContent(q *SendMessageQuote) (s string, obj map[string]any, isString bool) {
	if s, ok := q.Content.(string); ok {
		return s, nil, true
	}
	if m, ok := q.Content.(map[string]any); ok {
		return "", m, false
	}
	_ = remarshal(q.Content, &obj) // structs like TAttachmentContent
	return "", obj, false
}

// prepareQMSGAttach returns nil when TS would yield undefined (key dropped).
func prepareQMSGAttach(q *SendMessageQuote) any {
	_, obj, isString := sendMessageQuoteContent(q)
	if isString {
		if q.PropertyExt == nil {
			return nil
		}
		return q.PropertyExt
	}
	if q.MsgType == "chat.todo" {
		return map[string]any{"properties": map[string]any{
			"color": 0, "size": 0, "type": 0, "subType": 0, "ext": `{"shouldParseLinkOrContact":0}`,
		}}
	}
	out := make(map[string]any, len(obj)+3)
	for k, v := range obj {
		out[k] = v
	}
	if v, ok := obj["thumb"]; ok {
		out["thumbUrl"] = v
	}
	if v, ok := obj["href"]; ok {
		out["oriUrl"], out["normalUrl"] = v, v
	}
	return out
}

func prepareQMSG(q *SendMessageQuote) (any, error) {
	_, obj, isString := sendMessageQuoteContent(q)
	if q.MsgType == "chat.todo" && !isString && obj != nil {
		if p, ok := obj["params"].(string); ok {
			var v struct {
				Item struct {
					Content any `json:"content"`
				} `json:"item"`
			}
			if err := json.Unmarshal([]byte(p), &v); err != nil {
				return nil, err
			}
			return v.Item.Content, nil
		}
	}
	return "", nil
}

type sendMessageMention struct {
	Pos  int64  `json:"pos"`
	UID  string `json:"uid"`
	Len  int64  `json:"len"`
	Type int    `json:"type"`
}

func sendMessageMentions(threadType ThreadType, msg string, mentions []Mention) ([]sendMessageMention, error) {
	var out []sendMessageMention
	var total int64
	if threadType == ThreadTypeGroup {
		for _, m := range mentions {
			if m.Pos < 0 || m.UID == "" || m.Len <= 0 {
				continue
			}
			total += m.Len
			t := 0
			if m.UID == "-1" {
				t = 1
			}
			out = append(out, sendMessageMention{Pos: m.Pos, UID: m.UID, Len: m.Len, Type: t})
		}
	}
	if total > int64(len(utf16.Encode([]rune(msg)))) {
		return nil, newError("Invalid mentions: total mention characters exceed message length")
	}
	return out, nil
}

func sendMessageStyles(params map[string]any, styles []Style) {
	if styles == nil {
		return
	}
	out := make([]Style, len(styles))
	for i, s := range styles {
		if s.St == TextStyleIndent {
			s.St = TextStyle(strings.ReplaceAll(string(TextStyleIndent), "$", strconv.Itoa(max(s.IndentSize, 1))+"0"))
		}
		out[i] = s
	}
	params["textProperties"] = mustJSON(map[string]any{"styles": out, "ver": 0})
}

func sendMessageUrgency(params map[string]any, u Urgency) {
	if u == UrgencyImportant || u == UrgencyUrgent {
		params["metaData"] = map[string]any{"urgency": u}
	}
}

type sendMessageRequest struct {
	url    string
	body   []byte
	header http.Header
}

func mediaForm(enc string) []byte { return []byte(url.Values{"params": {enc}}.Encode()) }

func (a *API) sendMessageSend(ctx context.Context, reqs []sendMessageRequest) ([]SendMessageResult, error) {
	// ponytail: sequential, TS uses Promise.all; result order is the same.
	out := make([]SendMessageResult, 0, len(reqs))
	for _, r := range reqs {
		res, err := mediaPost[*SendMessageResult](ctx, a.Session, r.url, r.body, r.header)
		if err != nil {
			return nil, err
		}
		if res == nil {
			res = &SendMessageResult{}
		}
		out = append(out, *res)
	}
	return out, nil
}

func (a *API) sendMessageText(m MessageContent, threadID string, threadType ThreadType) (sendMessageRequest, error) {
	if m.Msg == "" {
		return sendMessageRequest{}, newError("Missing message content")
	}
	isGroup := threadType == ThreadTypeGroup
	mentions, err := sendMessageMentions(threadType, m.Msg, m.Mentions)
	if err != nil {
		return sendMessageRequest{}, err
	}
	q := m.Quote
	if q != nil {
		if _, _, isString := sendMessageQuoteContent(q); !isString && q.MsgType == "webchat" {
			return sendMessageRequest{}, newError("This kind of `webchat` quote type is not available")
		}
		if q.MsgType == "group.poll" {
			return sendMessageRequest{}, newError("The `group.poll` quote type is not available")
		}
	}

	params := map[string]any{"message": m.Msg, "clientId": nowMs(), "ttl": m.TTL}
	if len(mentions) > 0 && isGroup {
		params["mentionInfo"] = mustJSON(mentions)
	}
	if isGroup {
		params["grid"], params["visibility"] = threadID, 0
	} else {
		params["toid"], params["imei"] = threadID, a.IMEI
	}
	if q != nil {
		params["qmsgOwner"] = q.UIDFrom
		params["qmsgId"] = q.MsgID
		params["qmsgCliId"] = q.CliMsgID
		params["qmsgType"] = clientMessageType(q.MsgType)
		params["qmsgTs"] = q.TS
		params["qmsgTTL"] = q.TTL
		if s, _, isString := sendMessageQuoteContent(q); isString {
			params["qmsg"] = s
		} else if params["qmsg"], err = prepareQMSG(q); err != nil {
			return sendMessageRequest{}, err
		}
		if isGroup {
			if att := prepareQMSGAttach(q); att != nil {
				params["qmsgAttach"] = mustJSON(att)
			}
		}
	}
	sendMessageStyles(params, m.Styles)
	sendMessageUrgency(params, m.Urgency)

	enc, err := a.EncodeAES(mustJSON(params))
	if err != nil {
		return sendMessageRequest{}, newError("Failed to encrypt message")
	}
	base := a.svc("chat") + "/api/message"
	if isGroup {
		base = a.svc("group") + "/api/group"
	}
	switch {
	case q != nil:
		base += "/quote"
	case !isGroup:
		base += "/sms"
	case params["mentionInfo"] != nil:
		base += "/mention"
	default:
		base += "/sendmsg"
	}
	return sendMessageRequest{url: a.MakeURL(base, map[string]any{"nretry": 0}, true), body: mediaForm(enc)}, nil
}

type sendMessageUpthumb struct {
	HdURL        string         `json:"hdUrl"`
	ClientFileID StringOrNumber `json:"clientFileId"`
	URL          string         `json:"url"`
	FileID       StringOrNumber `json:"fileId"`
}

func (a *API) sendMessageUpthumb(ctx context.Context, buf []byte, baseURL string) (*sendMessageUpthumb, error) {
	body, header := mediaMultipart("fileContent", "blob", "image/png", buf)
	enc, err := a.EncodeAES(mustJSON(map[string]any{"clientId": nowMs(), "imei": a.IMEI}))
	if err != nil {
		return nil, newError("Failed to encrypt message")
	}
	res, err := mediaPost[*sendMessageUpthumb](ctx, a.Session, a.MakeURL(baseURL+"upthumb", map[string]any{"params": enc}, true), body, header)
	if err == nil && res == nil {
		res = &sendMessageUpthumb{}
	}
	return res, err
}

func (a *API) sendMessageAttachments(ctx context.Context, m MessageContent, threadID string, threadType ThreadType) ([]sendMessageRequest, error) {
	if len(m.Attachments) == 0 {
		return nil, newError("Missing attachments")
	}
	sf := a.Settings.Features.Sharefile
	isGroup := threadType == ThreadTypeGroup
	attURL := map[ThreadType]string{
		ThreadTypeUser:  a.svc("file") + "/api/message/",
		ThreadTypeGroup: a.svc("file") + "/api/group/",
	}
	firstExt := mediaFileExt(mediaSourceName(m.Attachments[0]))
	canBeDesc := len(m.Attachments) == 1 && slices.Contains([]string{"jpg", "jpeg", "png", "webp"}, firstExt)

	var gifs, files []AttachmentSource
	for _, s := range m.Attachments {
		if mediaFileExt(mediaSourceName(s)) == "gif" {
			gifs = append(gifs, s)
		} else {
			files = append(files, s)
		}
	}
	var uploaded []UploadAttachmentType
	if len(files) > 0 {
		var err error
		if uploaded, err = a.UploadAttachment(ctx, files, threadID, threadType); err != nil {
			return nil, err
		}
	}

	groupLayoutID := strconv.FormatInt(nowMs(), 10)
	mentions, err := sendMessageMentions(threadType, m.Msg, m.Mentions)
	if err != nil {
		return nil, err
	}
	mentionsValid := len(mentions) > 0 && isGroup && len(files) == 1
	isMultiFile := len(files) > 1
	clientID := nowMs()
	indexInGroup := 0
	var reqs []sendMessageRequest
	for _, at := range uploaded {
		var params map[string]any
		var urlType string
		switch at.FileType {
		case "image":
			urlType = "photo_original/send"
			params = map[string]any{
				"photoId": at.PhotoID, "clientId": strconv.FormatInt(clientID, 10), "desc": m.Msg,
				"width": at.Width, "height": at.Height,
				"rawUrl": at.NormalURL, "hdUrl": at.HdURL, "thumbUrl": at.ThumbURL,
				"hdSize": strconv.FormatInt(at.TotalSize, 10), "zsource": -1, "ttl": m.TTL,
				"jcp": `{"convertible":"jxl"}`,
			}
			clientID++
			if isGroup {
				params["grid"], params["oriUrl"] = threadID, at.NormalURL
			} else {
				params["toid"], params["normalUrl"] = threadID, at.NormalURL
			}
			if isMultiFile {
				params["groupLayoutId"] = groupLayoutID
				params["isGroupLayout"] = 1
				params["idInGroup"] = indexInGroup
				params["totalItemInGroup"] = len(uploaded)
				params["extMsgProp"] = `{"groupMediaMsg":{"groupLayoutId":"` + groupLayoutID + `"}}`
				indexInGroup++
			}
			if mentionsValid && canBeDesc && m.Quote == nil {
				params["mentionInfo"] = mustJSON(mentions)
			}
		case "video", "others":
			urlType = "asyncfile/msg"
			params = map[string]any{
				"fileId": at.FileID, "checksum": at.Checksum, "checksumSha": "",
				"extention": mediaFileExt(at.FileName), "totalSize": at.TotalSize, "fileName": at.FileName,
				"clientId": at.ClientFileID, "fType": 1, "fileCount": 0, "fdata": "{}",
				"fileUrl": at.FileURL, "zsource": -1, "ttl": m.TTL,
			}
			if isGroup {
				params["grid"] = threadID
			} else {
				params["toid"] = threadID
			}
		default:
			continue // ponytail: hole left by an upload that never reported (TS would throw)
		}
		sendMessageUrgency(params, m.Urgency)
		enc, err := a.EncodeAES(mustJSON(params))
		if err != nil {
			return nil, newError("Failed to encrypt message")
		}
		reqs = append(reqs, sendMessageRequest{
			url:  a.MakeURL(attURL[threadType]+urlType, map[string]any{"nretry": "0"}, true),
			body: mediaForm(enc),
		})
	}

	for _, gif := range gifs {
		meta, err := a.mediaImageMetadata(gif)
		if err != nil {
			return nil, err
		}
		fileName := gif.Filename
		if gif.Path != "" {
			fileName = filepath.Base(gif.Path)
		}
		if meta.TotalSize > sf.MaxSizeShareFileV3*1024*1024 {
			return nil, newError(fmt.Sprintf("File %s size exceed maximum size of %dMB", fileName, sf.MaxSizeShareFileV3))
		}
		buf, err := mediaSourceBytes(gif)
		if err != nil {
			return nil, err
		}
		thumb, err := a.sendMessageUpthumb(ctx, buf, attURL[ThreadTypeUser])
		if err != nil {
			return nil, err
		}
		body, header := mediaMultipart("chunkContent", fileName, "application/octet-stream", buf)
		params := map[string]any{
			"clientId": strconv.FormatInt(nowMs(), 10), "fileName": fileName, "totalSize": meta.TotalSize,
			"width": meta.Width, "height": meta.Height, "msg": m.Msg, "type": 1, "ttl": m.TTL,
			"thumb": thumb.URL, "checksum": mediaMD5(buf, meta.TotalSize), "totalChunk": 1, "chunkId": 1,
		}
		if isGroup {
			params["visibility"], params["grid"] = 0, threadID
		} else {
			params["toid"] = threadID
		}
		sendMessageUrgency(params, m.Urgency)
		enc, err := a.EncodeAES(mustJSON(params))
		if err != nil {
			return nil, newError("Failed to encrypt message")
		}
		reqs = append(reqs, sendMessageRequest{
			url:    a.MakeURL(attURL[threadType]+"gif", map[string]any{"nretry": "0", "params": enc, "type": "1"}, true),
			body:   body,
			header: header,
		})
	}
	return reqs, nil
}

// SendMessage sends text and/or attachments to a thread. File attachments other than images
// and gifs need a running Listener (see UploadAttachment).
func (a *API) SendMessage(ctx context.Context, msg MessageContent, threadID string, threadType ThreadType) (*SendMessageResponse, error) {
	if threadID == "" {
		return nil, newError("Missing threadId")
	}
	sf := a.Settings.Features.Sharefile
	text, mentions := msg.Msg, msg.Mentions
	if text == "" && len(msg.Attachments) == 0 {
		return nil, newError("Missing message content")
	}
	if len(msg.Attachments) > sf.MaxFile {
		return nil, newError(fmt.Sprintf("Exceed maximum file of %d", sf.MaxFile))
	}
	res := &SendMessageResponse{Attachment: []SendMessageResult{}}

	if len(msg.Attachments) > 0 {
		firstExt := mediaFileExt(mediaSourceName(msg.Attachments[0]))
		canBeDesc := len(msg.Attachments) == 1 && slices.Contains([]string{"jpg", "jpeg", "png", "webp"}, firstExt)
		if (!canBeDesc && len(text) > 0) || (len(text) > 0 && msg.Quote != nil) {
			// send message and attachment separately
			req, err := a.sendMessageText(msg, threadID, threadType)
			if err != nil {
				return nil, err
			}
			out, err := a.sendMessageSend(ctx, []sendMessageRequest{req})
			if err != nil {
				return nil, err
			}
			res.Message = &out[0]
			text, mentions = "", nil
		}
		att := msg
		att.Msg, att.Mentions = text, mentions
		reqs, err := a.sendMessageAttachments(ctx, att, threadID, threadType)
		if err != nil {
			return nil, err
		}
		if res.Attachment, err = a.sendMessageSend(ctx, reqs); err != nil {
			return nil, err
		}
		text = ""
	}

	if len(text) > 0 {
		req, err := a.sendMessageText(msg, threadID, threadType)
		if err != nil {
			return nil, err
		}
		out, err := a.sendMessageSend(ctx, []sendMessageRequest{req})
		if err != nil {
			return nil, err
		}
		res.Message = &out[0]
	}
	return res, nil
}
