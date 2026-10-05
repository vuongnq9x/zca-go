package zca

import (
	"context"
	"net/http"
)

type UpdateLangAvailableLanguages string

const (
	UpdateLangAvailableLanguagesVI UpdateLangAvailableLanguages = "VI"
	UpdateLangAvailableLanguagesEN UpdateLangAvailableLanguages = "EN"
)

// UpdateLang updates the account language ("" means VI). Calling it alone does not change the user's language.
func (a *API) UpdateLang(ctx context.Context, language UpdateLangAvailableLanguages) (string, error) {
	if language == "" {
		language = UpdateLangAvailableLanguagesVI
	}
	return call[string](ctx, a.Session, http.MethodGet,
		a.MakeURL(a.svc("profile")+"/api/social/profile/updatelang", nil, true), map[string]any{"language": language})
}
