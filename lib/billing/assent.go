package billing

// Assent represents an assent by a person on behalf of an account to an agreement
type Assent struct {
	AgreementID string `json:"-"`
	AccountID   int    `json:"account_id"`
	PersonID    int    `json:"person_id"`
	// Name is the full name of the person
	Name  string `json:"name"`
	Email string `json:"email"`
}
