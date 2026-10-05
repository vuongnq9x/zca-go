package zca

import "testing"

func TestGetFriendOnlinesUnwrap(t *testing.T) {
	d := &GetFriendOnlinesResponse{Onlines: []GetFriendOnlinesStatus{
		{UserID: "1", Status: `{"status":"busy"}`},
		{UserID: "2", Status: `{"other":1}`},
	}}
	if err := getFriendOnlinesUnwrap(d); err != nil {
		t.Fatal(err)
	}
	if d.Onlines[0].Status != "busy" || d.Onlines[1].Status != `{"other":1}` {
		t.Fatalf("got %+v", d.Onlines)
	}
	if getFriendOnlinesUnwrap(&GetFriendOnlinesResponse{Onlines: []GetFriendOnlinesStatus{{Status: "x"}}}) == nil {
		t.Fatal("want parse error")
	}
}
