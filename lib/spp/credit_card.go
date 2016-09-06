package spp

// CreditCard represents all the data for a credit card with SPP.
type CreditCard struct {
	Number   string `yaml:"account_number"`
	Name     string `yaml:"name"`
	Expiry   string `yaml:"expiry"`
	CVV      string `yaml:"cvv"`
	Street   string `yaml:"street,omitempty"`
	City     string `yaml:"city,omitempty"`
	Postcode string `yaml:"postcode,omitempty"`
	Country  string `yaml:"country,omitempty"`
}
