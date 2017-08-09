package brain

import (
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// DistributionDefinitions represents more than one definition in output.Outputtable form.
type DistributionDefinitions []DistributionDefinition

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (ds DistributionDefinitions) DefaultFields(f output.Format) string {
	return (DistributionDefinition{}).DefaultFields(f)
}

// PrettyPrint writes a human-readable summary of the definitions to wr at the given detail level.
func (ds DistributionDefinitions) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) (err error) {
	for _, d := range ds {
		err = d.PrettyPrint(wr, detail)
		if err != nil {
			return
		}
	}
	return
}

// HardwareProfileDefinitions represents more than one definition in output.Outputtable form.
type HardwareProfileDefinitions []HardwareProfileDefinition

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (ds HardwareProfileDefinitions) DefaultFields(f output.Format) string {
	return (HardwareProfileDefinition{}).DefaultFields(f)
}

// PrettyPrint writes a human-readable summary of the definitions to wr at the given detail level.
func (ds HardwareProfileDefinitions) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) (err error) {
	for _, d := range ds {
		err = d.PrettyPrint(wr, detail)
		if err != nil {
			return
		}
	}
	return
}

// StorageGradeDefinitions represents more than one definition in output.Outputtable form.
type StorageGradeDefinitions []StorageGradeDefinition

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (ds StorageGradeDefinitions) DefaultFields(f output.Format) string {
	return (StorageGradeDefinition{}).DefaultFields(f)
}

// PrettyPrint writes a human-readable summary of the definitions to wr at the given detail level.
func (ds StorageGradeDefinitions) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) (err error) {
	for _, d := range ds {
		err = d.PrettyPrint(wr, detail)
		if err != nil {
			return
		}
	}
	return
}

// ZoneDefinitions represents more than one definition in output.Outputtable form.
type ZoneDefinitions []ZoneDefinition

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (ds ZoneDefinitions) DefaultFields(f output.Format) string {
	return (ZoneDefinition{}).DefaultFields(f)
}

// PrettyPrint writes a human-readable summary of the definitions to wr at the given detail level.
func (ds ZoneDefinitions) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) (err error) {
	for _, d := range ds {
		err = d.PrettyPrint(wr, detail)
		if err != nil {
			return
		}
	}
	return
}
