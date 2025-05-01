package templates

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"ingnyusan/v2/internal/config"
)

// Global config reference for templates
var globalConfig *config.Config

// SetConfig sets the global config reference for templates
func SetConfig(cfg *config.Config) {
	globalConfig = cfg
}

// currentYear returns the current year as a string
func currentYear() string {
	return fmt.Sprintf("%d", time.Now().Year())
}

// formatDate formats a time.Time as a readable date string
func formatDate(t time.Time) string {
	return t.Format("January 2, 2006")
}

// getThemeFromRequest extracts theme preference from the request
func getThemeFromRequest(r *http.Request) string {
	// Check cookie first
	if cookie, err := r.Cookie("theme"); err == nil {
		theme := strings.ToLower(cookie.Value)
		if theme == "dark" || theme == "light" {
			return theme
		}
	}

	// Otherwise use default from config
	return getConfigTheme()
}

// getConfigTheme returns the default theme from config
func getConfigTheme() string {
	if globalConfig != nil && globalConfig.Theme.DefaultTheme != "" {
		return globalConfig.Theme.DefaultTheme
	}
	return "light" // Fallback default
}
