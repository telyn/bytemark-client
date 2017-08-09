package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
	"github.com/telyn/row"
)

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

// OutputFn is a function for outputting an object to the terminal in some way
// See the OutputFormatFns map to see examples
type OutputFn func(wr io.Writer, config Config, obj Outputtable) error

// OutputFormatFns is a map which contains all the supported output format functions -- except 'human' because that's implemented in the OutputInDesiredForm method, by necessity.
var OutputFormatFns = map[Format]OutputFn{
	Debug: func(wr io.Writer, cfg Config, obj Outputtable) error {
		fmt.Fprintf(wr, "%#v", obj)
		return nil
	},
	JSON: func(wr io.Writer, _ Config, obj Outputtable) error {
		encoder := json.NewEncoder(wr)
		encoder.SetIndent("", "    ")
		return encoder.Encode(obj)
	},
	List:  outputTable,
	Table: outputTable,
	Human: func(wr io.Writer, _ Config, obj Outputtable) error {
		return obj.PrettyPrint(wr, prettyprint.Full)
	},
}

// outputTable is an OutputFn, used by List and Table output types
func outputTable(wr io.Writer, cfg Config, obj Outputtable) error {
	if cfg.Fields[0] == "help" {
		fieldsList := row.FieldsFrom(obj)
		fmt.Fprintf(wr, "Fields available for this command: \r\n  %s\r\n\r\n", strings.Join(fieldsList, "\r\n  "))
		return nil
	}
	return RenderTable(wr, cfg, obj)
}

// FormatByName returns the Format for the given format name. If the name is not valid, returns Human
func FormatByName(name string) Format {
	name = strings.ToLower(name)
	for f := range OutputFormatFns {
		if string(f) == name {
			return f
		}
	}
	return Human
}

// SupportedOutputFormats returns a list of all suppported output forms, including 'human'
func SupportedOutputFormats() (outputFormats []string) {
	outputFormats = make([]string, 0, len(OutputFormatFns)+1)
	for k := range OutputFormatFns {
		outputFormats = append(outputFormats, string(k))
	}
	outputFormats = append(outputFormats, "human")
	return
}
