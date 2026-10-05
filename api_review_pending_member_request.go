package zca

import (
	"context"
	"net/http"
)

type ReviewPendingMemberRequestPayload struct {
	Members   []string
	IsApprove bool
}

type ReviewPendingMemberRequestStatus int

const (
	ReviewPendingMemberRequestStatusSuccess                ReviewPendingMemberRequestStatus = 0
	ReviewPendingMemberRequestStatusNotInPendingList       ReviewPendingMemberRequestStatus = 170
	ReviewPendingMemberRequestStatusAlreadyInGroup         ReviewPendingMemberRequestStatus = 178
	ReviewPendingMemberRequestStatusInsufficientPermission ReviewPendingMemberRequestStatus = 166
)

// ReviewPendingMemberRequestResponse maps member id -> status.
type ReviewPendingMemberRequestResponse map[string]ReviewPendingMemberRequestStatus

// ReviewPendingMemberRequest approves or rejects pending member requests (group leader/deputy only).
func (a *API) ReviewPendingMemberRequest(ctx context.Context, payload ReviewPendingMemberRequestPayload, groupID string) (ReviewPendingMemberRequestResponse, error) {
	approve := 0
	if payload.IsApprove {
		approve = 1
	}
	params := map[string]any{"grid": groupID, "members": payload.Members, "isApprove": approve}
	return call[ReviewPendingMemberRequestResponse](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("group")+"/api/group/pending-mems/review", nil, true), params)
}
