package config

import (
	"fmt"
	"strings"
)

// ConfigDirInvalidError is returned when the path specified as the config dir was not a directory.
type ConfigDirInvalidError struct {
	Path string
}

func (e *ConfigDirInvalidError) Error() string {
	return fmt.Sprintf("The config directory is '%s' but it doesn't seem to be a directory.", e.Path)
}

// CannotLoadDefinitionsError is unused. Planned to be used if bytemark-client starts caching definitions, but it doesn't at the moment.
type CannotLoadDefinitionsError struct {
	Err error
}

func (e *CannotLoadDefinitionsError) Error() string {
	return fmt.Sprintf("Unable to load the definitions file from the Bytemark API.")
}

// ConfigReadError is returned when a file containing a value for a configVar couldn't be read.
type ConfigReadError struct {
	Name string
	Path string
	Err  error
}

func (e *ConfigReadError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Unable to read config for %s from %s - %v", e.Name, e.Path, e.Err)
	}
	return fmt.Sprintf("Unable to read config for %s from %s.", e.Name, e.Path)
}

// ConfigWriteError is returned when a file containing a value for a configVar couldn't be written to.
type ConfigWriteError struct {
	Name string
	Path string
	Err  error
}

func (e *ConfigWriteError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Unable to write config for %s to %s.", e.Name, e.Path)
	}
	return fmt.Sprintf("Unable to write config for %s to %s - %s", e.Name, e.Path, e.Err.Error())
}

// InvalidVarError is used to inform the user that they variable they attempted to set / get doesn't exist
type InvalidVarError struct {
	ConfigVar string
}

func (e InvalidVarError) Error() string {
	vs := "'" + strings.Join(configVars[:], "','") + "'"
	return fmt.Sprintf("'%s' is not a valid config var. Valid config vars are: %s", e.ConfigVar, vs)
}
