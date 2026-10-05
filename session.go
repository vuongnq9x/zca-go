package zca

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// ImageMetadata is returned by Options.ImageMetadataGetter.
type ImageMetadata struct {
	Width, Height int64
	Size          int64
}

type Options struct {
	SelfListen bool
	APIType    int // default 30
	APIVersion int // default 685
	// HTTPClient lets callers set a proxy/transport. Its Jar and CheckRedirect are ignored.
	HTTPClient *http.Client
	// Logger defaults to slog.Default(); use slog.New(slog.DiscardHandler) to silence.
	Logger              *slog.Logger
	ImageMetadataGetter func(filePath string) (*ImageMetadata, error)
}

// ServiceMap is zpw_service_map_v3: service name -> base URLs.
type ServiceMap map[string][]string

type ShareFileSettings struct {
	BigFileDomainList     []string `json:"big_file_domain_list"`
	MaxSizeShareFileV2    int64    `json:"max_size_share_file_v2"`
	MaxSizeShareFileV3    int64    `json:"max_size_share_file_v3"`
	FileUploadShowIcon1GB bool     `json:"file_upload_show_icon_1GB"`
	RestrictedExt         string   `json:"restricted_ext"`
	NextFileTime          int64    `json:"next_file_time"`
	MaxFile               int      `json:"max_file"`
	MaxSizePhoto          int64    `json:"max_size_photo"`
	MaxSizeShareFile      int64    `json:"max_size_share_file"`
	MaxSizeResizePhoto    int64    `json:"max_size_resize_photo"`
	MaxSizeGif            int64    `json:"max_size_gif"`
	MaxSizeOriginalPhoto  int64    `json:"max_size_original_photo"`
	ChunkSizeFile         int64    `json:"chunk_size_file"`
	RestrictedExtFile     []string `json:"restricted_ext_file"`
}

type SocketSettings struct {
	RotateErrorCodes []int `json:"rotate_error_codes"`
	Retries          map[string]struct {
		Max   int `json:"max"`
		Times any `json:"times"` // number or []number (ms)
	} `json:"retries"`
	PingInterval       int64 `json:"ping_interval"`
	ResetEndpoint      int64 `json:"reset_endpoint"`
	CloseAndRetryCodes []int `json:"close_and_retry_codes"`
	MaxMsgSize         int64 `json:"max_msg_size"`
	EnableCtrlSocket   bool  `json:"enable_ctrl_socket"`
	EnableChatSocket   bool  `json:"enable_chat_socket"`
}

type Settings struct {
	Features struct {
		Sharefile ShareFileSettings `json:"sharefile"`
		Socket    SocketSettings    `json:"socket"`
	} `json:"features"`
	Keepalive struct {
		AlwayKeepalive    int64 `json:"alway_keepalive"`
		KeepaliveDuration int64 `json:"keepalive_duration"`
		TimeDeactive      int64 `json:"time_deactive"`
	} `json:"keepalive"`
	Raw map[string]any `json:"-"`
}

// UploadEventData is delivered by the listener when an async file upload finishes.
type UploadEventData struct {
	FileURL string `json:"fileUrl"`
	FileID  string `json:"fileId"`
}

// Session is the logged-in state (zca-js ContextSession). Exposed for custom API calls.
type Session struct {
	UID, IMEI, UserAgent, Language, SecretKey string
	APIType, APIVersion                       int
	ServiceMap                                ServiceMap
	Settings                                  Settings
	LoginInfo                                 map[string]any
	ExtraVer                                  map[string]any
	Options                                   Options

	jar    *cookieJar
	client *http.Client

	uploadMu      sync.Mutex
	uploadWaiters map[string]chan UploadEventData
}

func newSession(opts Options) *Session {
	if opts.APIType == 0 {
		opts.APIType = 30
	}
	if opts.APIVersion == 0 {
		opts.APIVersion = 685
	}
	client := &http.Client{}
	if opts.HTTPClient != nil {
		c := *opts.HTTPClient
		client = &c
	}
	client.Jar = nil
	// Redirects are followed manually so every hop's Set-Cookie lands in our jar.
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return &Session{
		APIType: opts.APIType, APIVersion: opts.APIVersion, Options: opts,
		jar: &cookieJar{}, client: client, uploadWaiters: map[string]chan UploadEventData{},
	}
}

func (s *Session) log() *slog.Logger {
	if s.Options.Logger != nil {
		return s.Options.Logger
	}
	return slog.Default()
}

// Cookies exports the current cookie jar (zca-js getCookie).
func (s *Session) Cookies() []Cookie { return s.jar.export() }

// CookieHeader returns the Cookie header value for a URL.
func (s *Session) CookieHeader(rawURL string) string { return s.jar.header(rawURL) }

