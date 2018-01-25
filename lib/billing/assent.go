package billing

// Assent represents an assent by a person on behalf of an account to an agreement
type Assent struct {
	AgreementID string
	AccountID   int
	PersonID    int
	// Name is the full name of the person
	Name string
}
