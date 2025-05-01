package middleware

import (
	"context"
	"net/http"
	"strings"

	"ingnyusan/v2/internal/i18n"
)

// Language context key
type langContextKey string

const (
	// LanguageKey is the key to store/retrieve language from context
	LanguageKey langContextKey = "language"
	// LanguageCookie is the name of the cookie storing language preference
	LanguageCookie = "lang"
)

// I18nMiddleware adds language context to requests
func I18nMiddleware(translator *i18n.Translator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Determine language from (in order of priority):
			// 1. Query parameter "lang"
			// 2. Cookie
			// 3. Accept-Language header
			// 4. Default to English

			// Check query parameter
			lang := r.URL.Query().Get("lang")

			// If not valid or not specified, check cookie
			if lang == "" || !translator.IsValidLanguage(lang) {
				if cookie, err := r.Cookie(LanguageCookie); err == nil {
					lang = cookie.Value
				}
			}

			// If still not valid, check Accept-Language header
			if lang == "" || !translator.IsValidLanguage(lang) {
				acceptLang := r.Header.Get("Accept-Language")
				preferredLangs := parseAcceptLanguage(acceptLang)

				// Find first supported language
				for _, l := range preferredLangs {
					if translator.IsValidLanguage(l) {
						lang = l
						break
					}
				}
			}

			// If still not set or not valid, use default
			if lang == "" || !translator.IsValidLanguage(lang) {
				lang = translator.GetDefaultLanguage()
			}

			// If language comes from query param, set cookie for future requests
			if r.URL.Query().Get("lang") != "" && translator.IsValidLanguage(r.URL.Query().Get("lang")) {
				cookie := &http.Cookie{
					Name:     LanguageCookie,
					Value:    lang,
					Path:     "/",
					MaxAge:   365 * 24 * 60 * 60, // 1 year
					HttpOnly: true,
					SameSite: http.SameSiteLaxMode,
				}
				http.SetCookie(w, cookie)
			}

			// Add language to request context
			ctx := context.WithValue(r.Context(), LanguageKey, lang)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetLanguage retrieves the language code from the request context
func GetLanguage(r *http.Request) string {
	if lang, ok := r.Context().Value(LanguageKey).(string); ok {
		return lang
	}
	return i18n.LangEnglish // Default to English
}

// parseAcceptLanguage parses the Accept-Language header and returns
// language codes in order of preference
func parseAcceptLanguage(header string) []string {
	languages := []string{}

	// Split by comma
	parts := strings.Split(header, ",")
	for _, part := range parts {
		// Split by semicolon to separate language from quality
		langPart := strings.Split(strings.TrimSpace(part), ";")[0]

		// Get just the language code (first two chars)
		if len(langPart) >= 2 {
			langCode := strings.ToLower(langPart[:2])
			languages = append(languages, langCode)
		}
	}

	return languages
}
