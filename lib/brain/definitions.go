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
	default: // also output.Table
		return "Name, Description"
	}
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

func (d DistributionDefinition) DefaultFields(f output.Format) string {
	return definitionDefaultFields(f)
}

func (d DistributionDefinition) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	return definitionPrettyPrint(d, wr, detail)
}

// HardwareProfileDefinition is an object we assemble from hardwareprofiles in the /*Definitions API call and some static data in lib/definitions.go
// in the future (bytemark-client 3.0?) a slice of these this will replace the Definitions.HardwareProfiles slice.
type HardwareProfileDefinition struct {
	Name        string
	Description string
}

func (hp HardwareProfileDefinition) DefaultFields(f output.Format) string {
	return definitionDefaultFields(f)
}

func (hp HardwareProfileDefinition) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	return definitionPrettyPrint(hp, wr, detail)
}

// StorageGradeDefinition is an object we assemble from storage_grades and storage_grade_descriptions in the /*Definitions API call
// in the future (bytemark-client 3.0?) a slice of these this will replace the Definitions.StorageGrades slice and Definitions.StorageGradeDescriptions map.
type StorageGradeDefinition struct {
	Name        string
	Description string
}

func (sg StorageGradeDefinition) DefaultFields(f output.Format) string {
	return definitionDefaultFields(f)
}

func (sg StorageGradeDefinition) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	return definitionPrettyPrint(sg, wr, detail)
}

// ZoneDefinition is an object we assemble from zone_names in the /*Definitions API call and some static data in lib/definitions.go
// in the future (bytemark-client 3.0?) a slice of these this will replace the Definitions.ZoneNames slice.
type ZoneDefinition struct {
	Name        string
	Description string
}

func (z ZoneDefinition) DefaultFields(f output.Format) string {
	return definitionDefaultFields(f)
}

func (z ZoneDefinition) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	return definitionPrettyPrint(z, wr, detail)
}
