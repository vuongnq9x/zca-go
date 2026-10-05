package zca

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/url"
	"sort"
	"strings"
)

var zeroIV = make([]byte, aes.BlockSize)

func md5Hex(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

func aesCBCEncrypt(key, plain []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	pad := aes.BlockSize - len(plain)%aes.BlockSize
	buf := append(bytes.Clone(plain), bytes.Repeat([]byte{byte(pad)}, pad)...)
	cipher.NewCBCEncrypter(block, zeroIV).CryptBlocks(buf, buf)
	return buf, nil
}

func aesCBCDecrypt(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 || len(data)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("invalid ciphertext length %d", len(data))
	}
	out := bytes.Clone(data)
	cipher.NewCBCDecrypter(block, zeroIV).CryptBlocks(out, out)
	pad := int(out[len(out)-1])
	if pad == 0 || pad > aes.BlockSize || pad > len(out) {
		return nil, fmt.Errorf("invalid padding")
	}
	return out[:len(out)-pad], nil
}

// encodeAES encrypts data with the session secret key (base64), AES-CBC zero IV, base64 output.
func encodeAES(secretKey, data string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		return "", err
	}
	enc, err := aesCBCEncrypt(key, []byte(data))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(enc), nil
}

// decodeAES reverses encodeAES; data may be URI-encoded base64.
func decodeAES(secretKey, data string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		return nil, err
	}
	return decodeCBCBase64(key, data)
}

func decodeCBCBase64(key []byte, data string) ([]byte, error) {
	// decodeURIComponent semantics: '+' stays '+'.
	if unescaped, err := url.PathUnescape(data); err == nil {
		data = unescaped
	}
	raw, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return aesCBCDecrypt(key, raw)
}

// getSignKey: md5("zsecure" + type + values sorted by key).
func getSignKey(typ string, params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	b.WriteString("zsecure" + typ)
	for _, k := range keys {
		b.WriteString(params[k])
	}
	return md5Hex(b.String())
}

// paramsEncryptor mirrors zca-js ParamsEncryptor used by the login endpoints.
type paramsEncryptor struct {
	zcid, zcidExt, encryptKey string
}

func newParamsEncryptor(typ int, imei string, firstLaunchTime int64) (*paramsEncryptor, error) {
	return buildParamsEncryptor(typ, imei, firstLaunchTime, randomHex(6, 12))
}

func buildParamsEncryptor(typ int, imei string, firstLaunchTime int64, zcidExt string) (*paramsEncryptor, error) {
	if typ == 0 || imei == "" || firstLaunchTime == 0 {
		return nil, newError("createZcid: missing params")
	}
	zcid, err := aesCBCEncrypt([]byte("3FC4F0D2AB50057BCE0D90D9187A22B1"), fmt.Appendf(nil, "%d,%s,%d", typ, imei, firstLaunchTime))
	if err != nil {
		return nil, err
	}
	p := &paramsEncryptor{zcid: strings.ToUpper(hex.EncodeToString(zcid)), zcidExt: zcidExt}
	n := strings.ToUpper(md5Hex(p.zcidExt))
	nEven, _ := splitEvenOdd(n)
	aEven, aOdd := splitEvenOdd(p.zcid)
	reverse(aOdd)
	p.encryptKey = string(nEven[:8]) + string(aEven[:12]) + string(aOdd[:12])
	return p, nil
}

func (p *paramsEncryptor) params() map[string]string {
	return map[string]string{"zcid": p.zcid, "zcid_ext": p.zcidExt, "enc_ver": "v2"}
}

func splitEvenOdd(s string) (even, odd []byte) {
	for i := 0; i < len(s); i++ {
		if i%2 == 0 {
			even = append(even, s[i])
		} else {
			odd = append(odd, s[i])
		}
	}
	return
}

func reverse(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

// randomHex returns a random lowercase hex string with length in [min, max].
func randomHex(min, max int) string {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	size := min + int(n.Int64())
	buf := make([]byte, (size+1)/2)
	rand.Read(buf)
	return hex.EncodeToString(buf)[:size]
}

// decryptResp decrypts login responses with the UTF-8 encrypt key.
func decryptResp(key, data string) ([]byte, error) {
	return decodeCBCBase64([]byte(key), data)
}

// decodeEventData decodes a websocket event payload {data, encrypt}.
func decodeEventData(data string, encrypt int, cipherKey string) ([]byte, error) {
	if encrypt < 0 || encrypt > 3 {
		return nil, newError(fmt.Sprintf("invalid encrypt type, expected 0-3 but got %d", encrypt))
	}
	if encrypt == 0 {
		return []byte(data), nil
	}
	if encrypt != 1 {
		if s, err := url.PathUnescape(data); err == nil {
			data = s
		}
	}
	buf, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	if encrypt != 1 {
		if cipherKey == "" || len(buf) < 48 {
			return nil, newError("invalid data length or missing cipher key")
		}
		key, err := base64.StdEncoding.DecodeString(cipherKey)
		if err != nil {
			return nil, err
		}
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}
		gcm, err := cipher.NewGCMWithNonceSize(block, 16)
		if err != nil {
			return nil, err
		}
		if buf, err = gcm.Open(nil, buf[:16], buf[32:], buf[16:32]); err != nil {
			return nil, err
		}
	}
	if encrypt == 3 {
		return buf, nil
	}
	return inflate(buf)
}

// inflate accepts zlib or gzip like pako.inflate's auto-detection.
func inflate(b []byte) ([]byte, error) {
	var r io.ReadCloser
	var err error
	if len(b) > 1 && b[0] == 0x1f && b[1] == 0x8b {
		r, err = gzip.NewReader(bytes.NewReader(b))
	} else {
		r, err = zlib.NewReader(bytes.NewReader(b))
	}
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}

func mustJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err) // ponytail: params are plain maps/structs built by us
	}
	return string(b)
}

// GenerateZaloUUID builds an imei like zca-js generateZaloUUID.
func GenerateZaloUUID(userAgent string) string {
	return newUUID() + "-" + md5Hex(userAgent)
}

func newUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// EncryptPin hashes a 4-digit PIN (md5 hex) as used by hidden conversations.
func EncryptPin(pin string) string { return md5Hex(pin) }

// ValidatePin reports whether pin matches the hash from GetHiddenConversations.
func ValidatePin(encryptedPin, pin string) bool { return md5Hex(pin) == encryptedPin }
