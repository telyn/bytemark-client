package prettyprint

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

// TemplateFragmentMapper is an interface that requires MapTemplateFragment to exist. The implementation of MapTemplateFragment should add "{{" and "}}" to the beginning and end of template fragment, and run it for each element in the receiver, collecting the output up into an array of strings.
type TemplateFragmentMapper interface {
	MapTemplateFragment(templateFrag string) (strs []string, err error)
}

func prefixEachLine(prefix string, input string) string {
	lines := strings.Split(input, "\n")
	for i := range lines {
		if lines[i] != "" {
			lines[i] = string(prefix) + lines[i]
		}
	}
	return strings.Join(lines, "\n")
}

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
	// indent indents every line of the input string by n spaces.
	// for example indent 1 " hi\nhello" would produce "  hi\n hello"
	"indent": func(amount int, input string) string {
		spaces := make([]rune, amount)
		for i := range spaces {
			spaces[i] = ' '
		}
		return prefixEachLine(string(spaces), input)

	},
	"prefixEachLine": prefixEachLine,
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
	// pluralize returns single if num == 1, plural if num == 2. For instance, {{ pluralize "horse" "horses" 2 }} will return "2 horses", while {{ pluralize "bacterium" "bacteria" 1 }} will return "1 bacterium".
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
	// map ... see TemplateFragmentMapper and brain.BackupSchedules
	"map": func(mapper TemplateFragmentMapper, fragment string) ([]string, error) {
		return mapper.MapTemplateFragment(fragment)
	},
	// join joins multiple strings together with a separator
	"join": strings.Join,
	// joinWithSpecialLast joins multiple strings together with a separator, except the last two, which are seperated by a different seperator. e.g. joinWithSpecialLast ", " " and " []string{"hi","hello","welcome","good evening"} would produce "hi, hello, welcome and good evening"
	"joinWithSpecialLast": func(sep string, fin string, strs []string) string {
		// special cases for when there are 0, 1 or 2 strings
		switch len(strs) {
		case 0:
			return ""
		case 1:
			return strs[0]
		case 2:
			return strs[0] + fin + strs[1]
		}
		// join with one seperator
		most := strings.Join(strs[0:len(strs)-1], sep)
		// and add the last with the 'fin' seperator.
		return most + fin + strs[len(strs)-1]
	},
}

// Funcs is not for general usage, only while in transition to PrettyPrint.
var Funcs = templateFuncMap
