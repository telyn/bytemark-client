package brain

import "fmt"

type AccountID int

func (id AccountID) AccountPath() (string, error) {
	return fmt.Sprintf("/accounts/%d", id), nil
}

type DiscID int

func (id DiscID) DiscPath() (string, error) {
	return fmt.Sprintf("/discs/%d", id), nil
}

type GroupID int

func (id GroupID) GroupPath() (string, error) {
	return fmt.Sprintf("/groups/%d", id), nil
}

type VirtualMachineID int

func (id VirtualMachineID) VirtualMachinePath() (string, error) {
	return fmt.Sprintf("/virtual_machines/%d", id), nil
}
