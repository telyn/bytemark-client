package main

import (
	"bytes"
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"testing"
)

func TestOutput(t *testing.T) {
	oldWriter := global.App.Writer

	humanFnOK := func() error {
		fmt.Fprint(global.App.Writer, "OK")
		return nil
	}

	humanFnErr := func() error {
		fmt.Fprint(global.App.Writer, "NOT OK")
		return fmt.Errorf("humanFnErr called")
	}

	tests := []struct {
		ShouldErr     bool
		DefaultFormat []string
		ConfigFormat  util.ConfigVar
		JSONFlag      bool
		TableFlag     bool
		HumanFn       func() error
		Object        interface{}
		Expected      string
		TableFields   string
	}{
		{ // 0
			// default to human output
			ConfigFormat: util.ConfigVar{"output-format", "human", "CODE"},
			HumanFn:      humanFnOK,
			Object: brain.Disc{
				StorageGrade: "sata",
				Size:         25660,
				ID:           123,
			},
			Expected: "OK",
		}, { // 1
			// default to human output with an error
			ConfigFormat: util.ConfigVar{"output-format", "human", "CODE"},
			HumanFn:      humanFnErr,
			Object:       nil,
			Expected:     "NOT OK",
			ShouldErr:    true,
		}, { // 2
			// when there's a default format specific to the command, use that instead of the uber-default
			ConfigFormat:  util.ConfigVar{"output-format", "human", "CODE"},
			DefaultFormat: []string{"table"},
			HumanFn:       humanFnErr,
			Object: brain.Disc{
				StorageGrade: "sata",
				Size:         25660,
				ID:           123,
			},
			TableFields: "ID",
			Expected:    "+-----+\n| ID  |\n+-----+\n| 123 |\n+-----+\n",
		}, { // 3
			// except when the JSON flag is set, then output JSON
			ConfigFormat:  util.ConfigVar{"output-format", "human", "CODE"},
			DefaultFormat: []string{"table"},
			JSONFlag:      true,
			HumanFn:       humanFnErr,
			Object: brain.Group{
				Name: "my-cool-group",
				ID:   11323,
			},
			Expected: "{\n    \"name\": \"my-cool-group\",\n    \"account_id\": 0,\n    \"id\": 11323,\n    \"virtual_machines\": null\n}",
		}, { // 4
			// or if output-format is set by a FILE
			ConfigFormat:  util.ConfigVar{"output-format", "json", "FILE"},
			DefaultFormat: []string{"table"},
			HumanFn:       humanFnErr,
			Object: brain.Group{
				Name: "my-cool-group",
				ID:   11323,
			},
			Expected: "{\n    \"name\": \"my-cool-group\",\n    \"account_id\": 0,\n    \"id\": 11323,\n    \"virtual_machines\": null\n}",
			// but the table and json flags should have precedence in every situation
		}, { // 5
			ConfigFormat:  util.ConfigVar{"output-format", "json", "FILE"},
			DefaultFormat: []string{"human"},
			HumanFn:       humanFnErr,
			TableFlag:     true,
			Object: brain.Group{
				Name: "my-cool-group",
				ID:   11323,
			},
			Expected: "+----------------------+---------------------------------------------+---------------+-----------+-------+-----------------+\n| CountVirtualMachines |                   String                    |     Name      | AccountID |  ID   | VirtualMachines |\n+----------------------+---------------------------------------------+---------------+-----------+-------+-----------------+\n|                    0 | group 11323 \"my-cool-group\" - has 0 servers | my-cool-group |         0 | 11323 |                 |\n+----------------------+---------------------------------------------+---------------+-----------+-------+-----------------+\n",
			// also, --table-fields being non-empty should imply --table and be case insensitive
		}, { // 6
			ConfigFormat:  util.ConfigVar{"output-format", "json", "FILE"},
			DefaultFormat: []string{"human"},
			HumanFn:       humanFnErr,
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
		config, _ := baseTestSetup(t, true)
		config.Reset()

		cliContext := &mocks.CliContext{}
		cliContext.When("App").Return(global.App)
		context := Context{Context: cliContext}

		config.When("GetBool", "admin").Return(true)
		config.When("GetV", "output-format").Return(test.ConfigFormat)
		cliContext.When("Bool", "json").Return(test.JSONFlag)
		cliContext.When("Bool", "table").Return(test.TableFlag)
		cliContext.When("GlobalString", "table-fields").Return(test.TableFields)
		cliContext.When("String", "table-fields").Return(test.TableFields)
		cliContext.When("IsSet", "table-fields").Return(test.TableFields != "")
		global.Config = config

		buf := bytes.Buffer{}
		global.App.Writer = &buf

		var err error
		if test.DefaultFormat == nil {
			err = context.OutputInDesiredForm(test.Object, test.HumanFn)
		} else {
			err = context.OutputInDesiredForm(test.Object, test.HumanFn, test.DefaultFormat...)
		}
		if err != nil && !test.ShouldErr {
			t.Errorf("TestOutput %d ERR: %s", i, err)
		} else if err == nil && test.ShouldErr {
			t.Errorf("TestOutput %d Didn't error", i)
		}

		output := buf.String()
		if output != test.Expected {
			t.Errorf("Output for %d didn't match expected.\r\nExpected: %q\r\nActual:   %q", i, test.Expected, output)
		}
		global.App.Writer = oldWriter
	}

}
