package zca

import (
	"context"
	"net/http"
)

// SendFriendRequest sends a friend request. Zalo codes: 225 already friends, 215 maybe blocked,
// 222 they already requested you (treated as acceptance).
func (a *API) SendFriendRequest(ctx context.Context, msg, userID string) (string, error) {
	params := map[string]any{
		"toid": userID, "msg": msg, "reqsrc": 30, "imei": a.IMEI, "language": a.Language,
		"srcParams": mustJSON(map[string]any{"uidTo": userID}),
	}
	return call[string](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("friend")+"/api/friend/sendreq", nil, true), params)
}
