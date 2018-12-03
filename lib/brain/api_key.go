package brain

// APIKey represents an api_key in the brain.
type APIKey struct {
	// ID must not be set during creation.
	ID     int `json:"id,omitempty"`
	UserID int `json:"user_id,omitempty"`
	// Label is a friendly display name for this API key
	Label string `json:"label,omitempty"`
	// API key is the actual key. To use it, it must be prepended with
	// 'apikey.' in the HTTP Authorization header. For example, if the api key
	// is xpq21, the HTTP headers should include `Authorization: Bearer apikey.xpq21`
	APIKey string `json:"api_key,omitempty"`
	// ExpiresAt should be a time or datetime in HH:MM:SS or
	// YYYY-MM-DDTHH:MM:SS.msZ where T is a literal T, .ms are optional
	// microseconds and Z is either literal Z (meaning UTC) or a timezone
	// specified like -600 or +1200
	ExpiresAt string `json:"expires_at,omitempty"`
	// Privileges cannot be set at creation or update time, but are returned by
	// the brain when view=overview.
	Privileges Privileges `json:"privileges,omitempty"`
}
