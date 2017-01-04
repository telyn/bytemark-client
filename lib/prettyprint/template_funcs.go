package prettyprint

import (
	"fmt"
	"strings"
	"unicode"
)

var templateFuncMap = map[string]interface{}{
	"mibgib": func(size int) string {
		mg := 'M'
		if size >= 1024 {
			size /= 1024
			mg = 'G'
		}
		return fmt.Sprintf("%d%ciB", size, mg)
	},
	"gibtib": func(size int) string {
		// lt is a less than sign if < 1GiB
		lt := ""
		if size < 1024 {
			lt = "< "
		}
		size /= 1024
		gt := 'G'
		if size >= 1024 {
			size /= 1024
			gt = 'T'
		}
		return fmt.Sprintf("%s%d%ciB", lt, size, gt)
	},
	"capitalize": func(str string) string {
		if len(str) == 0 {
			return str
		}

		runes := []rune(str)
		runes[0] = unicode.ToUpper(runes[0])
		return string(runes)
	},
	"pluralize": func(single string, plural string, num int) string {
		if num == 1 {
			return fmt.Sprintf("%d %s", num, single)
		}
		return fmt.Sprintf("%d %s", num, plural)
	},
	"percentage": func(num int, denom int) string {
		return fmt.Sprintf("%d%%", int(100*float64(num)/float64(denom)))
	},
	"join": strings.Join,
}
