package app

import (
	"bytes"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/config"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/urfave/cli"
)

func TestOutput(t *testing.T) {

	tests := []struct {
		ShouldErr     bool
		DefaultFormat []output.Format
		ConfigFormat  config.Var
		JSONFlag      bool
		TableFlag     bool
		Object        output.Outputtable
		Expected      string
		TableFields   string
	}{
		{ // 0
			// default to human output
			ConfigFormat: config.Var{Name: "output-format", Value: "human", Source: "CODE"},
			Object: brain.Disc{
				Label:        "disk-1",
				StorageGrade: "sata",
				Size:         25660,
				ID:           123,
			},
			Expected: "disk-1 - 25GiB, sata grade\n",
		}, { // 1
			// default to human output with an error
			ConfigFormat: config.Var{Name: "output-format", Value: "human", Source: "CODE"},
			Object:       nil,
			Expected:     "",
			ShouldErr:    true,
		}, { // 2
			// when there's a default format specific to the command, use that instead of the uber-default
			ConfigFormat:  config.Var{Name: "output-format", Value: "human", Source: "CODE"},
			DefaultFormat: []output.Format{output.Table},
			Object: brain.Disc{
				StorageGrade: "sata",
				Size:         25660,
				ID:           123,
			},
			TableFields: "ID",
			Expected:    "+-----+\n| ID  |\n+-----+\n| 123 |\n+-----+\n",
		}, { // 3
			// except when the JSON flag is set, then output JSON
			ConfigFormat:  config.Var{Name: "output-format", Value: "human", Source: "CODE"},
			DefaultFormat: []output.Format{output.Table},
			JSONFlag:      true,
			Object: brain.Group{
				Name: "my-cool-group",
				ID:   11323,
			},
			Expected: "{\n    \"name\": \"my-cool-group\",\n    \"account_id\": 0,\n    \"id\": 11323,\n    \"virtual_machines\": null\n}\n",
		}, { // 4
			// or if output-format is set by a FILE
			ConfigFormat:  config.Var{Name: "output-format", Value: "json", Source: "FILE"},
			DefaultFormat: []output.Format{output.Table},
			Object: brain.Group{
				Name: "my-cool-group",
				ID:   11323,
			},
			Expected: "{\n    \"name\": \"my-cool-group\",\n    \"account_id\": 0,\n    \"id\": 11323,\n    \"virtual_machines\": null\n}\n",
			// but the table and json flags should have precedence in every situation
		}, { // 5
			ConfigFormat:  config.Var{Name: "output-format", Value: "json", Source: "FILE"},
			DefaultFormat: []output.Format{output.Human},
			TableFlag:     true,
			Object: brain.Group{
				Name: "my-cool-group",
				ID:   11323,
			},
			Expected: "+---------------+-----------------+\n|     Name      | VirtualMachines |\n+---------------+-----------------+\n| my-cool-group |                 |\n+---------------+-----------------+\n",
			// also, --table-fields being non-empty should imply --table and be case insensitive
		}, { // 6
			ConfigFormat:  config.Var{Name: "output-format", Value: "json", Source: "FILE"},
			DefaultFormat: []output.Format{output.Human},
			TableFlag:     false,
			TableFields:   "AccountID,ID,Name,VirtualMachines",
			Object: brain.Group{
				Name: "my-cool-group",
				ID:   11323,
			},
			Expected: "+-----------+-------+---------------+-----------------+\n| AccountID |  ID   |     Name      | VirtualMachines |\n+-----------+-------+---------------+-----------------+\n|         0 | 11323 | my-cool-group |                 |\n+-----------+-------+---------------+-----------------+\n",
		},
	}

	for i, test := range tests {
		t.Logf("TestOutput %d\n", i)
		config, _, app := baseTestSetup(t, true, []cli.Command{})
		config.Reset()

		cliContext := &mocks.CliContext{}
		cliContext.When("App").Return(app)
		context := Context{Context: cliContext}

		config.When("GetBool", "admin").Return(true)
		config.When("GetV", "output-format").Return(test.ConfigFormat)
		cliContext.When("Bool", "json").Return(test.JSONFlag)
		cliContext.When("Bool", "table").Return(test.TableFlag)
		cliContext.When("GlobalString", "table-fields").Return(test.TableFields)
		cliContext.When("String", "table-fields").Return(test.TableFields)
		cliContext.When("IsSet", "table-fields").Return(test.TableFields != "")

		var err error
		if test.DefaultFormat == nil {
			err = context.OutputInDesiredForm(test.Object)
		} else {
			err = context.OutputInDesiredForm(test.Object, test.DefaultFormat...)
		}
		if err != nil && !test.ShouldErr {
			t.Errorf("TestOutput %d ERR: %s", i, err)
		} else if err == nil && test.ShouldErr {
			t.Errorf("TestOutput %d Didn't error", i)
		}

		buf := app.Metadata["buf"].(*bytes.Buffer)
		output := buf.String()
		if output != test.Expected {
			t.Errorf("Output for %d didn't match expected.\r\nExpected: %q\r\nActual:   %q", i, test.Expected, output)
		}
	}

}
