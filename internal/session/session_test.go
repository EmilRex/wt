package session

import (
	"regexp"
	"testing"
)

func TestGenerateSessionName(t *testing.T) {
	name := GenerateSessionName()

	// Should match YYYYMMDD-HHMMSS format
	pattern := `^\d{8}-\d{6}$`
	matched, err := regexp.MatchString(pattern, name)
	if err != nil {
		t.Fatalf("regex error: %v", err)
	}
	if !matched {
		t.Errorf("GenerateSessionName() = %q, want format YYYYMMDD-HHMMSS", name)
	}
}

func TestGetBranchName(t *testing.T) {
	tests := []struct {
		session string
		want    string
	}{
		{"foo", "wt-foo"},
		{"my-feature", "wt-my-feature"},
		{"20241215-143022", "wt-20241215-143022"},
	}

	for _, tt := range tests {
		got := GetBranchName(tt.session)
		if got != tt.want {
			t.Errorf("GetBranchName(%q) = %q, want %q", tt.session, got, tt.want)
		}
	}
}

func TestGetSessionFromBranch(t *testing.T) {
	tests := []struct {
		branch string
		want   string
	}{
		{"wt-foo", "foo"},
		{"wt-my-feature", "my-feature"},
		{"wt-20241215-143022", "20241215-143022"},
		{"other-branch", "other-branch"}, // no prefix, returns as-is
	}

	for _, tt := range tests {
		got := GetSessionFromBranch(tt.branch)
		if got != tt.want {
			t.Errorf("GetSessionFromBranch(%q) = %q, want %q", tt.branch, got, tt.want)
		}
	}
}
