package config

var configVars = [...]string{
	"endpoint",
	"billing-endpoint",
	"auth-endpoint",
	"spp-endpoint",
	"admin",
	"user",
	"account",
	"group",
	"output-format",
	"session-validity",
	"token",
	"debug-level",
	"yubikey",
}

// IsConfigVar checks to see if the named variable is actually one of the settable configVars.
func IsConfigVar(name string) bool {
	for _, v := range configVars {
		if v == name {
			return true
		}
	}
	return false
}
