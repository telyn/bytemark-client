package main

import (
	"encoding/json"
	"fmt"
	"github.com/BytemarkHosting/row"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	"reflect"
	"strings"
)

// OutputJSON is an OutputFn which outputs a nicely-indented JSON object that represents obj
func (c *Context) OutputJSON(obj interface{}) error {
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

func (c *Context) determineTableFields(obj interface{}) []string {
	chosenFields := strings.Split(c.String("table-fields"), ",")
	trimAllSpace(chosenFields)

	if len(chosenFields) > 0 && chosenFields[0] != "" {
		return chosenFields
	}

	fieldsList := row.FieldsFrom(obj)
	return fieldsList
}

// OutputTable is an OutputFn which outputs the object in table form, using github.com/BytemarkHosting/row and github.com/olekukonko/tablewriter
func (c *Context) OutputTable(obj interface{}) error {
	if c.String("table-fields") == "help" {
		fieldsList := row.FieldsFrom(obj)
		fmt.Fprintf(global.App.Writer, "Table fields available for this command: \r\n  %s\r\n\r\n", strings.Join(fieldsList, "\r\n  "))
		return nil
	}
	fields := c.determineTableFields(obj)
	return RenderTable(obj, fields)
}

// RenderTable creates a table for the given object. This makes
// most sense when it's an array, but a regular struct-y object works fine too.
func RenderTable(obj interface{}, fields []string) error {
	table := tablewriter.NewWriter(global.App.Writer)
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

const (
	// DefaultAccountTableFields is the default for --table-fields for lib.Account
	DefaultAccountTableFields = "BillingID, Name, Suspended, Groups"
	// DefaultBackupTableFields is the default for --table-fields for brain.Backup
	DefaultBackupTableFields = "ID, Manual, Label, StorageGrade, Size, BackupCount, BackupSchedules"
	// DefaultBackupScheduleTableFields is the default for --table-fields for brain.BackupSchedule
	DefaultBackupScheduleTableFields = "ID, StartDate, Interval"
	// DefaultDiscTableFields is the default for --table-fields for brain.Disc
	DefaultDiscTableFields = "ID, Label, StorageGrade, Size, BackupCount, BackupSchedules"
	// DefaultGroupTableFields is the default for --table-fields for brain.Group
	DefaultGroupTableFields = "ID, Name, VirtualMachines"
	// DefaultPrivilegeTableFields is the default for --table-fields for brain.Privilege
	DefaultPrivilegeTableFields = "ID, Username, Level, Target, YubikeyRequired"
	// DefaultServerTableFields is the default for --table-fields for brain.VirtualMachine
	DefaultServerTableFields = "ID, Hostname, ManagementAddress, Memory, Cores, Discs, CdromURL, Autoreboot, PowerOn, Deleted"

	// DefaultHeadTableFields is the default for --table-fields for brain.Head
	DefaultHeadTableFields = "ID, Label, IsOnline, UsageStrategy, UUID, CCAddress, VirtualMachineCount, MemoryFree, UsedCores, Memory, Note, Architecture, Models, ZoneName"
	// DefaultTailTableFields is the default for --table-fields for brain.Tail
	DefaultTailTableFields = "ID, Label, IsOnline, UUID, CCAddress, StoragePools, ZoneName"
	// DefaultStoragePoolTableFields is the default for --table-fields for brain.StoragePool
	DefaultStoragePoolTableFields = "Label, Discs, Name, Size, FreeSpace, StorageGrade, UsageStrategy, OvercommitRatio, Note, Zone"
	// DefaultIPRangeTableFields is the default for --table-fields for brain.IPRange
	DefaultIPRangeTableFields = "ID, Spec, VLANNum, Available, Zones"
	// DefaultVLANTableFields is the default for --table-fields for brain.VLAN
	DefaultVLANTableFields = "ID, Num, UsageType, IPRanges"

	// DefaultDefinitionTableFields is the default for --table-fields for the *Definition types, because they're all the same at the moment.
	DefaultDefinitionTableFields = "Name, Description"
)

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

// OutputFn is a function for outputting an object to the terminal in some way
// See the OutputFormatFns map to see examples
type OutputFn func(context *Context, obj interface{}) error

// OutputFormatFns is a map which contains all the supported output format functions -- except 'human' because that's implemented in the OutputInDesiredForm method, by necessity.
var OutputFormatFns = map[string]OutputFn{
	"debug": func(c *Context, obj interface{}) error {
		fmt.Fprintf(global.App.Writer, "%#v", obj)
		return nil
	},
	"json":  (*Context).OutputJSON,
	"table": (*Context).OutputTable,
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
