package output

import (
	"fmt"
	"io"
	"reflect"

	"github.com/olekukonko/tablewriter"
	"github.com/telyn/row"
)

// SetupTable creates a tablewriter.Table for the given writer and output config.
func SetupTable(wr io.Writer, cfg Config) (table *tablewriter.Table) {
	table = tablewriter.NewWriter(wr)
	// not sure SetHeader will ALWAYS be useful for every kind of table we wanna write, but it's useful for table and list, so use it here.
	table.SetHeader(cfg.Fields)
	switch cfg.Format {
	case List:
		// autowrap - lists want to be contained on a single line
		table.SetAutoWrapText(true)
		// no row line - we wanna be able to grep stuff
		table.SetRowLine(false)
		// CAPITAL headers because... well `docker ps` has them
		table.SetAutoFormatHeaders(true)
	default:
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
	}

	return
}

// RenderTable creates a table for the given object. This makes
// most sense when it's an array, but a regular struct-y object works fine too.
func RenderTable(wr io.Writer, cfg Config, obj interface{}) error {
	table := SetupTable(wr, cfg)

	v := reflect.ValueOf(obj)

	// indirect pointers so we can switch on Kind()
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// output a single table row for a struct, or several for a slice / array
	switch v.Kind() {
	case reflect.Struct:
		r, err := row.From(obj, cfg.Fields)
		if err != nil {
			return err
		}
		table.Append(r)
	case reflect.Slice, reflect.Array:
		length := v.Len()
		for i := 0; i < length; i++ {
			el := v.Index(i)
			r, err := row.From(el.Interface(), cfg.Fields)
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
