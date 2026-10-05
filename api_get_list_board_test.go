package zca

import "testing"

func TestGetListBoardDecode(t *testing.T) {
	raw := []byte(`{"count":2,"items":[
		{"boardType":1,"data":{"id":"n1","params":"{\"title\":\"hi\"}"}},
		{"boardType":3,"data":{"poll_id":7,"question":"q"}}]}`)
	got, err := getListBoardDecode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if n := got.Items[0].Data.(*NoteDetail); n.ID != "n1" || n.Params.Title != "hi" {
		t.Fatalf("note %+v", n)
	}
	if p := got.Items[1].Data.(*PollDetail); p.PollID != 7 || p.Question != "q" {
		t.Fatalf("poll %+v", p)
	}
}
