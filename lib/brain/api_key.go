package brain

import (
	"io"
	"time"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

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

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (key APIKey) DefaultFields(f output.Format) string {
	switch f {
	case output.List:
		return "ID, Label, Expired, ExpiresAt, Privileges"
	}
	return "ID, UserID, Label, Expired, ExpiresAt, Privileges"
}

// Expired returns true if ExpiresAt is in the past.
// This assumes ExpiresAt is in ISO8601 format, which it will be if the APIKey
// was set by a unmarshalling a response from the brain.
func (key APIKey) Expired() bool {
	if key.ExpiresAt == "" {
		return false
	}
	// TODO(telyn): not keen on this ad-hoc iso8601 parsing
	expiresAt, err := time.Parse("2006-01-02T15:04:05-0700", key.ExpiresAt)
	if err != nil {
		return false
	}
	return expiresAt.Before(time.Now())
}

// PrettyPrint outputs a nice human-readable overview of the api key, including
// what privileges it has, to the given writer.
func (key APIKey) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	accountTpl := `
{{ define "apikey_sgl" }}{{ .Label }}{{ if .Expired }} (expired){{ end }}{{ end }}
{{ define "apikey_medium" }}{{ template "apikey_sgl" . }}{{ end }}
{{ define "apikey_full" }}{{ template "apikey_sgl" . }}
  Expire{{ if .Expired }}d{{ else }}s{{ end }}: {{ if eq .ExpiresAt "" }}never{{ else }}{{ .ExpiresAt }}{{ end -}}
{{ if .APIKey}}
  Key: apikey.{{ .APIKey }}{{ end -}}
{{ if len .Privileges | le 1 }}

  Privileges:
{{ prettysprint .Privileges "_medium" | prefixEachLine "* " | indent 4 -}}
{{ else }}
{{ end -}}
{{ end }}
`
	return prettyprint.Run(wr, accountTpl, "apikey"+string(detail), key)
}
