package templates

import (
	"fmt"
	"time"
)

// currentYear returns the current year as a string
func currentYear() string {
	return fmt.Sprintf("%d", time.Now().Year())
}

// formatDate formats a time.Time as a readable date string
func formatDate(t time.Time) string {
	return t.Format("January 2, 2006")
} 