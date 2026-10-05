// Echo bot: logs in by QR (or credentials.json if present) and echoes user messages.
package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	zca "github.com/vuongnq9x/zca-go"
)

func main() {
	ctx := context.Background()
	z := zca.New(zca.Options{})

	var api *zca.API
	var err error
	if b, rerr := os.ReadFile("credentials.json"); rerr == nil {
		var cred zca.Credentials
		if err := json.Unmarshal(b, &cred); err != nil {
			log.Fatal(err)
		}
		api, err = z.Login(ctx, cred)
	} else {
		api, err = z.LoginQR(ctx, zca.LoginQROptions{}, func(ev zca.LoginQREvent) zca.LoginQRAction {
			switch ev.Type {
			case zca.LoginQREventQRCodeGenerated:
				ev.QRCode.SaveToFile("qr.png")
				log.Println("scan qr.png with the Zalo app")
			case zca.LoginQREventGotLoginInfo:
				b, _ := json.MarshalIndent(ev.LoginInfo, "", "  ")
				os.WriteFile("credentials.json", b, 0o600)
			}
			return zca.LoginQRContinue
		})
	}
	if err != nil {
		log.Fatal(err)
	}

	l := api.Listener
	l.OnMessage = func(m *zca.Message) {
		text, ok := m.Data.Content.(string)
		if !ok || m.IsSelf {
			return
		}
		if _, err := api.SendMessage(ctx, zca.MessageContent{Msg: text, Quote: zca.QuoteFromMessage(m)}, m.ThreadID, m.Type); err != nil {
			log.Println("send:", err)
		}
	}
	l.OnError = func(err error) { log.Println("listener:", err) }
	done := make(chan struct{})
	l.OnClosed = func(code zca.CloseReason, reason string) { log.Println("closed", code, reason); close(done) }
	if err := l.Start(true); err != nil {
		log.Fatal(err)
	}
	<-done
}
