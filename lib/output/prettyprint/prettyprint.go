package prettyprint

import (
	"io"
	"text/template"
)

// DetailLevel allows us to specify how much detail we want included when we call PrettyPrint
type DetailLevel string

const (
	// SingleLine will show you only as much detail as fits on a single good line.
	SingleLine DetailLevel = "_sgl"
	// Medium is a small amount of detail, but that runs on more than one line. Most likely
	// it will be two lines long. I just didn't want to guarantee that.
	Medium = "_medium"
	// Full is the usual level of detail for a 'show' command. A nice, multi-line thing
	// that contains all the detail a normal person could ever need.
	Full = "_full"
	// Ridiculous will include all the detail probably no one cares about.
	// For example: What head a VM is on, what storage pool a disc is in, every note on a head
	//Ridiculous = "_ridiculous"

	// List is a very minimal amount of information, good for scripting against. In most cases it will be a name. Sometimes followed by a bit of status information, after a space.
	// On collection types, such as Backups (rather than Backup)
	List = "_list"
)

// PrettyPrinter is the common interface used to output different entities in a user friendly way
type PrettyPrinter interface {
	PrettyPrint(wr io.Writer, detail DetailLevel) error
}

// Run is a convenience function for running templates with the standard prettyprint functions
func Run(wr io.Writer, templates string, templateToExecute string, object interface{}) error {
	tmpl, err := template.New("prettyprint").Funcs(templateFuncMap).Parse(templates)
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(wr, templateToExecute, object)
}
