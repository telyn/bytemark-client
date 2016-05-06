package lib

import (
	"encoding/json"
	"fmt"
)

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

type JSONDefinition struct {
	ID   string          `json:"id"`
	Data json.RawMessage `json:"data"`
}

type JSONDefinitions []*JSONDefinition

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

func (defs JSONDefinitions) Process() *Definitions {
	var out Definitions

	for _, def := range defs {
		if err := def.Process(&out); err != nil {
			fmt.Printf("WARN: %v\n", err)
		}
	}

	return &out
}

func (c *bytemarkClient) ReadDefinitions() (definitions *Definitions, err error) {
	var defs JSONDefinitions
	r, err := c.BuildRequestNoAuth("GET", EP_BRAIN, "/definitions")

	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &defs)

	definitions = defs.Process()

	return

}
