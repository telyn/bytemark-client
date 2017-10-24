package billing

// Definition is an admin-modifiable parameter for bmbilling
// examples include account-opening credit amount and trial length
type Definition struct {
	ID    int    `json:"id,omitempty"`
	Name  string `json:"name"`
	Value string `json:"value"`
	// Which auth group a user must be in to update the definition
	UpdateGroupReq string `json:"update_group_req,omitempty"`
}
