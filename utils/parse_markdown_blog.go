package utils

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v2"
	"time"
	"unicode"

	"github.com/Open-Code-Zone/cms/types"
)

func sanitizeFilename(title string) string {
	title = strings.ToLower(title)
	title = strings.ReplaceAll(title, " ", "-")

	// Remove any character that isn't alphanumeric or hyphen
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '-' {
			return r
		}
		return -1
	}, title)
}

func createFrontMatter(blogPost *types.BlogPost) string {
	return fmt.Sprintf(`---
title: "%s"
description: "%s"
date: %s
authors: [%s]
image: "%s"
tags: [%s]
---`,
		blogPost.Metadata.Title,
		blogPost.Metadata.Description,
		time.Now().Format("2006-01-02"),
		strings.Join(blogPost.Metadata.Authors, ", "),
		blogPost.Metadata.Image,
		strings.Join(blogPost.Metadata.Tags, ", "))
}

// ParseMarkdown takes a raw markdown string, separates the frontmatter from the content, and parses them.
func ParseMarkdown(rawMarkdown string) (*types.BlogPost, error) {
	// Split the frontmatter and content based on the "---" delimiter.
	parts := strings.SplitN(rawMarkdown, "---", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid markdown format")
	}

	// Parse frontmatter (YAML metadata).
	var metadata types.BlogMetadata
	err := yaml.Unmarshal([]byte(parts[1]), &metadata)
	if err != nil {
		return nil, fmt.Errorf("error parsing frontmatter: %v", err)
	}

	// Parse the main content markdown.
	var buf bytes.Buffer
	md := goldmark.New()
	if err := md.Convert([]byte(parts[2]), &buf); err != nil {
		return nil, fmt.Errorf("error parsing markdown content: %v", err)
	}

	return &types.BlogPost{
		Metadata: metadata,
		Content:  buf.String(),
	}, nil
}

func GenerateMarkdown(metadata types.BlogMetadata, content string) (string, error) {
	frontmatter, err := yaml.Marshal(&metadata)
	if err != nil {
		return "", fmt.Errorf("error generating frontmatter: %v", err)
	}

	markdown := fmt.Sprintf(`---
%s---
%s`, string(frontmatter), content)

	return markdown, nil
}
