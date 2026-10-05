package zca

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

// UploadAttachmentType is TS UploadAttachmentImageResponse | UploadAttachmentVideoResponse |
// UploadAttachmentFileResponse, discriminated by FileType ("image" | "video" | "others").
type UploadAttachmentType struct {
	FileType     string `json:"fileType"`
	Finished     any    `json:"finished"` // number | boolean
	ClientFileID int64  `json:"clientFileId"`
	ChunkID      int64  `json:"chunkId"`
	TotalSize    int64  `json:"totalSize"`

	// image
	NormalURL string `json:"normalUrl,omitempty"`
	PhotoID   string `json:"photoId,omitempty"`
	HdURL     string `json:"hdUrl,omitempty"`
	ThumbURL  string `json:"thumbUrl,omitempty"`
	Width     int64  `json:"width,omitempty"`
	Height    int64  `json:"height,omitempty"`
	HdSize    int64  `json:"hdSize,omitempty"`

	// video | others
	FileURL  string `json:"fileUrl,omitempty"`
	FileID   string `json:"fileId,omitempty"`
	Checksum string `json:"checksum,omitempty"`
	FileName string `json:"fileName,omitempty"`
}

type uploadAttachmentRaw struct {
	Finished     any            `json:"finished"`
	ClientFileID int64          `json:"clientFileId"`
	ChunkID      int64          `json:"chunkId"`
	FileID       StringOrNumber `json:"fileId"`
	PhotoID      StringOrNumber `json:"photoId"`
	NormalURL    string         `json:"normalUrl"`
	HdURL        string         `json:"hdUrl"`
	ThumbURL     string         `json:"thumbUrl"`
}

// uploadAttachmentWaitTTL mirrors the zca-js upload callback TTL.
const uploadAttachmentWaitTTL = 5 * time.Minute

// UploadAttachment uploads files to a thread so they can be sent with SendMessage.
//
// Video (mp4) and other files finish asynchronously: Zalo announces them with a websocket
// "file_done" event, so the Listener must be started (a.Listener.Start) before uploading them,
// otherwise this waits up to 5 minutes and fails. Images need Options.ImageMetadataGetter when
// given by Path; in-memory sources must set Filename and Metadata.
func (a *API) UploadAttachment(ctx context.Context, sources []AttachmentSource, threadID string, threadType ThreadType) ([]UploadAttachmentType, error) {
	sf := a.Settings.Features.Sharefile
	if len(sources) == 0 {
		return nil, newError("Missing sources")
	}
	if len(sources) > sf.MaxFile {
		return nil, newError(fmt.Sprintf("Exceed maximum file of %d", sf.MaxFile))
	}
	if threadID == "" {
		return nil, newError("Missing threadId")
	}
	chunkSize := sf.ChunkSizeFile
	if chunkSize <= 0 {
		return nil, newError("Invalid chunk_size_file setting")
	}
	isGroup := threadType == ThreadTypeGroup
	typeParam, urlBase := "2", a.svc("file")+"/api/message/"
	if isGroup {
		typeParam, urlBase = "11", a.svc("file")+"/api/group/"
	}

	type attachment struct {
		fileType string
		fileName string
		meta     AttachmentMetadata
		params   map[string]any
		buf      []byte
	}
	var atts []attachment
	clientID := nowMs()
	for _, src := range sources {
		isPath := src.Path != ""
		if !isPath && src.Data == nil {
			return nil, newError("Invalid source type")
		}
		if !isPath && src.Filename == "" {
			return nil, newError("Missing filename")
		}
		if isPath {
			if _, err := os.Stat(src.Path); err != nil {
				return nil, newError("File not found")
			}
		}
		ext := strings.ToLower(mediaFileExt(mediaSourceName(src)))
		fileName := src.Filename
		if isPath {
			fileName = filepath.Base(src.Path)
		}
		if slices.Contains(sf.RestrictedExtFile, ext) {
			return nil, newError(fmt.Sprintf("File extension %q is not allowed", ext))
		}

		at := attachment{fileName: fileName}
		switch ext {
		case "jpg", "jpeg", "png", "webp":
			meta, err := a.mediaImageMetadata(src)
			if err != nil {
				return nil, err
			}
			at.fileType, at.meta = "image", meta
		default:
			at.fileType = "others"
			if ext == "mp4" {
				at.fileType = "video"
			}
			at.meta.TotalSize = src.Metadata.TotalSize
			if isPath {
				st, err := os.Stat(src.Path)
				if err != nil {
					return nil, err
				}
				at.meta.TotalSize = st.Size()
			}
		}
		if at.meta.TotalSize > sf.MaxSizeShareFileV3*1024*1024 {
			return nil, newError(fmt.Sprintf("File %s size exceed maximum size of %dMB", fileName, sf.MaxSizeShareFileV3))
		}
		at.params = map[string]any{
			"totalChunk": (at.meta.TotalSize + chunkSize - 1) / chunkSize,
			"fileName":   fileName,
			"clientId":   clientID,
			"totalSize":  at.meta.TotalSize,
			"imei":       a.IMEI,
			"isE2EE":     0,
			"jxl":        0,
			"chunkId":    1,
		}
		clientID++
		if isGroup {
			at.params["grid"] = threadID
		} else {
			at.params["toid"] = threadID
		}
		buf, err := mediaSourceBytes(src)
		if err != nil {
			return nil, err
		}
		at.buf = buf
		atts = append(atts, at)
	}

	results := make([]UploadAttachmentType, len(atts))
	for i, at := range atts {
		totalChunk := at.params["totalChunk"].(int64)
		urlType := "asyncfile/upload"
		if at.fileType == "image" {
			urlType = "photo_original/upload"
		}
		// ponytail: chunks go out sequentially (TS fires them in parallel); only the last chunk's
		// response carries a real fileId, so ordering is safe.
		for c := int64(0); c < totalChunk; c++ {
			at.params["chunkId"] = c + 1
			enc, err := a.EncodeAES(mustJSON(at.params))
			if err != nil {
				return nil, newError("Failed to encrypt message")
			}
			start, end := min(c*chunkSize, int64(len(at.buf))), min((c+1)*chunkSize, int64(len(at.buf)))
			body, header := mediaMultipart("chunkContent", at.fileName, "application/octet-stream", at.buf[start:end])
			u := a.MakeURL(urlBase+urlType, map[string]any{"type": typeParam, "params": enc}, true)
			raw, err := mediaPost[*uploadAttachmentRaw](ctx, a.Session, u, body, header)
			if err != nil {
				return nil, err
			}
			if raw == nil || raw.FileID == "-1" || raw.PhotoID == "-1" {
				continue
			}
			if at.fileType == "image" {
				results[i] = UploadAttachmentType{
					FileType: "image", Width: at.meta.Width, Height: at.meta.Height,
					TotalSize: at.meta.TotalSize, HdSize: at.meta.TotalSize,
					Finished: raw.Finished, NormalURL: raw.NormalURL, HdURL: raw.HdURL, ThumbURL: raw.ThumbURL,
					ChunkID: raw.ChunkID, PhotoID: string(raw.PhotoID), ClientFileID: raw.ClientFileID,
				}
				continue
			}
			if raw.FileID == "" { // TS throws on fileId!.toString() here
				return nil, newError("Missing fileId in upload response")
			}
			ws, err := a.uploadAttachmentWait(ctx, string(raw.FileID))
			if err != nil {
				return nil, err
			}
			results[i] = UploadAttachmentType{
				FileType: at.fileType, Finished: raw.Finished, ClientFileID: raw.ClientFileID, ChunkID: raw.ChunkID,
				PhotoID: string(raw.PhotoID), NormalURL: raw.NormalURL, HdURL: raw.HdURL, ThumbURL: raw.ThumbURL,
				FileURL: ws.FileURL, FileID: ws.FileID,
				TotalSize: at.meta.TotalSize, FileName: at.fileName,
				Checksum: mediaMD5(at.buf, at.meta.TotalSize),
			}
		}
	}
	return results, nil
}

