package slug

import (
	"regexp"
	"strings"
)

var (
	nonAlphanumericRegex = regexp.MustCompile(`[^a-z0-9]+`)
	multipleHyphensRegex = regexp.MustCompile(`-+`)
)

func New(name string) string {
	lowerCaseName := strings.ToLower(name)
	slug := nonAlphanumericRegex.ReplaceAllString(lowerCaseName, "-")
	slug = multipleHyphensRegex.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}
