package config

import "strings"

// Var is a struct which contains a name-value-source triplet
// Source is up to two words separated by a space. The first word is the source type: FLAG, ENV, DIR, CODE.
// The second is the name of the flag/file/environment var used.
type Var struct {
	Name   string
	Value  string
	Source string
}

// SourceType returns one of the following:
// FLAG for a configVar whose value was set by passing a flag on the command line
// ENV for a configVar whose value was set from an environment variable
// DIR for a configVar whose value was set from a file in the config dir
//
func (v *Var) SourceType() string {
	bits := strings.Fields(v.Source)

	return bits[0]
}

// SourceBaseName returns the basename of the configVar's source.
// it's a bit stupid and so its output is only valid for configVars with SourceType() of DIR
func (v *Var) SourceBaseName() string {
	bits := strings.Split(v.Source, "/")
	return bits[len(bits)-1]
}
