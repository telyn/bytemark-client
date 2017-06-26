package billing

// Person represents a bmbilling person
type Person struct {
	ID int `json:"id,omitempty"`
	// Username is the name this person uses to log in to our services.
	Username    string `json:"username"`
	Email       string `json:"email"`
	BackupEmail string `json:"email_backup,omitempty"`

	// only set in the creation request
	Password string `json:"password"`

	FirstName   string `json:"firstname"`
	LastName    string `json:"surname"`
	Address     string `json:"address"`
	City        string `json:"city"`
	StateCounty string `json:"statecounty,omitempty"`
	Postcode    string `json:"postcode"`
	Country     string `json:"country"`
	Phone       string `json:"phone"`
	MobilePhone string `json:"phonemobile,omitempty"`

	Organization         string `json:"organization,omitempty"`
	OrganizationDivision string `json:"division,omitempty"`
	VATNumber            string `json:"vatnumber,omitempty"`
}

func (p Person) IsValid() bool {
	return p.Username != ""
}
