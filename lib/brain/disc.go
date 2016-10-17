package brain

// Disc is a representation of a VM's disc.
type Disc struct {
	Label        string `json:"label"`
	StorageGrade string `json:"storage_grade"`
	Size         int    `json:"size"`

	ID               int    `json:"id,omitempty"`
	VirtualMachineID int    `json:"virtual_machine_id,omitempty"`
	StoragePool      string `json:"storage_pool,omitempty"`
}

// Validate makes sure the disc has a storage grade. Doesn't modify the origin disc.
func (disc Disc) Validate() (*Disc, error) {
	if disc.StorageGrade == "" {
		newDisc := disc
		newDisc.StorageGrade = "sata"
		return &newDisc, nil
	}
	return &disc, nil
}
