package utils

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"
	"unicode"

	"github.com/Open-Code-Zone/cms/internal/database"
	"gopkg.in/yaml.v2"
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

// ParseMarkdown takes a raw markdown string, separates the frontmatter from the content, and parses them.
func ExtractFrontMatter(fileContentWithFrontMatter string) (string, string) {
	// Define regex pattern for front matter block
	frontMatterPattern := regexp.MustCompile(`(?s)^---\n(.*?)\n---\n`)
	frontMatter := make(map[string]interface{})

	// Find front matter block
	match := frontMatterPattern.FindStringSubmatch(fileContentWithFrontMatter)
	if len(match) < 2 {
		// No front matter found, return original content and empty JSON object
		emptyFrontMatterJSON, _ := json.Marshal(frontMatter)
		return fileContentWithFrontMatter, string(emptyFrontMatterJSON)
	}

	// Extract front matter and remaining content
	frontMatterYaml := match[1]
	content := strings.TrimPrefix(fileContentWithFrontMatter, match[0])

	// Parse YAML front matter into map
	if err := yaml.Unmarshal([]byte(frontMatterYaml), &frontMatter); err != nil {
		log.Println("Error parsing front matter:", err)
		return content, "{}" // Return empty JSON object on error
	}

	// Convert front matter to JSON string
	frontMatterJSON, err := json.Marshal(frontMatter)
	if err != nil {
		log.Println("Error converting front matter to JSON:", err)
		return content, "{}" // Return empty JSON object on error
	}

	return content, string(frontMatterJSON)
}

func GenerateMarkdownFile(collectionItem database.GetCollectionItemRow) string {
	var markdownFile strings.Builder

	yamlMetadata, err := convertJSONToYAML(collectionItem.Metadata)
	if err != nil {
		log.Println("Error converting JSON to YAML:", err)
		return ""
	}

	// Write front matter
	markdownFile.WriteString("---\n")
	markdownFile.WriteString(string(yamlMetadata))
	markdownFile.WriteString("\n---\n")

	// Write content
	markdownFile.WriteString(collectionItem.Content)

	return markdownFile.String()
}

func convertJSONToYAML(jsonStr string) ([]byte, error) {
	// Step 1: Unmarshal the JSON string into a Go map
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, err
	}

	// Step 2: Marshal the Go map to YAML
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}

	return yamlData, nil
}
