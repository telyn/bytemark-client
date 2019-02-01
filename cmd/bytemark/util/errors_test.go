package util

import (
	"fmt"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
)

func TestRecursiveDeleteGroupError(t *testing.T) {
	tests := []struct {
		//name of the test
		name string

		group  pathers.GroupName
		errors map[string]error
		output string
	}{{
		name:  "one error",
		group: pathers.GroupName{Group: "test", Account: "account"},
		errors: map[string]error{
			"vm1": fmt.Errorf("Deleting the vm totes failed, my dude"),
		},
		output: "Errors occurred while deleting VMs in group test.account: \n\tvm1: Deleting the vm totes failed, my dude",
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := RecursiveDeleteGroupError{
				Group:  test.group,
				Errors: test.errors,
			}
			if err.Error() != test.output {
				t.Errorf("Wanted %#v\nGot    %#v", test.output, err.Error())
			}
		})
	}
}
