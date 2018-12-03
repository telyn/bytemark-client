package brain

type APIKey struct {
	ID        int    `json:"id,omitempty"`
	UserID    int    `json:"user_id,omitempty"`
	Label     string `json:"label,omitempty"`
	APIKey    string `json:"api_key,omitempty"`
	ExpiresAt string `json:"expires_at,omitempty"`
}
