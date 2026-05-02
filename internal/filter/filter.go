package filter

import "strings"

// Matcher applies a display-only substring filter to rendered log lines.
type Matcher struct {
	query string
}

func NewMatcher(query string) Matcher {
	return Matcher{query: query}
}

func (m Matcher) Matches(line string) bool {
	return m.query == "" || strings.Contains(line, m.query)
}

func Lines(lines []string, query string) []string {
	matcher := NewMatcher(query)
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		if matcher.Matches(line) {
			filtered = append(filtered, line)
		}
	}
	return filtered
}
