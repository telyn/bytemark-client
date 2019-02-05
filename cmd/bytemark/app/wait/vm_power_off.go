package wait

import (
	"fmt"
	"time"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
)

// VMPowerOff waits for the named virtual machine to power off before returning
// a nil error. This is done by frequently polling the brain for info about the
// VM. If any calls fail, the error is returned.
func VMPowerOff(c *app.Context, name pathers.VirtualMachineName) (err error) {
	vm := brain.VirtualMachine{PowerOn: true}

	for vm.PowerOn {
		if !c.IsTest() {
			time.Sleep(5 * time.Second)
		}
		// this Fprint is actually needed as there's no way to get ctx.Log to not output a newline after
		_, _ = fmt.Fprint(c.App().Writer, ".")

		vm, err = c.Client().GetVirtualMachine(name)
		if err != nil {
			return
		}
	}
	return
}
