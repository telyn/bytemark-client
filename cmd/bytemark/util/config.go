package util

import (
	"flag"
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

var configVars = [...]string{
	"endpoint",
	"billing-endpoint",
	"auth-endpoint",
	"spp-endpoint",
	"admin",
	"user",
	"account",
	"group",
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

// InvalidConfigVarError is used to inform the user that they variable they attempted to set / get doesn't exist
type InvalidConfigVarError struct {
	ConfigVar string
}

func (e InvalidConfigVarError) Error() string {
	vs := "'" + strings.Join(configVars[:], "','") + "'"
	return fmt.Sprintf("'%s' is not a valid config var. Valid config vars are: %s", e.ConfigVar, vs)
}

// ConfigVar is a struct which contains a name-value-source triplet
// Source is up to two words separated by a space. The first word is the source type: FLAG, ENV, DIR, CODE.
// The second is the name of the flag/file/environment var used.
type ConfigVar struct {
	Name   string
	Value  string
	Source string
}

// SourceType returns one of the following:
// FLAG for a configVar whose value was set by passing a flag on the command line
// ENV for a configVar whose value was set from an environment variable
// DIR for a configVar whose value was set from a file in the config dir
//
func (v *ConfigVar) SourceType() string {
	bits := strings.Fields(v.Source)

	return bits[0]
}

// SourceBaseName returns the basename of the configVar's source.
// it's a bit stupid and so its output is only valid for configVars with SourceType() of DIR
func (v *ConfigVar) SourceBaseName() string {
	bits := strings.Split(v.Source, "/")
	return bits[len(bits)-1]
}

// ConfigManager is an interface defining a key->value store that also knows where the values were set from.
type ConfigManager interface {
	Get(string) (string, error)
	GetIgnoreErr(string) string
	GetBool(string) (bool, error)
	GetV(string) (ConfigVar, error)
	GetVirtualMachine() *lib.VirtualMachineName
	GetGroup() *lib.GroupName
	GetAll() ([]ConfigVar, error)
	Set(string, string, string)
	SetPersistent(string, string, string) error
	Unset(string) error
	GetDebugLevel() int
	EndpointName() string
	PanelURL() string

	ImportFlags(*flag.FlagSet) []string
}

// Params currently used:
// token - an OAuth 2.0 bearer token to use when authenticating
// username - the default username to use - if not present, $USER
// endpoint - the default endpoint to use - if not present, https://uk0.bigv.io
// auth-endpoint - the default auth API endpoint to use - if not present, https://auth.bytemark.co.uk
// account - account to use if not specified elsewhere§
// group - group to use if not specified

// A Config determines the configuration of the Bytemark client.
// It's responsible for handling things like the credentials to use and what endpoints to talk to.
//
// Each configuration item is read from the following places, falling back to successive places:
//
// Per-command command-line flags, global command-line flags, environment variables, configuration directory, hard-coded defaults
//
//The location of the configuration directory is read from global command-line flags, or is otherwise ~/.bytemark
//
type Config struct {
	debugLevel  int
	Dir         string
	Memo        map[string]ConfigVar
	Definitions map[string]string
}

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

// NewConfig sets up a new config struct. Pass in an empty string to default to ~/.bytemark
func NewConfig(configDir string) (config *Config, err error) {
	config = new(Config)
	config.Memo = make(map[string]ConfigVar)
	home := os.Getenv("HOME")
	if runtime.GOOS == "windows" {
		home = os.Getenv("APPDATA")

	}

	config.Dir = filepath.Join(home, "/.bytemark")
	if os.Getenv("BM_CONFIG_DIR") != "" {
		config.Dir = os.Getenv("BM_CONFIG_DIR")
	}

	if configDir != "" {
		config.Dir = strings.Replace(configDir, "~", os.Getenv("HOME"), -1)
	}

	err = os.MkdirAll(config.Dir, 0700)
	if err != nil {

		return nil, err
	}

	stat, err := os.Stat(config.Dir)
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return nil, &ConfigDirInvalidError{config.Dir}
	}

	dbgLog := config.GetPath("debug.log")
	_, err = os.Stat(dbgLog)
	if err == nil {
		_, err1 := os.Stat(dbgLog + ".1")
		if err1 == nil {
			os.Remove(dbgLog + ".1")
		}
		os.Rename(dbgLog, dbgLog+".1")
	}

	log.LogFile, err = os.Create(dbgLog)
	if err != nil {
		log.Errorf("Couldn't open %s for writing\r\n", config.GetPath("debug.log"))
	}

	strDL, err := config.Get("debug-level")
	if err != nil {
		log.Error(err)
	} else {
		debugLevel, err := strconv.ParseInt(strDL, 10, 0)
		if err == nil {
			config.debugLevel = int(debugLevel)
			log.DebugLevel = int(debugLevel)
		}
	}
	return config, nil
}

// ImportFlags reads all the flags from the passed FlagSet that have the same name as a valid configVar, and sets the configVar to that.
func (config *Config) ImportFlags(flags *flag.FlagSet) []string {
	if flags != nil {
		if flags.Parsed() {
			// dump all the flags into the memo
			// should be reet...reet?
			flags.Visit(func(f *flag.Flag) {
				config.Memo[f.Name] = ConfigVar{
					f.Name,
					f.Value.String(),
					"FLAG " + f.Name,
				}
			})

			strDL := config.GetIgnoreErr("debug-level")
			debugLevel, err := strconv.ParseInt(strDL, 10, 0)
			if err == nil {
				config.debugLevel = int(debugLevel)
				log.DebugLevel = int(debugLevel)
			}

			return flags.Args()
		}
	}
	return nil
}

// GetDebugLevel returns the current debug-level as an integer. This is used throughout the github.com/BytemarkHosting/bytemark-client library to determine verbosity of output.
func (config *Config) GetDebugLevel() int {
	return config.debugLevel
}

// GetPath joins the given string onto the end of the Config.Dir path
func (config *Config) GetPath(name string) string {
	return filepath.Join(config.Dir, name)
}

// Get returns the value of a ConfigVar. Used to simplify code when the source is unnecessary.
func (config *Config) Get(name string) (string, error) {
	v, err := config.GetV(name)
	return v.Value, err
}

// GetIgnoreErr returns the value of a ConfigVar or an empty string , if it was unable to read it for whatever reason.
func (config *Config) GetIgnoreErr(name string) string {
	s, _ := config.Get(name)
	return s
}

// GetV returns the ConfigVar for the given key.
func (config *Config) GetV(name string) (ConfigVar, error) {
	// try to read the Memo
	name = strings.ToLower(name)
	if val, ok := config.Memo[name]; ok {
		if val.Name == "" {
			val.Name = name
		}
		if val.Source == "" {
			val.Source = "SOURCE UNSET"
		}
		return val, nil
	}
	return config.read(name)
}

// GetVirtualMachine returns a VirtualMachineName with the config's default group and account set, and a blank VirtualMachine field
func (config *Config) GetVirtualMachine() (vm *lib.VirtualMachineName) {
	vm = new(lib.VirtualMachineName)
	vm.Account = config.GetIgnoreErr("account")
	vm.Group = config.GetIgnoreErr("group")
	vm.VirtualMachine = ""
	return vm
}

// GetGroup returns a GroupName with the config's default group and account
func (config *Config) GetGroup() (group *lib.GroupName) {
	group = new(lib.GroupName)
	group.Account = config.GetIgnoreErr("account")
	group.Group = config.GetIgnoreErr("group")
	return group
}

// GetAll returns all of the available ConfigVars in the Config.
func (config *Config) GetAll() (vars []ConfigVar, err error) {
	vars = make([]ConfigVar, len(configVars))
	for i, v := range configVars {
		vars[i], err = config.GetV(v)
		if err != nil {
			return nil, err
		}
	}
	return vars, nil
}

// GetDefault returns the default ConfigVar for the given key.
func (config *Config) GetDefault(name string) ConfigVar {
	// ideally most of these should just be	os.Getenv("BM_"+name.Upcase().Replace("-","_"))
	switch name {
	case "user":
		if os.Getenv("BM_USER") == "" {
			return ConfigVar{"user", os.Getenv("USER"), "ENV USER"}
		}
		return ConfigVar{"user", os.Getenv("BM_USER"), "ENV BM_USER"}
	case "endpoint":
		v := ConfigVar{"endpoint", "https://uk0.bigv.io", "CODE"}

		val := os.Getenv("BM_ENDPOINT")
		if val != "" {
			v.Value = val
			v.Source = "ENV BM_ENDPOINT"
		}
		return v
	case "billing-endpoint":
		if config.GetIgnoreErr("endpoint") == "https://staging.bigv.io" {
			return ConfigVar{"billing-endpoint", "", "CODE STAGING DEFAULT"}
		}
		v := ConfigVar{"billing-endpoint", "https://bmbilling.bytemark.co.uk", "CODE"}
		if val := os.Getenv("BM_BILLING_ENDPOINT"); val != "" {
			v.Value = val
			v.Source = "ENV BM_BILLING_ENDPOINT"
		}
		return v
	case "spp-endpoint":
		v := ConfigVar{"spp-endpoint", "https://spp-submissions.bytemark.co.uk", "CODE"}
		if val := os.Getenv("BM_SPP_ENDPOINT"); val != "" {
			v.Value = val
			v.Source = "ENV BM_SPP_ENDPOINT"
		}
		return v
	case "auth-endpoint":
		v := ConfigVar{"auth-endpoint", "https://auth.bytemark.co.uk", "CODE"}

		val := os.Getenv("BM_AUTH_ENDPOINT")
		if val != "" {
			v.Value = val
			v.Source = "ENV BM_AUTH_ENDPOINT"
		}
		return v
	case "account":
		val := os.Getenv("BM_ACCOUNT")
		if val != "" {
			return ConfigVar{
				"account",
				val,
				"ENV BM_ACCOUNT",
			}
		}
		return ConfigVar{
			"account",
			"",
			"CODE",
		}
	case "group":
		val := os.Getenv("BM_GROUP")
		if val != "" {
			return ConfigVar{
				"group",
				val,
				"ENV BM_GROUP",
			}
		}
		return ConfigVar{"group", "default", "CODE"}
	case "debug-level":
		v := ConfigVar{"debug-level", "0", "CODE"}
		if val := os.Getenv("BM_DEBUG_LEVEL"); val != "" {
			v.Value = val
		}
		return v
	case "force":
		return ConfigVar{"force", "false", "CODE"}
	}
	return ConfigVar{name, "", "UNSET"}
}

// GetBool returns the given configvar as a bool - true if it is set, not blank, and not equal to "false". false otherwise.
func (config *Config) GetBool(name string) (bool, error) {
	v, err := config.Get(name)
	if err != nil {
		return false, err
	}
	return !(v == "" || v == "false"), nil
}

func (config *Config) read(name string) (ConfigVar, error) {
	path := config.GetPath(name)
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return config.GetDefault(name), nil
		}

		return config.GetDefault(name), &ConfigReadError{Name: name, Path: path, Err: err}
	}

	return ConfigVar{name, strings.TrimSpace(string(contents)), "FILE " + path}, nil
}

