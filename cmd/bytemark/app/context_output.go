package app

import (
	"fmt"
	"io"
	"strings"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/urfave/cli"
)

func trimAllSpace(strs []string) {
	for i, s := range strs {
		strs[i] = strings.TrimSpace(s)
	}
}

// Debug runs fmt.Fprintf on the args, outputting to the App's debugWriter.
// In tests, this is a TestWriter. Otherwise it's nil for now - but might be
// changed to the debug.log File in the future.
func (c *Context) Debug(format string, values ...interface{}) {
	dw, ok := c.App().Metadata["debugWriter"]
	if !ok {
		return
	}
	if wr, ok := dw.(io.Writer); ok {
		fmt.Fprintf(wr, format, values...)
	}
}

// Log runs fmt.Fprintf on the args, outputting to the App's Writer
func (c *Context) Log(format string, values ...interface{}) {
	fmt.Fprintf(c.App().Writer, format+"\n", values...)
}

// LogErr runs fmt.Fprintf on the args, outputting to the App's Writer
func (c *Context) LogErr(format string, values ...interface{}) {
	fmt.Fprintf(c.App().ErrWriter, format+"\n", values...)
}

// OutputFormat attempts to figure out the output format needed, given the contents of the output-format config var,
// the json flag, and the table and table-fields flag. If there is an error reading the config, it is returned and
// human output is assumed.
func (c *Context) OutputFormat(defaultFormat ...output.Format) (output.Format, error) {
	format, err := c.Config().GetV("output-format")
	if err != nil {
		return output.Human, err
	}
	if len(defaultFormat) > 0 && format.Source == "CODE" {
		format.Value = string(defaultFormat[0])
	}

	if c.Bool("json") {
		format.Value = "json"
	} else if c.Bool("table") || c.Context.IsSet("table-fields") {
		format.Value = "table"
	}

	return output.FormatByName(format.Value), nil

}

func (c *Context) createOutputConfig(obj output.DefaultFieldsHaver, defaultFormat ...output.Format) (cfg output.Config, err error) {
	cfg = output.Config{}
	cfg.Format, err = c.OutputFormat(defaultFormat...)

	cfg.Fields = strings.Split(c.String("table-fields"), ",")
	trimAllSpace(cfg.Fields)

	if len(cfg.Fields) > 0 && cfg.Fields[0] != "" {
		return
	}
	cfg.Fields = strings.Split(obj.DefaultFields(cfg.Format), ",")
	trimAllSpace(cfg.Fields)

	return
}

// OutputFlags creates some cli.Flags for when you wanna use OutputInDesiredForm
// thing should be like "server", "servers", "group", "groups"
// jsonType should be "array" or "object"
func OutputFlags(thing string, jsonType string) []cli.Flag {
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
		},
	}
}

// OutputInDesiredForm outputs obj as a JSON object if --json is set,
// or as a table / table row if --table is set
// otherwise calls humanOutputFn (which should output it in a very human form - PrettyPrint or such
// defaultFormat is an optional string stating what the default format should be
func (c *Context) OutputInDesiredForm(obj output.Outputtable, defaultFormat ...output.Format) error {
	if obj == nil {
		return fmt.Errorf("Object passed to OutputInDesiredForm was nil")
	}
	cfg, err := c.createOutputConfig(obj, defaultFormat...)
	if err != nil {
		return err
	}
	return output.Write(c.App().Writer, cfg, obj)
}
