package output

import "strings"

// Format is a canonical name of output formats
type Format string

const (
	// List is the canonical name of the List output format
	List Format = "list"
	// Table is the canonical name of the Table output format
	Table = "table"
	// JSON is the canonical name of the JSON output format
	JSON = "json"
	// Human is the canonical name of the Human output format
	Human = "human"
	// Debug is the canonical name of the Debug output format
	Debug = "debug"
)

// FormatByName returns the Format for the given format name. If the name is not valid, returns Human
func FormatByName(name string) Format {
	formats := map[string]Format{
		"debug": Debug,
		"list":  List,
		"json":  JSON,
		"table": Table,
		"human": Human,
	}
	name = strings.ToLower(name)
	if f, ok := formats[name]; ok {
		return f
	} else {
		return Human
	}
}
