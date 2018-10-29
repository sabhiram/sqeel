package sqeel

import (
	"strings"
	"unicode"
)

func ToSnakeCase(s string) string {
	// If current letter is uppercase - grab letters until the next uppercase
	// letter that precedes a lowercase letter.
	runes := []rune(s)
	rlen := len(runes)
	words := []string{}
	lasti := 0
	for i := 0; i < rlen; i++ {
		if i > 1 && unicode.IsUpper(runes[i]) &&
			(unicode.IsLower(runes[i-1]) || (rlen > i+1 && unicode.IsLower(runes[i+1]))) {
			words = append(words, string(runes[lasti:i]))
			lasti = i
		}
	}
	if lasti < rlen {
		words = append(words, string(runes[lasti:rlen]))
	}
	return strings.ToLower(strings.Join(words, "_"))
}
