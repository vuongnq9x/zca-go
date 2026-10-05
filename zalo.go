// Package zca is an unofficial Zalo API client, a Go port of zca-js.
package zca

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type Credentials struct {
	IMEI      string   `json:"imei"`
	Cookie    []Cookie `json:"cookie"`
	UserAgent string   `json:"userAgent"`
	Language  string   `json:"language,omitempty"` // default "vi"
}

type Zalo struct {
	Options Options
}

func New(opts Options) *Zalo { return &Zalo{Options: opts} }

// API is a logged-in client. All Zalo endpoints are methods on it.
type API struct {
	*Session
	Listener *Listener
}

// Login logs in with existing cookies + imei + user agent.
func (z *Zalo) Login(ctx context.Context, cred Credentials) (*API, error) {
	return loginCookie(ctx, newSession(z.Options), cred)
}

func loginCookie(ctx context.Context, s *Session, cred Credentials) (*API, error) {
	if cred.IMEI == "" || len(cred.Cookie) == 0 || cred.UserAgent == "" {
		return nil, newError("Missing required params")
	}
	s.IMEI, s.UserAgent, s.Language = cred.IMEI, cred.UserAgent, cred.Language
	if s.Language == "" {
		s.Language = "vi"
	}
	s.jar = &cookieJar{}
	s.jar.load(cred.Cookie)

	info, err := login(ctx, s)
	if err != nil {
		return nil, err
	}
	srv, err := getServerInfo(ctx, s)
	if err != nil {
		return nil, err
	}

	s.SecretKey, _ = info["zpw_enk"].(string)
	s.UID = fmt.Sprint(info["uid"])
	if s.SecretKey == "" {
		return nil, newError("Đăng nhập thất bại")
	}
	settings := srv.Setttings
	if len(settings) == 0 || string(settings) == "null" {
		settings = srv.Settings
	}
	if err := json.Unmarshal(settings, &s.Settings); err != nil {
		return nil, fmt.Errorf("zca: parse settings: %w", err)
	}
	s.ExtraVer = srv.ExtraVer
	s.LoginInfo = info
	if err := remarshal(info["zpw_service_map_v3"], &s.ServiceMap); err != nil {
		return nil, fmt.Errorf("zca: parse service map: %w", err)
	}
	var wsURLs []string
	_ = remarshal(info["zpw_ws"], &wsURLs)

	s.log().Info("Logged in", "uid", s.UID)
	api := &API{Session: s}
	api.Listener = newListener(s, wsURLs)
	return api, nil
}

func remarshal(in, out any) error {
	b, err := json.Marshal(in)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}

type LoginQROptions struct {
	UserAgent string // default Firefox 133 UA
	Language  string // default "vi"
	QRPath    string // default "qr.png", used when callback is nil
}

// LoginQR logs in by QR code. With a nil callback the QR is saved to QRPath and
// regenerated on expiry.
func (z *Zalo) LoginQR(ctx context.Context, opts LoginQROptions, cb LoginQRCallback) (*API, error) {
	if opts.UserAgent == "" {
		opts.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:133.0) Gecko/20100101 Firefox/133.0"
	}
	if opts.Language == "" {
		opts.Language = "vi"
	}
	if opts.QRPath == "" {
		opts.QRPath = "qr.png"
	}
	s := newSession(z.Options)
	s.UserAgent = opts.UserAgent
	var cookies []Cookie
	for {
		var err error
		cookies, err = loginQROnce(ctx, s, opts.QRPath, cb)
		if errors.Is(err, errQRRetry) {
			continue
		}
		if err != nil {
			return nil, err
		}
		break
	}
	cred := Credentials{IMEI: GenerateZaloUUID(opts.UserAgent), Cookie: cookies, UserAgent: opts.UserAgent, Language: opts.Language}
	if cb != nil {
		cb(LoginQREvent{Type: LoginQREventGotLoginInfo, LoginInfo: &cred})
	}
	return loginCookie(ctx, s, cred)
}

// OwnID returns the logged-in user id (zca-js getOwnId).
func (a *API) OwnID() string { return a.UID }
