package brain

import (
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
	"io"
)

func definitionDefaultFields(f output.Format) string {
	switch f {
	case output.List:
		return "Name, Description"
	}
	return "Name, Description"
}

func definitionPrettyPrint(definition interface{}, wr io.Writer, detail prettyprint.DetailLevel) error {
	defTpl := `
	{{ define "definition_sgl" }}{{ .Name }}: {{ .Description }}{{ end }}
	{{ define "definition_medium" }}{{ template "definition_sgl" . }}{{ end }}
	{{ define "definition_full" }}{{ template "definition_medium" . }}{{ end }}
	`
	return prettyprint.Run(wr, defTpl, "definition"+string(detail), definition)
}

// DistributionDefinition is an object we assemble from distributions and distribution_descriptions from the /definitions API call
// in the future (bytemark-client 3.0?) a slice of these this will replace the Definitions.Distributions slice and Definitions.DistributionDescriptions map.
type DistributionDefinition struct {
	Name        string
	Description string
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/bytemark-client for this type.
func (d DistributionDefinition) DefaultFields(f output.Format) string {
	return definitionDefaultFields(f)
}

// PrettyPrint outputs a vaguely human-readable version of the definition to wr. Detail is ignored.
func (d DistributionDefinition) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	return definitionPrettyPrint(d, wr, detail)
}

// HardwareProfileDefinition is an object we assemble from hardwareprofiles from the /definitions API call
// in the future (bytemark-client 3.0?) a slice of these this will replace the Definitions.HardwareProfiles slice.
type HardwareProfileDefinition struct {
	Name        string
	Description string
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/bytemark-client for this type.
func (hp HardwareProfileDefinition) DefaultFields(f output.Format) string {
	return definitionDefaultFields(f)
}

// PrettyPrint outputs a vaguely human-readable version of the definition to wr. Detail is ignored.
func (hp HardwareProfileDefinition) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	return definitionPrettyPrint(hp, wr, detail)
}

// StorageGradeDefinition is an object we assemble from storage_grades and storage_grade_descriptions from the /definitions API call
// in the future (bytemark-client 3.0?) a slice of these this will replace the Definitions.StorageGrades slice and Definitions.StorageGradeDescriptions map.
type StorageGradeDefinition struct {
	Name        string
	Description string
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/bytemark-client for this type.
func (sg StorageGradeDefinition) DefaultFields(f output.Format) string {
	return definitionDefaultFields(f)
}

// PrettyPrint outputs a vaguely human-readable version of the definition to wr. Detail is ignored.
func (sg StorageGradeDefinition) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	return definitionPrettyPrint(sg, wr, detail)
}

// ZoneDefinition is an object we assemble from zone_names from the /definitions API call and some static data in lib/definitions.go
// in the future (bytemark-client 3.0?) a slice of these this will replace the Definitions.ZoneNames slice.
type ZoneDefinition struct {
	Name        string
	Description string
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/bytemark-client for this type.
func (z ZoneDefinition) DefaultFields(f output.Format) string {
	return definitionDefaultFields(f)
}

// PrettyPrint outputs a vaguely human-readable version of the definition to wr. Detail is ignored.
func (z ZoneDefinition) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	return definitionPrettyPrint(z, wr, detail)
}
