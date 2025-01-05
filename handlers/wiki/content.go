package wiki

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

// ContentItem represents a single parsed item from the MediaWiki text.
type ContentItem struct {
	Type     string        `json:"type"`
	Level    int           `json:"level,omitempty"`    // For headings only
	Text     string        `json:"text,omitempty"`     // For headings, bold, italic, and paragraph
	Title    string        `json:"title,omitempty"`    // For internal links
	URL      string        `json:"url,omitempty"`      // For external links
	Children []ContentItem `json:"children,omitempty"` // Nested content for headings
}

func (c ContentItem) DisplayTitle() string {
	return fmt.Sprintf("%s \n", c.Text)
}

// ParseMediaWiki parses MediaWiki-like syntax into JSON-structured content.
func ParseMediaWiki(text string) ([]ContentItem, error) {
	var content []ContentItem
	var headingStack []*ContentItem

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "="): // Headings
			heading := parseHeading(line)
			content, headingStack = addHeading(content, headingStack, heading)
		case strings.HasPrefix(line, "'''"): // Bold
			content = append(content, ContentItem{Type: "bold", Text: parseBold(line)})
		case strings.HasPrefix(line, "''"): // Italic
			content = append(content, ContentItem{Type: "italic", Text: parseItalic(line)})
		case strings.HasPrefix(line, "[["): // Internal Links
			content = append(content, parseInternalLink(line))
		case strings.HasPrefix(line, "["): // External Links
			content = append(content, parseExternalLink(line))
		default: // Regular text (paragraph)
			if len(headingStack) > 0 {
				// Add paragraphs under the last heading
				lastHeading := headingStack[len(headingStack)-1]
				lastHeading.Children = append(lastHeading.Children, ContentItem{Type: "paragraph", Text: line})
			} else {
				content = append(content, ContentItem{Type: "paragraph", Text: line})
			}
		}
	}

	return content, nil
}

// parseHeading parses MediaWiki-style headings, e.g., `== Heading ==`
func parseHeading(line string) ContentItem {
	headingRegex := regexp.MustCompile(`^(=+) *(.*?)\s*=+$`)
	matches := headingRegex.FindStringSubmatch(line)
	level := len(matches[1]) // Number of `=` symbols represents heading level
	return ContentItem{Type: "heading", Level: level, Text: matches[2]}
}

// addHeading adds a heading to the content tree based on its level.
func addHeading(content []ContentItem, stack []*ContentItem, heading ContentItem) ([]ContentItem, []*ContentItem) {
	// Remove headings from the stack that are deeper than the current level
	for len(stack) > 0 && stack[len(stack)-1].Level >= heading.Level {
		stack = stack[:len(stack)-1]
	}

	if len(stack) == 0 {
		// Top-level heading
		content = append(content, heading)
		stack = append(stack, &content[len(content)-1])
	} else {
		// Add as a child to the last heading in the stack
		parent := stack[len(stack)-1]
		parent.Children = append(parent.Children, heading)
		stack = append(stack, &parent.Children[len(parent.Children)-1])
	}

	return content, stack
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
	if matches == nil || len(matches) < 2 {
		// Если не удалось найти совпадение по регулярке, возвращаем описание ошибки
		log.Printf("Failed to parse external link: %s", line)
		return ContentItem{Type: "error", Text: "Invalid external link format"}
	}
	url := matches[1]
	return ContentItem{Type: "link", URL: url, Title: url}
}

func ConvertStringWikiToJSON(title string) []ContentItem {
	// Parse the content to JSON structure
	parsedContent, err := ParseMediaWiki(title)
	if err != nil {
		log.Fatal("Error parsing content:", err)
	}
	return parsedContent
}
