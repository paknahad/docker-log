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
