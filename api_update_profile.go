package zca

import (
	"context"
	"net/http"
)

type UpdateProfilePayload struct {
	Profile UpdateProfileInfo
	// Biz: for Business Accounts include Cate, otherwise the category is removed.
	// Empty fields are left out (unchanged).
	Biz *UpdateProfileBiz
}

type UpdateProfileInfo struct {
	Name   string
	Dob    string // YYYY-MM-DD
	Gender Gender
}

type UpdateProfileBiz struct {
	Cate        *BusinessCategory
	Description string
	Address     string
	Website     string
	Email       string
}

// UpdateProfile changes account profile information.
func (a *API) UpdateProfile(ctx context.Context, payload UpdateProfilePayload) (string, error) {
	biz := map[string]any{}
	if b := payload.Biz; b != nil {
		for k, v := range map[string]string{"desc": b.Description, "addr": b.Address, "website": b.Website, "email": b.Email} {
			if v != "" {
				biz[k] = v
			}
		}
		if b.Cate != nil {
			biz["cate"] = *b.Cate
		}
	}
	params := map[string]any{
		"profile": mustJSON(map[string]any{
			"name": payload.Profile.Name, "dob": payload.Profile.Dob, "gender": payload.Profile.Gender,
		}),
		"biz":      mustJSON(biz),
		"language": a.Language,
	}
	return call[string](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("profile")+"/api/social/profile/update", nil, true), params)
}
