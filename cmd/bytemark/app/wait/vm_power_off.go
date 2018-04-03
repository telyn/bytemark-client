package wait

import (
	"fmt"
	"time"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

func VMPowerOff(c *app.Context, name lib.VirtualMachineName) (err error) {
	vm := brain.VirtualMachine{PowerOn: true}

	for vm.PowerOn {
		if !c.IsTest() {
			time.Sleep(5 * time.Second)
		}
		fmt.Fprint(c.App().Writer, ".")

		vm, err = c.Client().GetVirtualMachine(name)
		if err != nil {
			return
		}
	}
	return
}
