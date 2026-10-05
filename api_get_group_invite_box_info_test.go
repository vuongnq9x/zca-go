package zca

import "testing"

func TestGetGroupInviteBoxInfoParseTopic(t *testing.T) {
	topic := &GroupTopic{Params: `{"title":"t","extra":"{\"a\":1}"}`}
	if err := getGroupInviteBoxInfoParseTopic(topic); err != nil {
		t.Fatal(err)
	}
	m := topic.Params.(map[string]any)
	if m["title"] != "t" || m["extra"].(map[string]any)["a"] != float64(1) {
		t.Fatalf("got %#v", m)
	}
}
