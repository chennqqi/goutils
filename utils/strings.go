package utils

import (
	"strings"
)

func Index(s, substr string, start int) int {
	return start + strings.Index(s[start:], substr)
}
