package brain

import "fmt"

type VLAN struct {
	ID        int    `json:"id"`
	Num       int    `json:"num"`
	UsageType string `json:"usage_type"`
}

func (v *VLAN) String() string {
	return fmt.Sprintf("%d: %s (Num: %d)", v.ID, v.UsageType, v.ID)
}
