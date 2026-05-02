package domain

import "testing"

func TestContainerDisplayName(t *testing.T) {
	tests := []struct {
		name      string
		container Container
		want      string
	}{
		{
			name:      "uses container name first",
			container: Container{ID: "abc123", Name: "api"},
			want:      "api",
		},
		{
			name:      "falls back to id",
			container: Container{ID: "abc123"},
			want:      "abc123",
		},
		{
			name:      "falls back to unknown",
			container: Container{},
			want:      "<unknown>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.container.DisplayName(); got != tt.want {
				t.Fatalf("DisplayName() = %q, want %q", got, tt.want)
			}
		})
	}
}
