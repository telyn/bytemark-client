package prettyprint

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

var templateFuncMap = map[string]interface{}{
	"capitalize": func(str string) string {
		if len(str) == 0 {
			return str
		}

		runes := []rune(str)
		runes[0] = unicode.ToUpper(runes[0])
		return string(runes)
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
	"mibgib": func(size int) string {
		mg := 'M'
		if size >= 1024 {
			size /= 1024
			mg = 'G'
		}
		return fmt.Sprintf("%d%ciB", size, mg)
	},
	"percentage": func(num int, denom int) string {
		return fmt.Sprintf("%d%%", int(100*float64(num)/float64(denom)))
	},
	"pluralize": func(single string, plural string, num int) string {
		if num == 1 {
			return fmt.Sprintf("%d %s", num, single)
		}
		return fmt.Sprintf("%d %s", num, plural)
	},
	"prettysprint": func(pp PrettyPrinter, detail DetailLevel) (string, error) {
		b := new(bytes.Buffer)
		err := pp.PrettyPrint(b, detail)
		if err != nil {
			return "", err
		}
		return b.String(), nil
	},
	"join": strings.Join,
}

// Funcs is not for general usage, only while in transition to PrettyPrint.
var Funcs = templateFuncMap
