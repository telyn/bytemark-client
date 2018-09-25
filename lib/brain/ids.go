package brain

import "fmt"

type AccountID int

func (id AccountID) Path() (string, error) {
	return fmt.Sprintf("/accounts/%d", id), nil
}

type DiscID int

func (id DiscID) Path() (string, error) {
	return fmt.Sprintf("/discs/%d", id), nil
}

type GroupID int

func (id GroupID) Path() (string, error) {
	return fmt.Sprintf("/groups/%d", id), nil
}

type VirtualMachineID int

func (id VirtualMachineID) Path() (string, error) {
	return fmt.Sprintf("/virtual_machines/%d", id), nil
}