// Set stores the given key-value pair in config's Memo. This storage does not persist once the program terminates.
func (config *Config) Set(name, value, source string) {
	config.Memo[name] = ConfigVar{name, value, source}
}

// SetPersistent writes a file to the config directory for the given key-value pair.
func (config *Config) SetPersistent(name, value, source string) error {
	found := false
	for _, v := range configVars {
		if v == name {
			found = true
		}
	}
	if !found {
		return InvalidConfigVarError{name}
	}
	path := config.GetPath(name)
	config.Set(name, value, source)
	err := ioutil.WriteFile(path, []byte(value), 0600)
	if err != nil {
		return &ConfigWriteError{Name: name, Path: path, Err: err}
	}
	return nil
}

// Unset removes the named key from both config's Memo and the user's config directory.
func (config *Config) Unset(name string) error {
	found := false
	for _, v := range configVars {
		if v == name {
			found = true
		}
	}
	if !found {
		return InvalidConfigVarError{name}
	}
	delete(config.Memo, name)
	return os.Remove(config.GetPath(name))
}

// PanelURL returns config's best guess at the correct URL for the bytemark panel for the cluster with the endpoint we're using. Basically it flips between panel-beta.bytemark and panel-int.
func (config *Config) PanelURL() string {
	endpoint := config.EndpointName()
	if strings.EqualFold(endpoint, "uk0.bigv.io") {
		return "https://panel-beta.bytemark.co.uk"
	}
	if strings.EqualFold(endpoint, "int.bigv.io") {
		// am i leaking a secret?
		return "https://panel-int.vlan863.bytemark.uk0.bigv.io"
	}
	panel := config.GetIgnoreErr("panel-address")
	if panel == "" {
		panel = "https://your.panel.address.example.com"
	}
	return panel
}

// EndpointName trims the URL scheme off the beginning of the endpoint.
// TODO(telyn): Why?
func (config *Config) EndpointName() string {
	endpoint := config.GetIgnoreErr("endpoint")
	endpoint = strings.TrimPrefix(endpoint, "https://")
	endpoint = strings.TrimPrefix(endpoint, "http://") // it never hurts to be prepared
	return endpoint
}
