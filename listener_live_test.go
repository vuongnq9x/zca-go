package zca

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// fakeWS counts connections; dropFirst closes the first TCP conn abruptly (=> 1006).
func fakeWS(t *testing.T, dropFirst bool) (*httptest.Server, *atomic.Int32) {
	var n atomic.Int32
	up := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := up.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		if n.Add(1) == 1 && dropFirst {
			c.UnderlyingConn().Close()
			return
		}
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				return
			}
		}
	}))
	t.Cleanup(srv.Close)
	return srv, &n
}

func liveListener(srvURL string) *Listener {
	s := newSession(Options{})
	s.UserAgent = "test"
	s.Settings.Features.Socket.Retries = map[string]struct {
		Max   int `json:"max"`
		Times any `json:"times"`
	}{"internal": {Max: 3, Times: []any{float64(10)}}}
	return newListener(s, []string{"ws" + strings.TrimPrefix(srvURL, "http")})
}

func TestListenerStopDoesNotReconnect(t *testing.T) {
	srv, n := fakeWS(t, false)
	l := liveListener(srv.URL)
	closed := make(chan CloseReason, 1)
	l.OnClosed = func(c CloseReason, _ string) { closed <- c }
	if err := l.Start(true); err != nil {
		t.Fatal(err)
	}
	l.Stop()
	select {
	case c := <-closed:
		if c != CloseReasonManualClosure {
			t.Fatalf("close code = %d, want manual", c)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("OnClosed not called")
	}
	time.Sleep(100 * time.Millisecond)
	if got := n.Load(); got != 1 {
		t.Fatalf("connections = %d, want 1 (Stop must not reconnect)", got)
	}
}

func TestListenerReconnectsOnAbnormalClosure(t *testing.T) {
	srv, n := fakeWS(t, true)
	l := liveListener(srv.URL)
	connected := make(chan struct{}, 2)
	l.OnConnected = func() { connected <- struct{}{} }
	if err := l.Start(true); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		select {
		case <-connected:
		case <-time.After(2 * time.Second):
			t.Fatalf("connections = %d, want reconnect after 1006", n.Load())
		}
	}
	l.Stop()
}
