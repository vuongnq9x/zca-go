package zca

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"reflect"
	"testing"
)

func TestSendMessageStyles(t *testing.T) {
	p := map[string]any{}
	sendMessageStyles(p, []Style{
		{Start: 0, Len: 3, St: TextStyleBold},
		{Start: 1, Len: 2, St: TextStyleIndent},
		{Start: 2, Len: 1, St: TextStyleIndent, IndentSize: 3},
	})
	want := `{"styles":[{"start":0,"len":3,"st":"b"},{"start":1,"len":2,"st":"ind_10"},{"start":2,"len":1,"st":"ind_30"}],"ver":0}`
	if p["textProperties"] != want {
		t.Fatalf("got %v", p["textProperties"])
	}
	p = map[string]any{}
	sendMessageStyles(p, nil)
	if _, ok := p["textProperties"]; ok {
		t.Fatal("nil styles must not set textProperties")
	}
}

func TestSendMessageMentions(t *testing.T) {
	ms := []Mention{{Pos: 0, UID: "1", Len: 2}, {Pos: -1, UID: "2", Len: 1}, {Pos: 3, UID: "-1", Len: 1}, {Pos: 0, UID: "", Len: 1}}
	got, err := sendMessageMentions(ThreadTypeGroup, "@a @all", ms)
	if err != nil {
		t.Fatal(err)
	}
	want := []sendMessageMention{{Pos: 0, UID: "1", Len: 2, Type: 0}, {Pos: 3, UID: "-1", Len: 1, Type: 1}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v", got)
	}
	if got, _ := sendMessageMentions(ThreadTypeUser, "hi", ms); len(got) != 0 {
		t.Fatal("user threads drop mentions")
	}
	if _, err := sendMessageMentions(ThreadTypeGroup, "ab", []Mention{{UID: "1", Len: 3}}); err == nil {
		t.Fatal("expected length error")
	}
	// UTF-16 length like JS: "😀" counts 2.
	if _, err := sendMessageMentions(ThreadTypeGroup, "😀", []Mention{{UID: "1", Len: 2}}); err != nil {
		t.Fatal(err)
	}
}

func TestClientMessageType(t *testing.T) {
	for in, want := range map[string]int{"webchat": 1, "chat.photo": 32, "chat.link": 38, "share.file": 46, "chat.gif": 49, "x": 1} {
		if got := clientMessageType(in); got != want {
			t.Errorf("%s: got %d want %d", in, got, want)
		}
	}
}

func TestPrepareQMSG(t *testing.T) {
	q := &SendMessageQuote{MsgType: "chat.todo", Content: map[string]any{"params": `{"item":{"content":"buy milk"}}`}}
	if v, err := prepareQMSG(q); err != nil || v != "buy milk" {
		t.Fatalf("todo: %v %v", v, err)
	}
	if v, _ := prepareQMSG(&SendMessageQuote{MsgType: "chat.photo", Content: map[string]any{}}); v != "" {
		t.Fatalf("got %v", v)
	}
	att := prepareQMSGAttach(q)
	if b, _ := json.Marshal(att); string(b) != `{"properties":{"color":0,"ext":"{\"shouldParseLinkOrContact\":0}","size":0,"subType":0,"type":0}}` {
		t.Fatalf("todo attach %s", b)
	}
	photo := &SendMessageQuote{MsgType: "chat.photo", Content: map[string]any{"thumb": "t", "href": "h", "title": ""}}
	want := map[string]any{"thumb": "t", "href": "h", "title": "", "thumbUrl": "t", "oriUrl": "h", "normalUrl": "h"}
	if got := prepareQMSGAttach(photo); !reflect.DeepEqual(got, want) {
		t.Fatalf("photo attach %v", got)
	}
	if prepareQMSGAttach(&SendMessageQuote{Content: "text"}) != nil {
		t.Fatal("string content without propertyExt must be dropped")
	}
	// Struct content (e.g. TAttachmentContent) is accepted too.
	if got := prepareQMSGAttach(&SendMessageQuote{Content: TAttachmentContent{Thumb: "t"}}).(map[string]any); got["thumbUrl"] != "t" {
		t.Fatalf("struct attach %v", got)
	}
}

func TestMediaHelpers(t *testing.T) {
	for in, want := range map[string]string{"a/b.JPG": "JPG", ".bashrc": "", "x": "", "a.tar.gz": "gz", "c.": ""} {
		if got := mediaFileExt(in); got != want {
			t.Errorf("%s: got %q want %q", in, got, want)
		}
	}
	if mediaMD5([]byte("abcdef"), 3) != md5Hex("abc") || mediaMD5([]byte("ab"), 9) != md5Hex("ab") {
		t.Fatal("md5 prefix")
	}
	body, h := mediaMultipart("chunkContent", `tệp "1".bin`, "application/octet-stream", []byte("data"))
	_, mp, _ := mime.ParseMediaType(h.Get("Content-Type"))
	part, err := multipart.NewReader(bytes.NewReader(body), mp["boundary"]).NextPart()
	if err != nil {
		t.Fatal(err)
	}
	got, _ := io.ReadAll(part)
	if part.FormName() != "chunkContent" || part.FileName() != `tệp "1".bin` || string(got) != "data" {
		t.Fatalf("multipart %q %q %q", part.FormName(), part.FileName(), got)
	}
}
