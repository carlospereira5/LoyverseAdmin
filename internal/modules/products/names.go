package products

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var _titleCaser = cases.Title(language.Spanish)

// standardizeName applies Title Case and normalizes internal whitespace.
// "  yerba   mate  " → "Yerba Mate"
func standardizeName(name string) string {
	normalized := strings.Join(strings.Fields(strings.TrimSpace(name)), " ")
	return _titleCaser.String(normalized)
}
