package brain

import (
	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
	"io"
)

// Disc is a representation of a VM's disc.
type Disc struct {
	Label        string `json:"label"`
	StorageGrade string `json:"storage_grade"`
	Size         int    `json:"size"`

	ID               int    `json:"id,omitempty"`
	VirtualMachineID int    `json:"virtual_machine_id,omitempty"`
	StoragePool      string `json:"storage_pool,omitempty"`
}

func (d Disc) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	tmpl := `{{ define "disc_sgl" }}{{ .Label }} - {{ gibtib .Size }}, {{ .StorageGrade }} grade{{ end }}
{{ define "disc_medium" }}{{ template "_sgl" . }}{{ end }}
{{ define "disc_full" }}{{ template "_medium" . }}{{ end }}`
	return prettyprint.Run(wr, tmpl, "disc"+string(detail), d)
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
