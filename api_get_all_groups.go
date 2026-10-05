package zca

import (
	"context"
	"net/http"
)

type GetAllGroupsResponse struct {
	Version string `json:"version"`
	// GridVerMap maps group id -> version.
	GridVerMap map[string]string `json:"gridVerMap"`
}

// GetAllGroups gets all group ids with their versions.
func (a *API) GetAllGroups(ctx context.Context) (*GetAllGroupsResponse, error) {
	resp, err := a.Request(ctx, http.MethodGet, a.MakeURL(a.svc("group_poll")+"/api/group/getlg/v4", nil, true), nil, nil)
	if err != nil {
		return nil, err
	}
	raw, err := a.Resolve(resp, true)
	if err != nil {
		return nil, err
	}
	return decodeData[*GetAllGroupsResponse](raw)
}
