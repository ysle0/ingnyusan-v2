package parser

import (
	"bytes"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// Post represents a blog post with metadata
type Post struct {
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	RawContent  string    `json:"rawContent"`
	Date        time.Time `json:"date"`
	Tags        []string  `json:"tags"`
	TimeToRead  int       `json:"timeToRead"` // in minutes
}

// MarkdownParser handles the parsing of markdown content
type MarkdownParser struct {
	md goldmark.Markdown
}

// NewMarkdownParser creates a new markdown parser
func NewMarkdownParser() *MarkdownParser {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	return &MarkdownParser{md}
}

// ParseMarkdown converts markdown to HTML
func (mp *MarkdownParser) ParseMarkdown(content string) (string, error) {
	var buf bytes.Buffer
	if err := mp.md.Convert([]byte(content), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ParsePost parses a markdown file with frontmatter
func (mp *MarkdownParser) ParsePost(content string, slug string) (Post, error) {
	// Extract frontmatter and content
	frontmatter, mdContent := extractFrontmatter(content)

	// Convert markdown to HTML
	htmlContent, err := mp.ParseMarkdown(mdContent)
	if err != nil {
		return Post{}, err
	}

	// Parse frontmatter
	title := extractFrontmatterValue(frontmatter, "title")
	description := extractFrontmatterValue(frontmatter, "description")
	dateStr := extractFrontmatterValue(frontmatter, "date")
	tagsStr := extractFrontmatterValue(frontmatter, "tags")
	language := extractFrontmatterValue(frontmatter, "language")

	// Parse date
	date, err := time.Parse("YYYY-MM-dd", dateStr)
	if err != nil {
		// Use current time if date is invalid
		date = time.Now()
	}

	// Parse tags
	tags := []string{}
	if tagsStr != "" {
		for _, tag := range strings.Split(tagsStr, ",") {
			tags = append(tags, strings.TrimSpace(tag))
		}
	}

	// Calculate time to read based on language
	// Default reading speeds:
	// - English: ~200-250 words per minute
	// - Korean: ~150-180 chars per minute (different measurement due to language structure)
	words := len(strings.Fields(mdContent))
	timeToRead := 0

	if strings.ToLower(language) == "korean" || strings.ToLower(language) == "ko" {
		// For Korean: count characters instead of words
		// Average reading speed ~150-180 characters per minute for native speakers
		chars := len(strings.ReplaceAll(mdContent, " ", ""))
		timeToRead = chars / 180
	} else {
		// For English and other languages: use word count
		// Average reading speed ~200-250 words per minute
		timeToRead = words / 200
	}

	if timeToRead < 1 {
		timeToRead = 1
	}

	return Post{
		Slug:        slug,
		Title:       title,
		Description: description,
		Content:     htmlContent,
		RawContent:  mdContent,
		Date:        date,
		Tags:        tags,
		TimeToRead:  timeToRead,
	}, nil
}

// extractFrontmatter separates frontmatter from markdown content
func extractFrontmatter(content string) (string, string) {
	const frontmatterDelimiter = "---"

	// Check if content has frontmatter
	lines := strings.Split(content, "\n")
	if len(lines) < 2 || lines[0] != frontmatterDelimiter {
		return "", content
	}

	// Find the end of frontmatter
	endIndex := -1
	for i := 1; i < len(lines); i++ {
		if lines[i] == frontmatterDelimiter {
			endIndex = i
			break
		}
	}

	if endIndex == -1 {
		return "", content
	}

	frontmatter := strings.Join(lines[1:endIndex], "\n")
	markdown := strings.Join(lines[endIndex+1:], "\n")

	return frontmatter, markdown
}

// extractFrontmatterValue extracts a value from frontmatter by key
func extractFrontmatterValue(frontmatter string, key string) string {
	lines := strings.Split(frontmatter, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		k := strings.TrimSpace(parts[0])
		if k == key {
			return strings.TrimSpace(parts[1])
		}
	}
	return ""
}
