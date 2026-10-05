package zca

import (
	"context"
	"net/http"
	"strconv"
)

type ReportReason int

const (
	ReportReasonSensitive ReportReason = 1
	ReportReasonAnnoy     ReportReason = 2
	ReportReasonFraud     ReportReason = 3
	ReportReasonOther     ReportReason = 0
)

// SendReportOptions: Content is used only with ReportReasonOther.
type SendReportOptions struct {
	Reason  ReportReason
	Content string
}

type SendReportResponse struct {
	ReportID string `json:"reportId"`
}

// SendReport reports a user or group to Zalo.
func (a *API) SendReport(ctx context.Context, options SendReportOptions, threadID string, threadType ThreadType) (*SendReportResponse, error) {
	var params map[string]any
	var u string
	if threadType == ThreadTypeUser {
		params = map[string]any{"idTo": threadID, "objId": "person.profile", "reason": strconv.Itoa(int(options.Reason))}
		if options.Reason == ReportReasonOther {
			params["content"] = options.Content
		}
		u = a.svc("profile") + "/api/report/abuse-v2"
	} else {
		content := ""
		if options.Reason == ReportReasonOther {
			content = options.Content
		}
		params = map[string]any{"uidTo": threadID, "type": 14, "reason": options.Reason, "content": content, "imei": a.IMEI}
		u = a.svc("profile") + "/api/social/profile/reportabuse"
	}
	return call[*SendReportResponse](ctx, a.Session, http.MethodPost, a.MakeURL(u, nil, true), params)
}
