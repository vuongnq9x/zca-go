package zca

import (
	"context"
	"net/http"
)

type GetBizAccountResponse struct {
	Biz *struct {
		Desc    *string          `json:"desc"`
		Cate    BusinessCategory `json:"cate"`
		Addr    string           `json:"addr"`
		Website string           `json:"website"`
		Email   string           `json:"email"`
	} `json:"biz,omitempty"`
	SettingStartPage *struct {
		EnableBizLabel int64 `json:"enable_biz_label"`
		EnableCate     int64 `json:"enable_cate"`
		EnableAdd      int64 `json:"enable_add"`
		CtaProfile     int64 `json:"cta_profile"`
		// CtaCatalog is a relative path: https://catalog.zalo.me/${cta_catalog}
		CtaCatalog *string `json:"cta_catalog"`
	} `json:"setting_start_page,omitempty"`
	PkgID int64 `json:"pkgId"`
}

// GetBizAccount gets the business account info of a friend.
func (a *API) GetBizAccount(ctx context.Context, friendID string) (*GetBizAccountResponse, error) {
	return call[*GetBizAccountResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("profile")+"/api/social/friend/get-bizacc", nil, true), map[string]any{"fid": friendID})
}
