package zmx

import (
	"strings"
	"testing"
)

func TestStripANSI(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"\x1b[31mRed\x1b[0m", "\x1b[31mRed\x1b[0m"},
		{"\x1b[1;32mBold Green\x1b[0m", "\x1b[1;32mBold Green\x1b[0m"},
		{"Plain text", "Plain text"},
		{"\x1b[H\x1b[2JClear screen", "Clear screen"}, // CSI H and J should be stripped
	}

	for _, tt := range tests {
		got := stripANSI(tt.input)
		if got != tt.expected {
			t.Errorf("stripANSI(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestScrollPreviewANSI(t *testing.T) {
	input := "\x1b[31mRed\x1b[0m and \x1b[32mGreen\x1b[0m"

	// Test basic display (no scroll)
	got := ScrollPreview(input, 0, 10)
	// Expect colors to be preserved
	if !strings.Contains(got, "\x1b[31m") || !strings.Contains(got, "\x1b[32m") {
		t.Errorf("ScrollPreview(input, 0, 10) lost colors: %q", got)
	}

	// Test scrolling (skip "Red and ")
	// "Red and " is 8 chars
	got = ScrollPreview(input, 8, 5)
	// Should contain Green color and "Green"
	if !strings.Contains(got, "\x1b[32m") || !strings.Contains(got, "Green") {
		t.Errorf("ScrollPreview(input, 8, 5) lost Green color or text: %q", got)
	}
	// Should STILL contain Red color start because we preserve all SGR
	if !strings.Contains(got, "\x1b[31m") {
		t.Errorf("ScrollPreview(input, 8, 5) lost Red color start: %q", got)
	}
}
