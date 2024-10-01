package types

import "time"

type MarkdownFile struct {
	FileName    string
	LastUpdated time.Time
}

type BlogMetadata struct {
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Date        string   `yaml:"date"`
	Authors     []string `yaml:"authors"`
	Image       string   `yaml:"image"`
	Tags        []string `yaml:"tags"`
	Featured    bool     `yaml:"featured"`
}

// BlogPost represents the parsed content of a blog post with its metadata.
type BlogPost struct {
	Metadata BlogMetadata
	Content  string
}
