package brain

// AdminNote represents the request to create a note on a sotrage pool or head
type AdminNote struct {
	On   string `json:"on"`
	Spec string `json:"spec"`
	Note string `json:"note"`
}
