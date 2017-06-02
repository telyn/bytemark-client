package brain

// DistributionDefinition is an object we assemble from distributions and distribution_descriptions from the /definitions API call
// in the future (bytemark-client 3.0?) a slice of these this will replace the Definitions.Distributions slice and Definitions.DistributionDescriptions map.
type DistributionDefinition struct {
	Name        string
	Description string
}

// HardwareProfileDefinition is an object we assemble from hardwareprofiles in the /definitions API call and some static data in lib/definitions.go
// in the future (bytemark-client 3.0?) a slice of these this will replace the Definitions.HardwareProfiles slice.
type HardwareProfileDefinition struct {
	Name        string
	Description string
}

// StorageGradeDefinition is an object we assemble from storage_grades and storage_grade_descriptions in the /definitions API call
// in the future (bytemark-client 3.0?) a slice of these this will replace the Definitions.StorageGrades slice and Definitions.StorageGradeDescriptions map.
type StorageGradeDefinition struct {
	Name        string
	Description string
}

// ZoneDefinition is an object we assemble from zone_names in the /definitions API call and some static data in lib/definitions.go
// in the future (bytemark-client 3.0?) a slice of these this will replace the Definitions.ZoneNames slice.
type ZoneDefinition struct {
	Name        string
	Description string
}
