package main

import (
	"fmt"
	"strings"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/urfave/cli"
)

func trimAllSpace(strs []string) {
	for i, s := range strs {
		strs[i] = strings.TrimSpace(s)
	}
}

func (c *Context) determineOutputFormat(defaultFormat ...string) (output.Format, error) {
	format, err := global.Config.GetV("output-format")
	if err != nil {
		return output.Human, err
	}
	if len(defaultFormat) > 0 && format.Source == "CODE" {
		format.Value = defaultFormat[0]
	}

	if c.Bool("json") {
		format.Value = "json"
	} else if c.Bool("table") || c.Context.IsSet("table-fields") {
		format.Value = "table"
	}

	return output.FormatByName(format.Value), nil

}

func (c *Context) CreateOutputConfig(obj output.DefaultFieldsHaver, defaultFormat ...string) (cfg output.Config, err error) {
	cfg = output.Config{}
	cfg.Format, err = c.determineOutputFormat(defaultFormat...)

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

// OutputInDesiredForm outputs obj as a JSON object if --json is set,
// or as a table / table row if --table is set
// otherwise calls humanOutputFn (which should output it in a very human form - PrettyPrint or such
// defaultFormat is an optional string stating what the default format should be
func (c *Context) OutputInDesiredForm(obj output.Outputtable, defaultFormat ...string) error {
	cfg, err := c.CreateOutputConfig(obj, defaultFormat...)
	if err != nil {
		return err
	}
	return output.Write(global.App.Writer, cfg, obj)
}
