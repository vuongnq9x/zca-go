package zca

import (
	"strings"
	"testing"
)

// Vectors generated with Node's crypto, matching CryptoJS behaviour in zca-js.
func TestCryptoVectors(t *testing.T) {
	key := "BwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwc="
	msg := `{"a":1,"b":"xin chào"}`
	enc, err := encodeAES(key, msg)
	if err != nil || enc != "BKIw23EDAsYBuQzJDjrLzFMkmkPUyKGAxrpeAxjedr8=" {
		t.Fatalf("encodeAES = %q, %v", enc, err)
	}
	dec, err := decodeAES(key, strings.ReplaceAll(enc, "=", "%3D"))
	if err != nil || string(dec) != msg {
		t.Fatalf("decodeAES = %q, %v", dec, err)
	}

	p, err := buildParamsEncryptor(30, "imei-x", 1700000000000, "abc123def")
	if err != nil {
		t.Fatal(err)
	}
	if p.zcid != "4B24483CF41B7E2C66180976DE9AF8A264A4DB7299EEE5D7268AA95AFE51FBC8" || p.encryptKey != "E139DD174243F17261078B1EA9A675E9" {
		t.Fatalf("zcid/encryptKey = %s / %s", p.zcid, p.encryptKey)
	}

	sign := getSignKey("getlogininfo", map[string]string{"zcid": "Z", "type": "30", "client_version": "685", "params": "P"})
	if sign != "b9d0c7d8a5449edfbaddce76ad794884" {
		t.Fatalf("sign = %s", sign)
	}

	ev, err := decodeEventData("RyblC9wQkny87E6W0m5mIHfeJL2aDMjR9mtz%2BD9Bf6I5npCgR0xKSNxTm2duEeVqftgNYHpdXZaimXh1k1Ft9qVNA32mLhlHYsES%2BFs4%2FRi1B911nl1dT9f6zVw%3D", 2, "lZE0jRHgLMQHzzslT6kdKw==")
	if err != nil || string(ev) != `{"data":{"msgs":[]},"x":12345678901234567890}` {
		t.Fatalf("decodeEventData = %q, %v", ev, err)
	}
	if h := randomHex(6, 12); len(h) < 6 || len(h) > 12 {
		t.Fatalf("randomHex len %d", len(h))
	}
}
