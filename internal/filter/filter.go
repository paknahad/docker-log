package filter

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// State describes how a display-only filter should match rendered log content.
type State struct {
	Text          string
	Regex         bool
	CaseSensitive bool
}

// Matcher applies a display-only filter to rendered log lines.
type Matcher struct {
	state State
	regex *regexp.Regexp
	err   error
}

func NewMatcher(query string) Matcher {
	matcher, _ := NewMatcherForState(NewState(query))
	return matcher
}

func NewState(text string) State {
	return State{Text: text, CaseSensitive: true}
}

func NewMatcherForState(state State) (Matcher, error) {
	matcher := Matcher{state: state}
	if state.Text == "" || !state.Regex {
		return matcher, nil
	}

	pattern := state.Text
	if !state.CaseSensitive {
		pattern = "(?i:" + pattern + ")"
	}
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		validationErr := ValidationError{Pattern: state.Text, Err: err}
		matcher.err = validationErr
		return matcher, validationErr
	}
	matcher.regex = compiled
	return matcher, nil
}

func (m Matcher) Matches(line string) bool {
	if m.err != nil {
		return false
	}
	if m.state.Text == "" {
		return true
	}
	if m.regex != nil {
		return m.regex.MatchString(line)
	}
	if !m.state.CaseSensitive {
		return strings.Contains(strings.ToLower(line), strings.ToLower(m.state.Text))
	}
	return strings.Contains(line, m.state.Text)
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

func LinesWithState(lines []string, state State) ([]string, error) {
	matcher, err := NewMatcherForState(state)
	if err != nil {
		return nil, err
	}

	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		if matcher.Matches(line) {
			filtered = append(filtered, line)
		}
	}
	return filtered, nil
}

// ValidationError reports filter input that cannot be used for the selected mode.
type ValidationError struct {
	Pattern string
	Err     error
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("invalid regex pattern %q: %v", e.Pattern, e.Err)
}

func (e ValidationError) Unwrap() error {
	return e.Err
}

func IsValidationError(err error) bool {
	var validationErr ValidationError
	return errors.As(err, &validationErr)
}
