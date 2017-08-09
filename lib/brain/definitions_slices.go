package brain

import (
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

type DistributionDefinitions []DistributionDefinition

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (ds DistributionDefinitions) DefaultFields(f output.Format) string {
	return (DistributionDefinition{}).DefaultFields(f)
}

func (ds DistributionDefinitions) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) (err error) {
	for _, d := range ds {
		err = d.PrettyPrint(wr, detail)
		if err != nil {
			return
		}
	}
	return
}

type HardwareProfileDefinitions []HardwareProfileDefinition

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (ds HardwareProfileDefinitions) DefaultFields(f output.Format) string {
	return (HardwareProfileDefinition{}).DefaultFields(f)
}

func (ds HardwareProfileDefinitions) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) (err error) {
	for _, d := range ds {
		err = d.PrettyPrint(wr, detail)
		if err != nil {
			return
		}
	}
	return
}

type StorageGradeDefinitions []StorageGradeDefinition

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (ds StorageGradeDefinitions) DefaultFields(f output.Format) string {
	return (StorageGradeDefinition{}).DefaultFields(f)
}

func (ds StorageGradeDefinitions) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) (err error) {
	for _, d := range ds {
		err = d.PrettyPrint(wr, detail)
		if err != nil {
			return
		}
	}
	return
}

type ZoneDefinitions []ZoneDefinition

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (ds ZoneDefinitions) DefaultFields(f output.Format) string {
	return (ZoneDefinition{}).DefaultFields(f)
}

func (ds ZoneDefinitions) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) (err error) {
	for _, d := range ds {
		err = d.PrettyPrint(wr, detail)
		if err != nil {
			return
		}
	}
	return
}
