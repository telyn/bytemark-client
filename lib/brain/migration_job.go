package brain

import (
	"encoding/json"
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
	"github.com/BytemarkHosting/bytemark-client/lib/util"
)

//MigrationJobModification represents the modifications possible on a migration job
type MigrationJobModification struct {
	Cancel  MigrationJobLocations `json:"cancel,omitempty"`
	Options MigrationJobOptions   `json:"options,omitempty"`
}

// MigrationJobQueue is a list of disc IDs that are still to be migrated as
// part of a migration job.
type MigrationJobQueue struct {
	Discs []int `json:"discs,omitempty"`
}

// MigrationJobLocations represents source or target locations for a migration
// job. Discs, pools and tails maybe represented by ID number, label or UUID.
type MigrationJobLocations struct {
	Discs []util.NumberOrString `json:"discs,omitempty"`
	Pools []util.NumberOrString `json:"pools,omitempty"`
	Tails []util.NumberOrString `json:"tails,omitempty"`
}

// MigrationJobDestinations represents available desintations for a migration
// job. Unlike MigrationJobLocations, these are represented using ID number
// only.
type MigrationJobDestinations struct {
	Pools []int `json:"pools,omitempty"`
}

// MigrationJobOptions represents options on a migration job.
type MigrationJobOptions struct {
	Priority int `json:"priority,omitempty"`
}

// MigrationJobDiscStatus represents the current status of a migration job.
// Each entry is a list of disc IDs indicating the fate of discs that
// have been removed from the queue.
type MigrationJobDiscStatus struct {
	Done      []int `json:"done,omitempty"`
	Errored   []int `json:"errored,omitempty"`
	Cancelled []int `json:"cancelled,omitempty"`
	Skipped   []int `json:"skipped,omitempty"`
}

// MigrationJobStatus captures the status of a migration job, currently only
// discs.
type MigrationJobStatus struct {
	Discs MigrationJobDiscStatus `json:"discs,omitempty"`
}

// MigrationJobSpec is a specification of a migration job to be created
type MigrationJobSpec struct {
	Options      MigrationJobOptions   `json:"options,omitempty"`
	Sources      MigrationJobLocations `json:"sources,omitempty"`
	Destinations MigrationJobLocations `json:"destinations,omitempty"`
}

// MigrationJob is a representation of a migration job.
type MigrationJob struct {
	ID           int                      `json:"id,omitempty"`
	Args         MigrationJobSpec         `json:"args,omitempty"`
	Queue        MigrationJobQueue        `json:"queue,omitempty"`
	Destinations MigrationJobDestinations `json:"destinations,omitempty"`
	Status       MigrationJobStatus       `json:"status,omitempty"`
	Priority     int                      `json:"priority,omitempty"`
	StartedAt    string                   `json:"started_at,omitempty"`
	FinishedAt   string                   `json:"finished_at,omitempty"`
	CreatedAt    string                   `json:"created_at,omitempty"`
	UpdatedAt    string                   `json:"updated_at,omitempty"`
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (mj MigrationJob) DefaultFields(f output.Format) string {
	return "ID, Queue, Destinations, Status, Priority, StartedAt, FinishedAt, CreatedAt, UpdatedAt"
}

// PrettyPrint formats a MigrationJobQueue for display
func (mjq MigrationJobQueue) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	const template = `{{ define "migration_job_queue_full" }}
{{- with .Discs }}
     discs:
{{- range . }}
       • {{.}}
{{- end }}
{{- end -}}
{{- end -}}{{- define "migration_job_queue_sgl" -}}{{ range .Discs }} {{.}}{{ end }}{{- end -}}`
	return prettyprint.Run(wr, template, "migration_job_queue"+string(detail), mjq)
}

// PrettyPrint formats a MigrationJobStatus for display
func (mjs MigrationJobStatus) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	const template = `{{- define "migration_job_status_full" -}}
{{- with .Discs -}}
{{- with .Done }}
     done:
    {{- range . }}
       • {{.}}
    {{- end }}
{{- end -}}
{{- with .Errored }}
     errored:
    {{- range . }}
       • {{.}}
    {{- end }}
{{- end -}}
{{- with .Cancelled }}
     cancelled:
    {{- range . }}
       • {{.}}
    {{- end }}
{{- end -}}
{{- with .Skipped }}
     skipped:
    {{- range . }}
       • {{.}}
    {{- end }}
{{- end -}}
{{- end -}}
{{- end -}}`
	return prettyprint.Run(wr, template, "migration_job_status"+string(detail), mjs)
}

// PrettyPrint outputs a nice human-readable overview of the migration
func (mj MigrationJob) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	const template = `{{ define "migration_job_full" }} ▸ {{ .ID }}
{{ with .Queue }}   queue: {{ prettysprint . "_full" }}
{{ end -}}{{- with .Status }}   status: {{ prettysprint . "_full" }}
{{ end -}}{{- with .Priority }}   priority: {{ . }}
{{ end -}}{{- with .StartedAt }}   started_at: {{ . }}
{{ end -}}{{- with .FinishedAt }}   finished_at: {{ . }}
{{ end -}}{{- with .CreatedAt }}   created_at: {{ . }}
{{ end -}}{{- with .UpdatedAt }}   updated_at: {{ . }}
{{ end -}}{{- end -}}{{- define "migration_job_medium" }} ▸ {{ .ID }} queue:{{ prettysprint .Queue "_sgl" }}
{{- end -}}{{- define "migration_job_sgl" -}}{{ template "migration_job_medium" . }}{{- end -}}`
	return prettyprint.Run(wr, template, "migration_job"+string(detail), mj)
}

// DefaultFields appeases quality tests
func (x MigrationJobDestinations) DefaultFields(f output.Format) string {
	return ""
}

// PrettyPrint appeases quality tests
func (x MigrationJobDestinations) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	return nil
}

// DefaultFields appeases quality tests
func (x MigrationJobDiscStatus) DefaultFields(f output.Format) string {
	return ""
}

// PrettyPrint appeases quality tests
func (x MigrationJobDiscStatus) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	return nil
}

// DefaultFields appeases quality tests
func (x MigrationJobLocations) DefaultFields(f output.Format) string {
	return ""
}

// PrettyPrint appeases quality tests
func (x MigrationJobLocations) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	return nil
}

// DefaultFields appeases quality tests
func (x MigrationJobOptions) DefaultFields(f output.Format) string {
	return ""
}

// PrettyPrint appeases quality tests
func (x MigrationJobOptions) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	return nil
}

// DefaultFields appeases quality tests
func (mjq MigrationJobQueue) DefaultFields(f output.Format) string {
	return ""
}

// DefaultFields appeases quality tests
func (x MigrationJobSpec) DefaultFields(f output.Format) string {
	return ""
}

// PrettyPrint appeases quality tests
func (x MigrationJobSpec) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	return nil
}

// DefaultFields appeases quality tests
func (mjs MigrationJobStatus) DefaultFields(f output.Format) string {
	return ""
}
