package wiki

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
)

// ContentItem represents a single parsed item from the MediaWiki text.
type ContentItem struct {
	Type  string `json:"type"`
	Level int    `json:"level,omitempty"` // For headings only
	Text  string `json:"text,omitempty"`  // For headings, bold, italic, and paragraph
	Title string `json:"title,omitempty"` // For internal links
	URL   string `json:"url,omitempty"`   // For external links
}

// ParseMediaWiki parses MediaWiki-like syntax into JSON-structured content.
func ParseMediaWiki(text string) ([]ContentItem, error) {
	var content []ContentItem

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "="): // Headings
			content = append(content, parseHeading(line))
		case strings.HasPrefix(line, "'''"): // Bold
			content = append(content, ContentItem{Type: "bold", Text: parseBold(line)})
		case strings.HasPrefix(line, "''"): // Italic
			content = append(content, ContentItem{Type: "italic", Text: parseItalic(line)})
		case strings.HasPrefix(line, "[["): // Internal Links
			content = append(content, parseInternalLink(line))
		case strings.HasPrefix(line, "["): // External Links
			content = append(content, parseExternalLink(line))
		default: // Regular text (paragraph)
			content = append(content, ContentItem{Type: "paragraph", Text: line})
		}
	}

	return content, nil
}

// parseHeading parses MediaWiki-style headings, e.g., `== Heading ==`
func parseHeading(line string) ContentItem {
	headingRegex := regexp.MustCompile(`^(=+)\s*(.*?)\s*=+$`)
	matches := headingRegex.FindStringSubmatch(line)
	level := len(matches[1]) // Number of `=` symbols represents heading level
	return ContentItem{Type: "heading", Level: level, Text: matches[2]}
}

// parseBold parses `”'bold”'` text
func parseBold(line string) string {
	boldRegex := regexp.MustCompile(`'''(.*?)'''`)
	return boldRegex.ReplaceAllString(line, "$1")
}

// parseItalic parses `”italic”` text
func parseItalic(line string) string {
	italicRegex := regexp.MustCompile(`''(.*?)''`)
	return italicRegex.ReplaceAllString(line, "$1")
}

// parseInternalLink parses `[[PageName]]` links to internal pages
func parseInternalLink(line string) ContentItem {
	internalLinkRegex := regexp.MustCompile(`\[\[([^\]]+?)\]\]`)
	matches := internalLinkRegex.FindStringSubmatch(line)
	title := matches[1]
	return ContentItem{Type: "link", Title: title, URL: fmt.Sprintf("/page/%s", title)}
}

// parseExternalLink parses `[http://example.com]` links
func parseExternalLink(line string) ContentItem {
	externalLinkRegex := regexp.MustCompile(`\[(http[^\s]+)\]`)
	matches := externalLinkRegex.FindStringSubmatch(line)
	url := matches[1]
	return ContentItem{Type: "link", URL: url, Title: url}
}

func test() {
	// Example MediaWiki content
	content := `= Welcome to the Wiki =
This is a simple page about ''Golang''. Visit the '''Golang Page''' by clicking [[Golang]].
To learn more, visit [https://golang.org].`

	// Parse the content to JSON structure
	parsedContent, err := ParseMediaWiki(content)
	if err != nil {
		log.Fatal("Error parsing content:", err)
	}

	// Convert parsed content to JSON
	jsonOutput, err := json.MarshalIndent(parsedContent, "", "  ")
	if err != nil {
		log.Fatal("Error generating JSON:", err)
	}

	// Display JSON output
	fmt.Println(string(jsonOutput))
}
