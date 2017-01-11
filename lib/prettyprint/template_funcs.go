package prettyprint

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

var templateFuncMap = map[string]interface{}{
	// capitalize the first letter of str
	"capitalize": func(str string) string {
		if len(str) == 0 {
			return str
		}

		runes := []rune(str)
		runes[0] = unicode.ToUpper(runes[0])
		return string(runes)
	},
	// gibtib takes a size in megabytes and formats it to be in GiB (with the unit). If size is less than 1024, outputs "< 1GiB".
	// If 1 TiB or more, outputs the size in TiB.
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
	// mibgib takes a size in megabytes and formats it with a unit in either MiB or GiB, if size >= 1024.
	"mibgib": func(size int) string {
		mg := 'M'
		if size >= 1024 {
			size /= 1024
			mg = 'G'
		}
		return fmt.Sprintf("%d%ciB", size, mg)
	},
	// percentage takes a fractions represented as two ints and returns a string showing the percentage. For instance, {{ percentage 1 2 }} will return "50%"
	"percentage": func(num int, denom int) string {
		return fmt.Sprintf("%d%%", int(100*float64(num)/float64(denom)))
	},
	// pluralize returns single if num == 1, plural if num == 2. For instance, {{ pluralize "horse" "horses" 2 }} will return "horses", while {{ pluralize "bacterium" "bacteria" 1 }} will return "bacterium".
	"pluralize": func(single string, plural string, num int) string {
		if num == 1 {
			return fmt.Sprintf("%d %s", num, single)
		}
		return fmt.Sprintf("%d %s", num, plural)
	},
	// prettysprint calls PrettyPrint on the prettyprinter, and returns its output.
	"prettysprint": func(pp PrettyPrinter, detail DetailLevel) (string, error) {
		b := new(bytes.Buffer)
		err := pp.PrettyPrint(b, detail)
		if err != nil {
			return "", err
		}
		return b.String(), nil
	},
	// join joins multiple strings together with a seperator
	"join": strings.Join,
}

// Funcs is not for general usage, only while in transition to PrettyPrint.
var Funcs = templateFuncMap
