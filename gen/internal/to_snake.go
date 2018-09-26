package internal

import (
	"fmt"
	"unicode"
)

// converts AStringLikeThis into a_string_like_this
func ToSnake(camel string) string {
	snake := make([]rune, 0, len(camel))

	prevprev := ' '
	prev := ' '
	for _, c := range camel {
		fmt.Printf("prev: '%c' c: '%c'\n", prev, c)
		if unicode.IsLower(prev) && unicode.IsUpper(c) {
			snake = append(snake, '_')
		}
		if unicode.IsUpper(prevprev) && unicode.IsUpper(prev) && unicode.IsLower(c) {
			snake[len(snake)-1] = '_'
			snake = append(snake, unicode.ToLower(prev))
		}
		snake = append(snake, unicode.ToLower(c))
		prevprev = prev
		prev = c
	}
	return string(snake)
}
