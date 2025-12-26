package utils

import (
	"strings"

	"github.com/google/uuid"
)

// Returns mock uuid
func GetMockUUID() uuid.UUID {
	return uuid.New()
}

// TrimString trims leading and trailing spaces from a string
func TrimString(s string) string {
	return strings.Trim(s, " ")
}
