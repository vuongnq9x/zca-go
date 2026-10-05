package zca

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Cookie accepts both browser-extension exports (name, expirationDate) and
// tough-cookie serialized cookies (key, expires) like zca-js Credentials.cookie.
type Cookie struct {
	Name           string  `json:"name,omitempty"`
	Key            string  `json:"key,omitempty"`
	Value          string  `json:"value"`
	Domain         string  `json:"domain"`
	Path           string  `json:"path,omitempty"`
	Expires        string  `json:"expires,omitempty"`
	ExpirationDate float64 `json:"expirationDate,omitempty"`
	Secure         bool    `json:"secure,omitempty"`
	HTTPOnly       bool    `json:"httpOnly,omitempty"`
	HostOnly       bool    `json:"hostOnly,omitempty"`
}

// ParseCookies parses a JSON cookie export: an array, or {"url": ..., "cookies": [...]}.
func ParseCookies(data []byte) ([]Cookie, error) {
	var arr []Cookie
	if err := json.Unmarshal(data, &arr); err == nil {
		return arr, nil
	}
	var wrapped struct {
		Cookies []Cookie `json:"cookies"`
	}
	if err := json.Unmarshal(data, &wrapped); err != nil {
		return nil, err
	}
	return wrapped.Cookies, nil
}

// cookieJar is a minimal jar that, unlike net/http/cookiejar, can export its cookies.
// ponytail: every cookie is domain-matched (no hostOnly/public-suffix rules); fine for *.zalo.me.
type cookieJar struct {
	mu      sync.Mutex
	cookies []*http.Cookie
}

func (j *cookieJar) set(c *http.Cookie, host string) {
	c.Domain = strings.TrimPrefix(c.Domain, ".")
	if c.Domain == "" {
		c.Domain = host
	}
	if c.Path == "" {
		c.Path = "/"
	}
	j.mu.Lock()
	defer j.mu.Unlock()
	for i, old := range j.cookies {
		if old.Name == c.Name && old.Domain == c.Domain && old.Path == c.Path {
			j.cookies = append(j.cookies[:i], j.cookies[i+1:]...)
			break
		}
	}
	if c.MaxAge < 0 || (!c.Expires.IsZero() && c.Expires.Before(time.Now())) {
		return
	}
	j.cookies = append(j.cookies, c)
}

func (j *cookieJar) header(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	host, path := u.Hostname(), u.Path
	if path == "" {
		path = "/"
	}
	j.mu.Lock()
	defer j.mu.Unlock()
	var parts []string
	now := time.Now()
	for _, c := range j.cookies {
		if !c.Expires.IsZero() && c.Expires.Before(now) {
			continue
		}
		if host != c.Domain && !strings.HasSuffix(host, "."+c.Domain) {
			continue
		}
		if !strings.HasPrefix(path, c.Path) {
			continue
		}
		parts = append(parts, c.Name+"="+c.Value)
	}
	return strings.Join(parts, "; ")
}

func (j *cookieJar) load(cookies []Cookie) {
	for _, c := range cookies {
		name := c.Key
		if name == "" {
			name = c.Name
		}
		hc := &http.Cookie{Name: name, Value: c.Value, Domain: c.Domain, Path: c.Path, Secure: c.Secure, HttpOnly: c.HTTPOnly}
		if c.ExpirationDate > 0 {
			hc.Expires = time.Unix(int64(c.ExpirationDate), 0)
		} else if t, err := time.Parse(time.RFC3339, c.Expires); err == nil {
			hc.Expires = t
		}
		host := strings.TrimPrefix(c.Domain, ".")
		if host == "" {
			host = "chat.zalo.me"
		}
		j.set(hc, host)
	}
}

func (j *cookieJar) export() []Cookie {
	j.mu.Lock()
	defer j.mu.Unlock()
	out := make([]Cookie, 0, len(j.cookies))
	for _, c := range j.cookies {
		ec := Cookie{Name: c.Name, Key: c.Name, Value: c.Value, Domain: c.Domain, Path: c.Path, Secure: c.Secure, HTTPOnly: c.HttpOnly}
		if !c.Expires.IsZero() {
			ec.Expires = c.Expires.UTC().Format(time.RFC3339)
			ec.ExpirationDate = float64(c.Expires.Unix())
		}
		out = append(out, ec)
	}
	return out
}
