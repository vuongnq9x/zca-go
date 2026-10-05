package zca

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"
)

func nowMs() int64 { return time.Now().UnixMilli() }

// login calls getLoginInfo with encrypted params and returns the decrypted login info.
func login(ctx context.Context, s *Session) (map[string]any, error) {
	enc, err := newParamsEncryptor(s.APIType, s.IMEI, nowMs())
	if err != nil {
		return nil, err
	}
	data := struct {
		ComputerName string `json:"computer_name"`
		IMEI         string `json:"imei"`
		Language     string `json:"language"`
		TS           int64  `json:"ts"`
	}{"Web", s.IMEI, s.Language, nowMs()}
	encoded, err := aesCBCEncrypt([]byte(enc.encryptKey), []byte(mustJSON(data)))
	if err != nil {
		return nil, newError("Failed to encrypt params: " + err.Error())
	}
	params := enc.params()
	params["params"] = base64.StdEncoding.EncodeToString(encoded)
	params["type"] = strconv.Itoa(s.APIType)
	params["client_version"] = strconv.Itoa(s.APIVersion)
	params["signkey"] = getSignKey("getlogininfo", params)

	q := map[string]any{"nretry": 0}
	for k, v := range params {
		q[k] = v
	}
	resp, err := s.Request(ctx, http.MethodGet, s.MakeURL("https://wpa.chat.zalo.me/api/login/getLoginInfo", q, true), nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if !isHTTPOK(resp) {
		return nil, newError("Failed to fetch login info: " + resp.Status)
	}
	var body struct {
		Data string `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	plain, err := decryptResp(enc.encryptKey, body.Data)
	if err != nil {
		return nil, newError("Đăng nhập thất bại")
	}
	var out struct {
		Data map[string]any `json:"data"`
	}
	dec := json.NewDecoder(bytes.NewReader(plain))
	dec.UseNumber() // keep uid and other big ids exact
	if err := dec.Decode(&out); err != nil || out.Data == nil {
		return nil, newError("Đăng nhập thất bại")
	}
	return out.Data, nil
}

type serverInfo struct {
	Settings  json.RawMessage `json:"settings"`
	Setttings json.RawMessage `json:"setttings"` // Zalo's typo, preferred like zca-js
	ExtraVer  map[string]any  `json:"extra_ver"`
}

func getServerInfo(ctx context.Context, s *Session) (*serverInfo, error) {
	sign := getSignKey("getserverinfo", map[string]string{
		"imei": s.IMEI, "type": strconv.Itoa(s.APIType), "client_version": strconv.Itoa(s.APIVersion), "computer_name": "Web",
	})
	u := s.MakeURL("https://wpa.chat.zalo.me/api/login/getServerInfo", map[string]any{
		"imei": s.IMEI, "type": s.APIType, "client_version": s.APIVersion, "computer_name": "Web", "signkey": sign,
	}, false)
	resp, err := s.Request(ctx, http.MethodGet, u, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if !isHTTPOK(resp) {
		return nil, newError("Failed to fetch server info: " + resp.Status)
	}
	var body struct {
		Data         *serverInfo `json:"data"`
		ErrorMessage string      `json:"error_message"`
	}
	dec := json.NewDecoder(resp.Body)
	dec.UseNumber()
	if err := dec.Decode(&body); err != nil {
		return nil, err
	}
	if body.Data == nil {
		return nil, newError("Failed to fetch server info: " + body.ErrorMessage)
	}
	return body.Data, nil
}

// ---- QR login ----

type LoginQREventType int

const (
	LoginQREventQRCodeGenerated LoginQREventType = iota
	LoginQREventQRCodeExpired
	LoginQREventQRCodeScanned
	LoginQREventQRCodeDeclined
	LoginQREventGotLoginInfo
)

// LoginQRAction is what the callback wants to do next. Ignored for Generated/Scanned/GotLoginInfo
// unless it is LoginQRAbort/LoginQRRetry.
type LoginQRAction int

const (
	LoginQRContinue LoginQRAction = iota
	LoginQRRetry
	LoginQRAbort
)

type QRCodeData struct {
	Code    string `json:"code"`
	Image   string `json:"image"` // base64 PNG without data URL prefix
	Options struct {
		EnabledCheckOCR   bool `json:"enabledCheckOCR"`
		EnabledMultiLayer bool `json:"enabledMultiLayer"`
	} `json:"options"`
	Token string `json:"token"`
}

// SaveToFile writes the QR PNG to path.
func (q *QRCodeData) SaveToFile(path string) error {
	img, err := base64.StdEncoding.DecodeString(q.Image)
	if err != nil {
		return err
	}
	return os.WriteFile(path, img, 0o644)
}

type QRScanData struct {
	Avatar      string `json:"avatar"`
	DisplayName string `json:"display_name"`
}

type LoginQREvent struct {
	Type      LoginQREventType
	QRCode    *QRCodeData  // Generated
	Scanned   *QRScanData  // Scanned
	Code      string       // Declined
	LoginInfo *Credentials // GotLoginInfo
}

type LoginQRCallback func(LoginQREvent) LoginQRAction

var errQRRetry = errors.New("qr retry")

var (
	reWindowsUA = regexp.MustCompile(`(?i)Windows`)
	reMacUA     = regexp.MustCompile(`(?i)Macintosh|Mac OS X`)
	reLinuxUA   = regexp.MustCompile(`(?i)Linux|X11`)
	reChromeUA  = regexp.MustCompile(`Chrome/(\d+)`)
)

// platformFromUA derives sec-ch-ua-platform from the User-Agent so client hints match it.
func platformFromUA(ua string) string {
	switch {
	case reWindowsUA.MatchString(ua):
		return "Windows"
	case reMacUA.MatchString(ua):
		return "macOS"
	case reLinuxUA.MatchString(ua):
		return "Linux"
	}
	return "Windows"
}

// chromeVersionFromUA returns the Chrome major version in ua, or "130" when absent.
func chromeVersionFromUA(ua string) string {
	if m := reChromeUA.FindStringSubmatch(ua); m != nil {
		return m[1]
	}
	return "130"
}

func qrHeaders(ua, referer string, document bool) http.Header {
	ver := chromeVersionFromUA(ua)
	h := http.Header{}
	h.Set("Accept", "*/*")
	h.Set("Accept-Language", "vi-VN,vi;q=0.9,fr-FR;q=0.8,fr;q=0.7,en-US;q=0.6,en;q=0.5")
	h.Set("sec-ch-ua", `"Chromium";v="`+ver+`", "Google Chrome";v="`+ver+`", "Not?A_Brand";v="99"`)
	h.Set("sec-ch-ua-mobile", "?0")
	h.Set("sec-ch-ua-platform", `"`+platformFromUA(ua)+`"`)
	h.Set("sec-fetch-dest", "empty")
	h.Set("sec-fetch-mode", "cors")
	h.Set("sec-fetch-site", "same-origin")
	h.Set("Priority", "u=1, i")
	h.Set("Referer", referer)
	h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
	if document {
		h.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		h.Set("sec-fetch-dest", "document")
		h.Set("sec-fetch-mode", "navigate")
		h.Set("upgrade-insecure-requests", "1")
		h.Set("Priority", "u=0, i")
		h.Del("Content-Type")
	}
	return h
}

const (
	refPC   = "https://id.zalo.me/account?continue=https%3A%2F%2Fzalo.me%2Fpc"
	refChat = "https://id.zalo.me/account?continue=https%3A%2F%2Fchat.zalo.me%2F"
)

func qrPost(ctx context.Context, s *Session, endpoint, referer string, form url.Values, out any) error {
	resp, err := s.Request(ctx, http.MethodPost, endpoint, []byte(form.Encode()), qrHeaders(s.UserAgent, referer, false))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(out)
}

type qrResult[T any] struct {
	Data         *T     `json:"data"`
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

// qrPoll repeats a long-poll request while Zalo answers error_code 8 (still waiting).
func qrPoll[T any](ctx context.Context, s *Session, endpoint string, form url.Values) (*qrResult[T], error) {
	for {
		var r qrResult[T]
		if err := qrPost(ctx, s, endpoint, refChat, form, &r); err != nil {
			return nil, err
		}
		if r.ErrorCode != 8 {
			return &r, nil
		}
	}
}

var loginVersionRe = regexp.MustCompile(`https://stc-zlogin\.zdn\.vn/main-([\d.]+)\.js`)

// loginQR runs one QR attempt; returns errQRRetry when the caller should start over.
func loginQROnce(ctx context.Context, s *Session, qrPath string, cb LoginQRCallback) ([]Cookie, error) {
	s.jar = &cookieJar{}

	h := qrHeaders(s.UserAgent, "https://chat.zalo.me/", true)
	h.Set("cache-control", "max-age=0")
	h.Set("sec-fetch-site", "same-site")
	h.Set("sec-fetch-user", "?1")
	resp, err := s.Request(ctx, http.MethodGet, "https://id.zalo.me/account?continue=https%3A%2F%2Fchat.zalo.me%2F", nil, h)
	if err != nil {
		return nil, err
	}
	html, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	m := loginVersionRe.FindSubmatch(html)
	if m == nil {
		return nil, newError("Cannot get API login version")
	}
	version := string(m[1])
	s.log().Info("Got login version", "version", version)

	var ignored any
	base := url.Values{"continue": {"https://zalo.me/pc"}, "v": {version}}
	if err := qrPost(ctx, s, "https://id.zalo.me/account/logininfo", refPC, base, &ignored); err != nil {
		s.log().Error("logininfo failed", "err", err)
	}
	verify := url.Values{"type": {"device"}, "continue": {"https://zalo.me/pc"}, "v": {version}}
	if err := qrPost(ctx, s, "https://id.zalo.me/account/verify-client", refPC, verify, &ignored); err != nil {
		s.log().Error("verify-client failed", "err", err)
	}
	var gen qrResult[QRCodeData]
	if err := qrPost(ctx, s, "https://id.zalo.me/account/authen/qr/generate", refPC, base, &gen); err != nil || gen.Data == nil {
		return nil, newError(fmt.Sprintf("Unable to generate QRCode: %v (code %d %s)", err, gen.ErrorCode, gen.ErrorMessage))
	}
	qr := gen.Data
	qr.Image = trimDataURL(qr.Image)

	act := LoginQRContinue
	if cb != nil {
		act = cb(LoginQREvent{Type: LoginQREventQRCodeGenerated, QRCode: qr})
	} else {
		if err := qr.SaveToFile(qrPath); err != nil {
			return nil, err
		}
		s.log().Info("Scan the QR code to proceed with login", "path", qrPath)
	}
	if err := qrAction(act); err != nil {
		return nil, err
	}

	waitCtx, cancel := context.WithTimeout(ctx, 100*time.Second)
	defer cancel()
	expired := func(err error) ([]Cookie, error) {
		if ctx.Err() == nil && errors.Is(waitCtx.Err(), context.DeadlineExceeded) {
			s.log().Info("QR expired!")
			if cb == nil {
				return nil, errQRRetry
			}
			if cb(LoginQREvent{Type: LoginQREventQRCodeExpired}) == LoginQRAbort {
				return nil, ErrLoginQRAborted
			}
			return nil, errQRRetry
		}
		return nil, err
	}

	codeForm := url.Values{"code": {qr.Code}, "continue": {"https://chat.zalo.me/"}, "v": {version}}
	scan, err := qrPoll[QRScanData](waitCtx, s, "https://id.zalo.me/account/authen/qr/waiting-scan", codeForm)
	if err != nil {
		return expired(err)
	}
	if scan.Data == nil {
		return nil, newError("Cannot get scan result")
	}
	if cb != nil {
		if err := qrAction(cb(LoginQREvent{Type: LoginQREventQRCodeScanned, Scanned: scan.Data})); err != nil {
			return nil, err
		}
	}

	s.log().Info("Please confirm on your phone")
	confirmForm := url.Values{"code": {qr.Code}, "gToken": {""}, "gAction": {"CONFIRM_QR"}, "continue": {"https://chat.zalo.me/"}, "v": {version}}
	confirm, err := qrPoll[json.RawMessage](waitCtx, s, "https://id.zalo.me/account/authen/qr/waiting-confirm", confirmForm)
	if err != nil {
		return expired(err)
	}
	cancel()
	if confirm.ErrorCode == -13 {
		if cb == nil {
			s.log().Error("QRCode login declined")
			return nil, ErrLoginQRDeclined
		}
		if cb(LoginQREvent{Type: LoginQREventQRCodeDeclined, Code: qr.Code}) == LoginQRRetry {
			return nil, errQRRetry
		}
		return nil, ErrLoginQRDeclined
	} else if confirm.ErrorCode != 0 {
		return nil, newError(fmt.Sprintf("An error has occurred: %d %s", confirm.ErrorCode, confirm.ErrorMessage))
	}

	ch := qrHeaders(s.UserAgent, refChat, true)
	resp, err = s.Request(ctx, http.MethodGet, "https://id.zalo.me/account/checksession?continue=https%3A%2F%2Fchat.zalo.me%2Findex.html", nil, ch)
	if err != nil {
		return nil, newError("Cannot get session, login failed")
	}
	resp.Body.Close()
	s.log().Info("Successfully logged into the account", "name", scan.Data.DisplayName)

	uh := qrHeaders(s.UserAgent, "https://chat.zalo.me/", false)
	uh.Set("sec-fetch-site", "same-site")
	uh.Del("Content-Type")
	resp, err = s.Request(ctx, http.MethodGet, "https://jr.chat.zalo.me/jr/userinfo", nil, uh)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var info qrResult[struct {
		Logged bool `json:"logged"`
	}]
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil || info.Data == nil {
		return nil, newError("Can't get account info")
	}
	if !info.Data.Logged {
		return nil, newError("Can't login")
	}
	return s.jar.export(), nil
}

func qrAction(a LoginQRAction) error {
	switch a {
	case LoginQRAbort:
		return ErrLoginQRAborted
	case LoginQRRetry:
		return errQRRetry
	}
	return nil
}
