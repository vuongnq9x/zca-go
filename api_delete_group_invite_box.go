package zca

import (
	"context"
	"net/http"
)

type DeleteGroupInviteBoxResponse struct {
	DelInvitaionIDs []string `json:"delInvitaionIds"`
	ErrMap          map[string]struct {
		Err int64 `json:"err"`
	} `json:"errMap"`
}

// DeleteGroupInviteBox deletes group invitation(s), optionally blocking future invites.
func (a *API) DeleteGroupInviteBox(ctx context.Context, groupIDs []string, blockFutureInvite bool) (*DeleteGroupInviteBoxResponse, error) {
	inv := make([]map[string]string, len(groupIDs))
	for i, id := range groupIDs {
		inv[i] = map[string]string{"grid": id}
	}
	block := 0
	if blockFutureInvite {
		block = 1
	}
	return call[*DeleteGroupInviteBoxResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group")+"/api/group/inv-box/mdel-inv", nil, true),
		map[string]any{"invitations": mustJSON(inv), "block": block})
}
