# zca-go

Bản port Go của [zca-js](https://github.com/RFS-ADRENO/zca-js) v2.2.0 (Unofficial Zalo API). Package `zca`, module `github.com/vuongnq9x/zca-go`.

```go
z := zca.New(zca.Options{})                       // Logger, HTTPClient (proxy), SelfListen, ImageMetadataGetter
api, err := z.Login(ctx, zca.Credentials{IMEI: imei, Cookie: cookies, UserAgent: ua})
// hoặc: api, err := z.LoginQR(ctx, zca.LoginQROptions{}, nil)  // lưu qr.png

api.SendMessage(ctx, zca.MessageContent{Msg: "hi"}, threadID, zca.ThreadTypeUser)
api.Listener.OnMessage = func(m *zca.Message) { ... }
api.Listener.Start(true)
```

Ví dụ đầy đủ: `examples/echo`. Cookie JSON xuất từ trình duyệt đọc bằng `zca.ParseCookies`.

## Khác với zca-js

- Mọi API là method trên `*API`, tham số đầu `context.Context`; `string | string[]` → `[]string`.
- Listener dùng các hàm `On*` thay EventEmitter; handler chạy trên goroutine của listener.
- Gửi file (video/others) chờ sự kiện `file_done` qua websocket ⇒ phải `Listener.Start` trước.
- `getOwnId` → `api.OwnID()`, `getContext` → `api.Session`, `getCookie` → `api.Cookies()`.
  `custom` không cần: dùng `api.MakeURL`, `api.EncodeAES`, `api.Request`, `api.Resolve`.
- Bỏ `checkUpdate` (kiểm npm). Upload chunk gửi tuần tự thay vì song song.
- Trường số nhận 0 = dùng mặc định của TS (Go không phân biệt "không truyền").
