package zca

import "testing"

func TestColorAndHolderName(t *testing.T) {
	if got := HexToNegativeColor("#00FF00"); got != -16711936 {
		t.Errorf("HexToNegativeColor = %d", got)
	}
	if got := HexToNegativeColor("7f00ff00"); got != 0x7f00ff00 {
		t.Errorf("HexToNegativeColor alpha = %d", got)
	}
	if got := NegativeColorToHex(-16711936); got != "#00ff00" {
		t.Errorf("NegativeColorToHex = %q", got)
	}
	if got := NegativeColorToHex(-16777216); got != "#000000" {
		t.Errorf("NegativeColorToHex black = %q", got)
	}
	for in, want := range map[string]string{
		"Nguyễn Văn Đức":     "NGUYEN VAN DUC",
		"  trần-thị  ỷ 123 ": "TRAN THI Y 123",
		"Lê":                 "",
		"":                   "",
	} {
		if got := NormalizeHolderName(in); got != want {
			t.Errorf("NormalizeHolderName(%q) = %q, want %q", in, got, want)
		}
	}
}
