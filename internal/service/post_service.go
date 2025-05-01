package service

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"blog/internal/parser"
)

// PostService handles operations related to blog posts
type PostService struct {
	Parser        *parser.MarkdownParser
	PostsDir      string
	posts         []parser.Post
	postsBySlug   map[string]parser.Post
	postsByTag    map[string][]parser.Post
}

// NewPostService creates a new post service
func NewPostService(postsDir string) *PostService {
	return &PostService{
		Parser:      parser.NewMarkdownParser(),
		PostsDir:    postsDir,
		postsBySlug: make(map[string]parser.Post),
		postsByTag:  make(map[string][]parser.Post),
	}
}

// LoadPosts loads all posts from the posts directory
func (ps *PostService) LoadPosts() error {
	// Clear existing posts
	ps.posts = []parser.Post{}
	ps.postsBySlug = make(map[string]parser.Post)
	ps.postsByTag = make(map[string][]parser.Post)

	// Walk through the posts directory
	err := filepath.Walk(ps.PostsDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-markdown files
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		// Read the file
		content, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Error reading file %s: %v", path, err)
			return nil
		}

		// Get the slug from the filename
		slug := strings.TrimSuffix(filepath.Base(path), ".md")

		// Parse the post
		post, err := ps.Parser.ParsePost(string(content), slug)
		if err != nil {
			log.Printf("Error parsing post %s: %v", path, err)
			return nil
		}

		// Add the post to the collections
		ps.posts = append(ps.posts, post)
		ps.postsBySlug[post.Slug] = post

		// Index posts by tag
		for _, tag := range post.Tags {
			ps.postsByTag[tag] = append(ps.postsByTag[tag], post)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Sort posts by date (newest first)
	sort.Slice(ps.posts, func(i, j int) bool {
		return ps.posts[i].Date.After(ps.posts[j].Date)
	})

	// Sort posts in each tag
	for tag, posts := range ps.postsByTag {
		sort.Slice(posts, func(i, j int) bool {
			return posts[i].Date.After(posts[j].Date)
		})
		ps.postsByTag[tag] = posts
	}

	return nil
}

// GetAllPosts returns all posts sorted by date (newest first)
func (ps *PostService) GetAllPosts() []parser.Post {
	return ps.posts
}

// GetPostBySlug returns a post by its slug
func (ps *PostService) GetPostBySlug(slug string) (parser.Post, error) {
	post, ok := ps.postsBySlug[slug]
	if !ok {
		return parser.Post{}, errors.New("post not found")
	}
	return post, nil
}

// GetPostsByTag returns posts with a specific tag
func (ps *PostService) GetPostsByTag(tag string) []parser.Post {
	return ps.postsByTag[tag]
}

// GetRecentPosts returns the n most recent posts
func (ps *PostService) GetRecentPosts(n int) []parser.Post {
	if n > len(ps.posts) {
		n = len(ps.posts)
	}
	return ps.posts[:n]
}

// GetAllTags returns all unique tags
func (ps *PostService) GetAllTags() []string {
	tags := make([]string, 0, len(ps.postsByTag))
	for tag := range ps.postsByTag {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

// GenerateRSS generates an RSS feed XML for the blog posts
func (ps *PostService) GenerateRSS(title, description, link string) string {
	rss := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>` + title + `</title>
    <description>` + description + `</description>
    <link>` + link + `</link>
    <atom:link href="` + link + `/rss" rel="self" type="application/rss+xml" />`

	for _, post := range ps.GetRecentPosts(10) {
		pubDate := post.Date.Format("Mon, 02 Jan 2006 15:04:05 -0700")
		postLink := link + "/posts/" + post.Slug
		
		rss += `
    <item>
      <title>` + post.Title + `</title>
      <description>` + post.Description + `</description>
      <link>` + postLink + `</link>
      <guid>` + postLink + `</guid>
      <pubDate>` + pubDate + `</pubDate>
    </item>`
	}

	rss += `
  </channel>
</rss>`

	return rss
} 