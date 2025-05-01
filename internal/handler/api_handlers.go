package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

// RegisterAPI registers all API routes with the provided router
func (h *Handler) RegisterAPI(r chi.Router) {
	r.Post("/api/preview", h.PreviewHandler)
	r.Post("/api/save-post", h.SavePostHandler)
	r.Get("/api/github-profile", h.GitHubProfileHandler)
	r.Get("/api/toggle-theme", h.ToggleThemeHandler)
}

// PreviewHandler handles rendering markdown preview
func (h *Handler) PreviewHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	// Get the content
	content := r.FormValue("content")

	// Parse the markdown
	htmlContent, err := h.PostService.Parser.ParseMarkdown(content)
	if err != nil {
		http.Error(w, "Could not parse markdown", http.StatusInternalServerError)
		return
	}

	// Write the HTML
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlContent))
}

// SavePostHandler handles saving a post
func (h *Handler) SavePostHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	// Get the form values
	title := r.FormValue("title")
	description := r.FormValue("description")
	tagsStr := r.FormValue("tags")
	slug := r.FormValue("slug")
	content := r.FormValue("content")

	// Validate
	if title == "" || description == "" || slug == "" || content == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Format the tags
	tags := []string{}
	if tagsStr != "" {
		for tag := range strings.SplitSeq(tagsStr, ",") {
			tags = append(tags, strings.TrimSpace(tag))
		}
	}

	// Create the frontmatter
	frontmatter := fmt.Sprintf("---\ntitle: %s\ndescription: %s\ndate: %s\ntags: %s\n---\n\n",
		title,
		description,
		time.Now().Format("2006-01-02"),
		tagsStr,
	)

	// Combine
	fullContent := frontmatter + content

	// Save to file
	filePath := filepath.Join(h.PostService.PostsDir, slug+".md")
	err = os.WriteFile(filePath, []byte(fullContent), 0644)
	if err != nil {
		http.Error(w, "Could not save post", http.StatusInternalServerError)
		return
	}

	// Reload posts
	err = h.PostService.LoadPosts()
	if err != nil {
		http.Error(w, "Could not reload posts", http.StatusInternalServerError)
		return
	}

	// Redirect to the post
	http.Redirect(w, r, "/posts/"+slug, http.StatusSeeOther)
}

// GitHubProfileHandler fetches and returns GitHub profile information
func (h *Handler) GitHubProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch GitHub profile (simplified for demo)
	resp, err := http.Get("https://api.github.com/users/ysle0")
	if err != nil {
		http.Error(w, "Could not fetch GitHub profile", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Could not read GitHub profile", http.StatusInternalServerError)
		return
	}

	// Parse the JSON
	var profile map[string]interface{}
	err = json.Unmarshal(body, &profile)
	if err != nil {
		http.Error(w, "Could not parse GitHub profile", http.StatusInternalServerError)
		return
	}

	// Create a simple HTML response
	html := fmt.Sprintf(`
	<div class="github-profile">
		<img src="%s" alt="Profile Image" class="profile-image" />
		<h3>%s</h3>
		<p>%s</p>
		<p>Repositories: %v</p>
		<p>Followers: %v</p>
		<p>Following: %v</p>
		<p><a href="%s" target="_blank">View on GitHub</a></p>
	</div>
	`,
		fmt.Sprintf("%v", profile["avatar_url"]),
		fmt.Sprintf("%v", profile["name"]),
		fmt.Sprintf("%v", profile["bio"]),
		profile["public_repos"],
		profile["followers"],
		profile["following"],
		fmt.Sprintf("%v", profile["html_url"]),
	)

	// Write the HTML
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// ToggleThemeHandler toggles between light and dark theme
func (h *Handler) ToggleThemeHandler(w http.ResponseWriter, r *http.Request) {
	// Get current theme from cookie
	var newTheme string
	if cookie, err := r.Cookie("theme"); err == nil {
		// Toggle theme
		if cookie.Value == "dark" {
			newTheme = "light"
		} else {
			newTheme = "dark"
		}
	} else {
		// No cookie found, get default from config and toggle it
		if h.Config.Theme.DefaultTheme == "dark" {
			newTheme = "light"
		} else {
			newTheme = "dark"
		}
	}

	// Set cookie with new theme
	http.SetCookie(w, &http.Cookie{
		Name:     "theme",
		Value:    newTheme,
		Path:     "/",
		MaxAge:   365 * 24 * 60 * 60, // 1 year
		HttpOnly: false,              // Allow JavaScript access
		SameSite: http.SameSiteLaxMode,
	})

	// Return the new theme
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"theme": newTheme})
}
