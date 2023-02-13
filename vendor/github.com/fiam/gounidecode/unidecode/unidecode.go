// Unidecode implements a unicode transliterator, which
// replaces non-ASCII characters with their ASCII
// counterparts

package unidecode

import (
	"unicode"
)

// Given an unicode encoded string, returns
// another string with non-ASCII characters replaced
// with their closest ASCII counterparts.
// e.g. Unicode("áéíóú") => "aeiou" 
func Unidecode(s string) string {
	str := ""
	for _, c := range s {
		if c <= unicode.MaxASCII {
			str += string(c)
			continue
		}
		if c > unicode.MaxRune {
			/* Ignore reserved chars */
			continue
		}
		if d, ok := transliterations[c]; ok {
			str += d
		}
	}
	return str
}
