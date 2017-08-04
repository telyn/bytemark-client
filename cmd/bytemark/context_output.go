package main

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
	"github.com/BytemarkHosting/row"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

type outputtable interface {
	prettyprint.PrettyPrinter
	output.DefaultFieldsHaver
}

// OutputJSON is an OutputFn which outputs a nicely-indented JSON object that represents obj
func (c *Context) OutputJSON(obj outputtable) error {
	js, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return err
	}
	fmt.Fprintf(global.App.Writer, string(js))
	return nil
}

func trimAllSpace(strs []string) {
	for i, s := range strs {
		strs[i] = strings.TrimSpace(s)
	}
}

func (c *Context) determineTableFields(obj output.DefaultFieldsHaver) []string {
	fields := strings.Split(c.String("table-fields"), ",")
	trimAllSpace(fields)

	if len(fields) > 0 && fields[0] != "" {
		return fields
	}
	fields = strings.Split(obj.DefaultFields(output.Table), ",")
	trimAllSpace(fields)
	return fields
}

// OutputTable is an OutputFn which outputs the object in table form, using github.com/BytemarkHosting/row and github.com/olekukonko/tablewriter
func OutputTable(wr io.Writer, cfg OutputConfig, obj output.DefaultFieldsHaver) error {
	if cfg.Fields == "help" {
		fieldsList := row.FieldsFrom(obj)
		fmt.Fprintf(global.App.Writer, "Fields available for this command: \r\n  %s\r\n\r\n", strings.Join(fieldsList, "\r\n  "))
		return nil
	}
	fields := determineFields(cfg, obj)
	return RenderTable(wr, fields, obj)
}

func SetupTable(wr io.Writer, cfg OutputConfig) *table.Table {

}

// RenderTable creates a table for the given object. This makes
// most sense when it's an array, but a regular struct-y object works fine too.
func RenderTable(wr io.Writer, fields []string, obj interface{}) error {
	table := tablewriter.NewWriter(wr)
	// don't autowrap because fields that are slices output one element per line
	// and autowrap
	table.SetAutoWrapText(false)
	// lines between rows!
	table.SetRowLine(true)
	// don't autoformat the headers - autoformat makes them ALLCAPS which makes
	// it hard to figure out what to set --table-fields to.
	// with autoformat off, --table-fields can be set by copying and pasting
	// from the table header.
	table.SetAutoFormatHeaders(false)

	table.SetHeader(fields)
	v := reflect.ValueOf(obj)

	// indirect pointers so we can switch on Kind()
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// output a single table row for a struct, or several for a slice / array
	switch v.Kind() {
	case reflect.Struct:
		r, err := row.From(obj, fields)
		if err != nil {
			return err
		}
		table.Append(r)
	case reflect.Slice, reflect.Array:
		length := v.Len()
		for i := 0; i < length; i++ {
			el := v.Index(i)
			r, err := row.From(el.Interface(), fields)
			if err != nil {
				return err
			}
			table.Append(r)
		}
	default:
		return fmt.Errorf("%T is not a struct or slice type - please file a bug report", obj)
	}

	table.Render()
	return nil
}

// OutputFlags creates some cli.Flags for when you wanna use OutputInDesiredForm
// thing should be like "server", "servers", "group", "groups"
// jsonType should be "array" or "object"
func OutputFlags(thing string, jsonType string, defaultTableFields string) []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "json",
			Usage: fmt.Sprintf("Output the %s as a JSON %s", thing, jsonType),
		},
		cli.BoolFlag{
			Name:  "table",
			Usage: fmt.Sprintf("Output the %s as a table", thing),
		},
		cli.StringFlag{
			Name:  "table-fields",
			Usage: fmt.Sprintf("The fields of the %s to output in the table, comma separated. set to 'help' for a list of fields for this command", thing),
			Value: defaultTableFields,
		},
	}
}

type OutputConfig struct {
	Fields string
	Format output.Format
}

// OutputFn is a function for outputting an object to the terminal in some way
// See the OutputFormatFns map to see examples
type OutputFn func(wr io.Writer, fields string, obj outputtable) error

// OutputFormatFns is a map which contains all the supported output format functions -- except 'human' because that's implemented in the OutputInDesiredForm method, by necessity.
var OutputFormatFns = map[string]OutputFn{
	"debug": func(wr io.Writer, cfg OutputConfig, obj outputtable) error {
		fmt.Fprintf(wr, "%#v", obj)
		return nil
	},
	"json":  OutputJSON,
	"table": OutputTable,
	"human": func(wr io.Writer, cfg OutputConfig, obj outputtable) error {
		return obj.PrettyPrint(wr, prettyprint.Full)
	},
}

// SupportedOutputTypes returns a list of all suppported output forms, including 'human'
func SupportedOutputTypes() (outputTypes []string) {
	outputTypes = make([]string, 0, len(OutputFormatFns)+1)
	for k := range OutputFormatFns {
		outputTypes = append(outputTypes, k)
	}
	outputTypes = append(outputTypes, "human")
	return
}

// OutputInDesiredForm outputs obj as a JSON object if --json is set,
// or as a table / table row if --table is set
// otherwise calls humanOutputFn (which should output it in a very human form - PrettyPrint or such
// defaultFormat is an optional string stating that the default format should be
func (c *Context) OutputInDesiredForm(obj interface{}, humanOutputFn func() error, defaultFormat ...string) error {
	format, err := global.Config.GetV("output-format")
	if err != nil {
		return err
	}
	if len(defaultFormat) > 0 && format.Source == "CODE" {
		format.Value = defaultFormat[0]
	}

	if c.Bool("json") {
		format.Value = "json"
	} else if c.Bool("table") || c.Context.IsSet("table-fields") {
		format.Value = "table"
	}

	if format.Value == "" || format.Value == "human" {
		return humanOutputFn()
	}

	if fn, ok := OutputFormatFns[format.Value]; ok {
		return fn(c, obj)
	}

	return fmt.Errorf("%s isn't a supported output type. Use one of the following instead:\r\n%s", format.Value, strings.Join(SupportedOutputTypes(), "\r\n"))
}
