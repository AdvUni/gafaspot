package util

import (
	"regexp"
	"strings"
)

// CreatePlainIdentifier replaces all characters which are not ascii letters oder numbers through an underscore
func CreatePlainIdentifier(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	return strings.ToLower(re.ReplaceAllString(name, "_"))
}
