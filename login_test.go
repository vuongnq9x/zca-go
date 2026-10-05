package zca

import "testing"

func TestQRHeadersFromUA(t *testing.T) {
	for _, c := range []struct{ ua, platform, ver string }{
		{"", "Windows", "130"},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36", "Windows", "133"},
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36", "macOS", "131"},
		{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36", "Linux", "129"},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:133.0) Gecko/20100101 Firefox/133.0", "Windows", "130"},
	} {
		h := qrHeaders(c.ua, "https://chat.zalo.me/", false)
		wantUA := `"Chromium";v="` + c.ver + `", "Google Chrome";v="` + c.ver + `", "Not?A_Brand";v="99"`
		if h.Get("sec-ch-ua-platform") != `"`+c.platform+`"` || h.Get("sec-ch-ua") != wantUA {
			t.Errorf("%q: platform=%s ua=%s", c.ua, h.Get("sec-ch-ua-platform"), h.Get("sec-ch-ua"))
		}
	}
}
