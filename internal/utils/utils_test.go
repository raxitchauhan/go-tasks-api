package utils_test

import (
	"testing"

	"go-tasks-api/internal/utils"

	"github.com/google/uuid"
)

func TestGetMockUUID(t *testing.T) {
	id := utils.GetMockUUID()
	if id == uuid.Nil {
		t.Errorf("expected non-nil UUID, got %v", id)
	}
}

func TestTrimString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"trim spaces", "  hello  ", "hello"},
		{"no spaces", "world", "world"},
		{"all spaces", "   ", ""},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.TrimString(tt.input)
			if got != tt.expected {
				t.Errorf("TrimString(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}