func (a *API) uploadAttachmentWait(ctx context.Context, fileID string) (UploadEventData, error) {
	ch := a.Session.waitUpload(fileID)
	select {
	case d := <-ch:
		return d, nil
	case <-ctx.Done():
		a.Session.cancelUpload(fileID)
		return UploadEventData{}, ctx.Err()
	case <-time.After(uploadAttachmentWaitTTL):
		a.Session.cancelUpload(fileID)
		return UploadEventData{}, newError("Upload timed out waiting for file_done (is the Listener started?)")
	}
}

// mediaSourceName is the name used for extension checks: Path, else Filename.
func mediaSourceName(src AttachmentSource) string {
	if src.Path != "" {
		return src.Path
	}
	return src.Filename
}

func mediaSourceBytes(src AttachmentSource) ([]byte, error) {
	if src.Path != "" {
		return os.ReadFile(src.Path)
	}
	return src.Data, nil
}

// mediaFileExt is node path.extname(name).slice(1): dotfiles like ".bashrc" have no extension.
func mediaFileExt(name string) string {
	base := filepath.Base(name)
	ext := filepath.Ext(base)
	if ext == base {
		return ""
	}
	return strings.TrimPrefix(ext, ".")
}

// mediaImageMetadata is zca-js getImageMetaData for Path sources, else src.Metadata.
func (a *API) mediaImageMetadata(src AttachmentSource) (AttachmentMetadata, error) {
	if src.Path == "" {
		return src.Metadata, nil
	}
	if a.Options.ImageMetadataGetter == nil {
		return AttachmentMetadata{}, ErrMissingImageMetadataGetter
	}
	md, err := a.Options.ImageMetadataGetter(src.Path)
	if err != nil {
		return AttachmentMetadata{}, err
	}
	if md == nil {
		return AttachmentMetadata{}, newError("Failed to get image metadata")
	}
	return AttachmentMetadata{TotalSize: md.Size, Width: md.Width, Height: md.Height}, nil
}

// mediaMD5 is zca-js getMd5LargeFileObject: md5 hex of the first size bytes.
func mediaMD5(buf []byte, size int64) string {
	h := md5.Sum(buf[:max(0, min(size, int64(len(buf))))])
	return hex.EncodeToString(h[:])
}

var mediaQuoteEscaper = strings.NewReplacer(`\`, `\\`, `"`, `\"`)

// mediaMultipart builds a single-file multipart/form-data body like form-data's getBuffer/getHeaders.
func mediaMultipart(field, filename, contentType string, data []byte) ([]byte, http.Header) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	h := textproto.MIMEHeader{}
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, field, mediaQuoteEscaper.Replace(filename)))
	h.Set("Content-Type", contentType)
	p, _ := w.CreatePart(h) // writes to a bytes.Buffer cannot fail
	p.Write(data)
	w.Close()
	return b.Bytes(), http.Header{"Content-Type": {w.FormDataContentType()}}
}

// mediaPost POSTs body and decodes the encrypted Zalo response into T.
func mediaPost[T any](ctx context.Context, s *Session, u string, body []byte, header http.Header) (T, error) {
	var zero T
	resp, err := s.Request(ctx, http.MethodPost, u, body, header)
	if err != nil {
		return zero, err
	}
	raw, err := s.Resolve(resp, true)
	if err != nil {
		return zero, err
	}
	return decodeData[T](raw)
}
