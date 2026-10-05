package zca

import (
	"fmt"
	"strconv"
	"strings"
)

// HexToNegativeColor converts "#RRGGBB"/"AARRGGBB" (with or without #) to Zalo's signed ARGB int.
// 6-digit input gets alpha FF. Invalid hex returns 0 (TS returns NaN).
func HexToNegativeColor(hex string) int64 {
	h := strings.TrimPrefix(hex, "#")
	if len(h) == 6 {
		h = "FF" + h
	}
	v, err := strconv.ParseInt(h, 16, 64)
	if err != nil {
		return 0
	}
	if v > 0x7fffffff {
		v -= 1 << 32
	}
	return v
}

// NegativeColorToHex converts Zalo's signed ARGB int to "#rrggbb" (alpha dropped, lowercase).
func NegativeColorToHex(color int64) string {
	s := strconv.FormatInt(color+1<<32, 16)
	if len(s) > 6 {
		s = s[len(s)-6:]
	}
	return "#" + fmt.Sprintf("%06s", s)
}

// ponytail: hand table instead of x/text NFD; covers Vietnamese + Latin-1, other accents become spaces.
var holderNameFold = func() *strings.Replacer {
	groups := map[string]string{
		"ÀÁẢÃẠĂẰẮẲẴẶÂẦẤẨẪẬÄÅĀ": "A", "ÈÉẺẼẸÊỀẾỂỄỆËĒ": "E", "ÌÍỈĨỊÎÏĪ": "I",
		"ÒÓỎÕỌÔỒỐỔỖỘƠỜỚỞỠỢÖŌ": "O", "ÙÚỦŨỤƯỪỨỬỮỰÛÜŪ": "U", "ỲÝỶỸỴŸ": "Y",
		"Đ": "D", "Ç": "C", "Ñ": "N",
	}
	var pairs []string
	for from, to := range groups {
		for _, r := range from {
			pairs = append(pairs, string(r), to)
		}
	}
	return strings.NewReplacer(pairs...)
}()

// NormalizeHolderName converts a bank holder name to unaccented uppercase A-Z0-9 words.
// Returns "" when the result is shorter than 5 characters (TS undefined).
func NormalizeHolderName(input string) string {
	s := holderNameFold.Replace(strings.ToUpper(input))
	s = strings.Join(strings.FieldsFunc(s, func(r rune) bool {
		return !(r >= 'A' && r <= 'Z' || r >= '0' && r <= '9')
	}), " ")
	if len(s) < 5 {
		return ""
	}
	return s
}