// MakeURL appends params and, if apiVersion, zpw_ver/zpw_type when missing.
func (s *Session) MakeURL(base string, params map[string]any, apiVersion bool) string {
	u, err := url.Parse(base)
	if err != nil {
		return base
	}
	q := u.Query()
	for k, v := range params {
		q.Add(k, fmt.Sprint(v))
	}
	if apiVersion {
		if !q.Has("zpw_ver") {
			q.Set("zpw_ver", fmt.Sprint(s.APIVersion))
		}
		if !q.Has("zpw_type") {
			q.Set("zpw_type", fmt.Sprint(s.APIType))
		}
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// EncodeAES encrypts data with the session secret key.
func (s *Session) EncodeAES(data string) (string, error) { return encodeAES(s.SecretKey, data) }

// Request sends an HTTP request with Zalo default headers + cookies, storing Set-Cookie
// and following redirects manually (zca-js utils.request).
func (s *Session) Request(ctx context.Context, method, rawURL string, body []byte, header http.Header) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, rawURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if body == nil {
		req.Body, req.ContentLength = nil, 0
	}
	if s.UserAgent == "" {
		return nil, newError("user agent is not available")
	}
	h := req.Header
	h.Set("Accept", "application/json, text/plain, */*")
	h.Set("Accept-Language", "en-US,en;q=0.9")
	h.Set("Content-Type", "application/x-www-form-urlencoded")
	h.Set("Origin", "https://chat.zalo.me")
	h.Set("Referer", "https://chat.zalo.me/")
	h.Set("User-Agent", s.UserAgent)
	for k, v := range header {
		h[k] = v
	}
	if c := s.jar.header(rawURL); c != "" {
		h.Set("Cookie", c)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	for _, c := range resp.Cookies() {
		s.jar.set(c, req.URL.Hostname())
	}
	if loc := resp.Header.Get("Location"); loc != "" {
		resp.Body.Close()
		next, err := req.URL.Parse(loc)
		if err != nil {
			return nil, err
		}
		h2 := header.Clone()
		if h2 == nil {
			h2 = http.Header{}
		}
		h2.Set("Referer", "https://id.zalo.me/")
		return s.Request(ctx, http.MethodGet, next.String(), nil, h2)
	}
	return resp, nil
}

type zaloEnvelope struct {
	ErrorCode    int             `json:"error_code"`
	ErrorMessage string          `json:"error_message"`
	Data         json.RawMessage `json:"data"`
}

// Resolve reads a Zalo response and returns its (decrypted) data field (zca-js resolveResponse).
func (s *Session) Resolve(resp *http.Response, encrypted bool) (json.RawMessage, error) {
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, newError(fmt.Sprintf("Request failed with status code %d", resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var env zaloEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, newError("Failed to parse response data")
	}
	if env.ErrorCode != 0 {
		return nil, &ZaloAPIError{Message: env.ErrorMessage, Code: env.ErrorCode}
	}
	if !encrypted {
		return env.Data, nil
	}
	var enc string
	if err := json.Unmarshal(env.Data, &enc); err != nil {
		return nil, newError("Failed to parse response data")
	}
	plain, err := decodeAES(s.SecretKey, enc)
	if err != nil {
		return nil, newError("Failed to parse response data")
	}
	var inner zaloEnvelope
	if err := json.Unmarshal(plain, &inner); err != nil {
		return nil, newError("Failed to parse response data")
	}
	if inner.ErrorCode != 0 {
		return nil, &ZaloAPIError{Message: inner.ErrorMessage, Code: inner.ErrorCode}
	}
	return inner.Data, nil
}

// encryptedForm encrypts params and returns the url-encoded body "params=<enc>".
func (s *Session) encryptedForm(params any) ([]byte, error) {
	enc, err := s.EncodeAES(mustJSON(params))
	if err != nil {
		return nil, newError("Failed to encrypt params")
	}
	return []byte(url.Values{"params": {enc}}.Encode()), nil
}

// call is the common shape of most zca-js APIs: encrypt params, send them as the
// "params" query (GET) or form field (POST), resolve the encrypted response into T.
func call[T any](ctx context.Context, s *Session, method, serviceURL string, params any) (T, error) {
	var zero T
	enc, err := s.EncodeAES(mustJSON(params))
	if err != nil {
		return zero, newError("Failed to encrypt params")
	}
	var body []byte
	if method == http.MethodGet {
		serviceURL = s.MakeURL(serviceURL, map[string]any{"params": enc}, true)
	} else {
		body = []byte(url.Values{"params": {enc}}.Encode())
	}
	resp, err := s.Request(ctx, method, serviceURL, body, nil)
	if err != nil {
		return zero, err
	}
	raw, err := s.Resolve(resp, true)
	if err != nil {
		return zero, err
	}
	return decodeData[T](raw)
}

// decodeData unmarshals raw into T; null/empty data yields T's zero value.
func decodeData[T any](raw json.RawMessage) (T, error) {
	var v T
	if len(raw) == 0 || string(raw) == "null" {
		return v, nil
	}
	switch p := any(&v).(type) {
	case *json.RawMessage:
		*p = raw
		return v, nil
	case *string:
		// TS types many responses as "" but Zalo sometimes sends a number/object; keep it as text.
		if raw[0] != '"' {
			*p = string(raw)
			return v, nil
		}
	}
	if err := json.Unmarshal(raw, &v); err != nil {
		return v, fmt.Errorf("zca: decode response: %w", err)
	}
	return v, nil
}

// svc returns the first base URL of a service ("" if missing).
func (s *Session) svc(name string) string {
	if u := s.ServiceMap[name]; len(u) > 0 {
		return u[0]
	}
	return ""
}

func (s *Session) waitUpload(fileID string) chan UploadEventData {
	ch := make(chan UploadEventData, 1)
	s.uploadMu.Lock()
	s.uploadWaiters[fileID] = ch
	s.uploadMu.Unlock()
	return ch
}

func (s *Session) cancelUpload(fileID string) {
	s.uploadMu.Lock()
	delete(s.uploadWaiters, fileID)
	s.uploadMu.Unlock()
}

func (s *Session) notifyUpload(d UploadEventData) {
	s.uploadMu.Lock()
	ch := s.uploadWaiters[d.FileID]
	delete(s.uploadWaiters, d.FileID)
	s.uploadMu.Unlock()
	if ch != nil {
		ch <- d
	}
}

func isHTTPOK(resp *http.Response) bool { return resp.StatusCode >= 200 && resp.StatusCode <= 299 }

func trimDataURL(img string) string { return strings.TrimPrefix(img, "data:image/png;base64,") }
