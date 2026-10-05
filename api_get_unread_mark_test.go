package zca

import "testing"

func TestGetUnreadMarkDecode(t *testing.T) {
	for _, raw := range []string{
		`{"data":"{\"convsGroup\":[{\"id\":1}],\"convsUser\":[]}","status":2}`,
		`{"data":{"convsGroup":[{"id":1}],"convsUser":[]},"status":2}`,
	} {
		var out GetUnreadMarkResponse
		if err := getUnreadMarkDecode([]byte(raw), &out); err != nil {
			t.Fatal(err)
		}
		if out.Status != 2 || len(out.Data.ConvsGroup) != 1 || out.Data.ConvsGroup[0].ID != 1 {
			t.Fatalf("bad decode: %+v", out)
		}
	}
}
