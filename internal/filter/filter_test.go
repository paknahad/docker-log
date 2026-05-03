package filter

import "testing"

func TestMatcherMatchesAllLinesWhenQueryEmpty(t *testing.T) {
	matcher := NewMatcher("")

	for _, line := range []string{"api: started", "", "worker: ready"} {
		if !matcher.Matches(line) {
			t.Fatalf("Matches(%q) = false, want true", line)
		}
	}
}

func TestMatcherIsCaseSensitive(t *testing.T) {
	matcher := NewMatcher("Error")

	if !matcher.Matches("api: Error opening stream") {
		t.Fatal("Matches() = false, want true for exact-case query")
	}
	if matcher.Matches("api: error opening stream") {
		t.Fatal("Matches() = true, want false for different-case query")
	}
}

func TestStateDefaultsToPlainTextCaseSensitiveMatching(t *testing.T) {
	state := NewState("Error")

	if state.Text != "Error" {
		t.Fatalf("Text = %q, want Error", state.Text)
	}
	if state.Regex {
		t.Fatal("Regex = true, want false")
	}
	if !state.CaseSensitive {
		t.Fatal("CaseSensitive = false, want true")
	}

	matcher, err := NewMatcherForState(state)
	if err != nil {
		t.Fatalf("NewMatcherForState() error = %v, want nil", err)
	}
	if !matcher.Matches("api: Error opening stream") {
		t.Fatal("Matches() = false, want true for exact-case query")
	}
	if matcher.Matches("api: error opening stream") {
		t.Fatal("Matches() = true, want false for different-case query")
	}
}

func TestMatcherSupportsPlainTextCaseInsensitiveMatching(t *testing.T) {
	matcher, err := NewMatcherForState(State{Text: "Error", CaseSensitive: false})
	if err != nil {
		t.Fatalf("NewMatcherForState() error = %v, want nil", err)
	}

	if !matcher.Matches("api: error opening stream") {
		t.Fatal("Matches() = false, want true for different-case plain-text query")
	}
}

func TestMatcherSupportsRegexCaseSensitiveMatching(t *testing.T) {
	matcher, err := NewMatcherForState(State{Text: `err(or)?`, Regex: true, CaseSensitive: true})
	if err != nil {
		t.Fatalf("NewMatcherForState() error = %v, want nil", err)
	}

	if !matcher.Matches("api: error opening stream") {
		t.Fatal("Matches() = false, want true for regex match")
	}
	if matcher.Matches("api: ERROR opening stream") {
		t.Fatal("Matches() = true, want false for different-case regex query")
	}
}

func TestMatcherSupportsRegexCaseInsensitiveMatching(t *testing.T) {
	matcher, err := NewMatcherForState(State{Text: `err(or)?`, Regex: true, CaseSensitive: false})
	if err != nil {
		t.Fatalf("NewMatcherForState() error = %v, want nil", err)
	}

	if !matcher.Matches("api: ERROR opening stream") {
		t.Fatal("Matches() = false, want true for case-insensitive regex query")
	}
}

func TestMatcherReturnsValidationErrorForInvalidRegex(t *testing.T) {
	matcher, err := NewMatcherForState(State{Text: `[`, Regex: true, CaseSensitive: true})
	if err == nil {
		t.Fatal("NewMatcherForState() error = nil, want validation error")
	}
	if !IsValidationError(err) {
		t.Fatalf("IsValidationError(%T) = false, want true", err)
	}
	if matcher.Matches("api: [") {
		t.Fatal("Matches() = true for invalid matcher, want false")
	}
}

func TestLinesReturnsOnlyMatchingLinesInOrder(t *testing.T) {
	got := Lines([]string{
		"api: started",
		"worker: ready",
		"api: accepted request",
	}, "api")

	want := []string{"api: started", "api: accepted request"}
	if len(got) != len(want) {
		t.Fatalf("len(Lines()) = %d, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Lines()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestLinesWithStateReturnsValidationErrorForInvalidRegex(t *testing.T) {
	got, err := LinesWithState([]string{"api: started"}, State{Text: `[`, Regex: true, CaseSensitive: true})
	if err == nil {
		t.Fatal("LinesWithState() error = nil, want validation error")
	}
	if got != nil {
		t.Fatalf("LinesWithState() lines = %#v, want nil", got)
	}
}
