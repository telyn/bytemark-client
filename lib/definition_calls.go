package lib

import (
	"encoding/json"
	"fmt"
	"sort"
)

// Definitions represent all the possible things that can be returned as part of BigV's /definitions endpoint.
type Definitions struct {
	Distributions            []string
	StorageGrades            []string
	ZoneNames                []string
	DistributionDescriptions map[string]string
	StorageGradeDescriptions map[string]string
	HardwareProfiles         []string
	Keymaps                  []string
	Sendkeys                 []string
}

// JSONDefinition is an intermediate type used for converting BigV's JSON output for /definitions into the beautiful Definitions struct above. It should not be exported.
type JSONDefinition struct {
	ID   string          `json:"id"`
	Data json.RawMessage `json:"data"`
}

// JSONDefinitions should not be exported.
type JSONDefinitions []*JSONDefinition

// Process unmarshals the data from this JSONDefinition into the right field of the Definitions object.
func (d *JSONDefinition) Process(into *Definitions) error {
	switch d.ID {
	case "distributions":
		return json.Unmarshal(d.Data, &into.Distributions)
	case "storage_grades":
		return json.Unmarshal(d.Data, &into.StorageGrades)
	case "zone_names":
		return json.Unmarshal(d.Data, &into.ZoneNames)
	case "hardware_profiles":
		return json.Unmarshal(d.Data, &into.HardwareProfiles)
	case "keymaps":
		return json.Unmarshal(d.Data, &into.Keymaps)
	case "sendkeys":
		return json.Unmarshal(d.Data, &into.Sendkeys)
	case "distribution_descriptions":
		return json.Unmarshal(d.Data, &into.DistributionDescriptions)
	case "storage_grade_descriptions":
		return json.Unmarshal(d.Data, &into.StorageGradeDescriptions)
	}

	// Shouldn't be a fatal error. into may still be useful
	return fmt.Errorf("Unknown definition returned: %v", d.ID)

}

// Process processes our intermediate JSONDefinitions into a Definitions object.
func (defs JSONDefinitions) Process() *Definitions {
	var out Definitions

	for _, def := range defs {
		if err := def.Process(&out); err != nil {
			fmt.Printf("WARN: %v\n", err)
		}
	}
	sort.StringSlice(out.Distributions).Sort()
	sort.StringSlice(out.StorageGrades).Sort()
	sort.StringSlice(out.ZoneNames).Sort()
	sort.StringSlice(out.HardwareProfiles).Sort()

	return &out
}

// ReadDefinitions queries the brain for its definitions
func (c *bytemarkClient) ReadDefinitions() (definitions *Definitions, err error) {
	var defs JSONDefinitions
	r, err := c.BuildRequestNoAuth("GET", BrainEndpoint, "/definitions")
	if err != nil {
		return
	}
	_, _, err = r.Run(nil, &defs)
	definitions = defs.Process()
	return
}
