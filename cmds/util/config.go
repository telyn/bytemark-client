package util

import (
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/util/log"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

var configVars = [...]string{
	"endpoint",
	"auth-endpoint",
	"user",
	"account",
	"group",
	"token",
	"debug-level",
	"yubikey",
}

// ConfigVar is a struct which contains a name-value-source triplet
// Source is up to two words separated by a space. The first word is the source type: FLAG, ENV, DIR, CODE.
// The second is the name of the flag/file/environment var used.
type ConfigVar struct {
	Name   string
	Value  string
	Source string
}

// ConfigManager is an interface defining a key->value store that also knows where the values were set from.
type ConfigManager interface {
	Get(string) (string, error)
	GetIgnoreErr(string) string
	GetBool(string) (bool, error)
	GetV(string) (ConfigVar, error)
	GetVirtualMachine() lib.VirtualMachineName
	GetGroup() lib.GroupName
	GetAll() ([]ConfigVar, error)
	Set(string, string, string)
	SetPersistent(string, string, string) error
	Unset(string) error
	GetDebugLevel() int
	Force() bool
	Silent() bool
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
	mainFlags   *flag.FlagSet
	Dir         string
	Memo        map[string]ConfigVar
	Definitions map[string]string
}

type ConfigDirInvalidError struct {
	Path string
}

func (e *ConfigDirInvalidError) Error() string {
	return fmt.Sprintf("The config directory is '%s' but it doesn't seem to be a directory.", e.Path)
}

type CannotLoadDefinitionsError struct {
	Err error
}

func (e *CannotLoadDefinitionsError) Error() string {
	return fmt.Sprintf("Unable to load the definitions file from the Bytemark API.")
}

type ConfigReadError struct {
	Name string
	Path string
	Err  error
}

func (e *ConfigReadError) Error() string {
	return fmt.Sprintf("Unable to read config for %s from %s.", e.Name, e.Path)
}

type ConfigWriteError struct {
	Name string
	Path string
	Err  error
}

func (e *ConfigWriteError) Error() string {
	return fmt.Sprintf("Unable to write persistent config for %s (%s).", e.Name, e.Path)
}

// Do I really need to have the flags passed in here?
// Yes. Doing commands will be sorted out in a different place, and I don't want to touch it here.

// NewConfig sets up a new config struct. Pass in an empty string to default to ~/.bytemark
func NewConfig(configDir string, flags *flag.FlagSet) (config *Config, err error) {
	config = new(Config)
	config.Memo = make(map[string]ConfigVar)
	home := os.Getenv("HOME")
	if runtime.GOOS == "windows" {
		home = os.Getenv("APPDATA")

	}
	config.Dir = filepath.Join(home, "/.bytemark")
	config.mainFlags = flags
	if os.Getenv("BM_CONFIG_DIR") != "" {
		config.Dir = os.Getenv("BM_CONFIG_DIR")
	}

	if configDir != "" {
		config.Dir = configDir
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

	log.LogFile, err = os.Create(config.GetPath("debug.log"))
	if err != nil {
		log.Errorf("Couldn't open %s for writing\r\n", config.GetPath("debug.log"))
	}

	config.ImportFlags(flags)
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

			if flags != config.mainFlags {

				args := flags.Args()
				for _, arg := range args {
					if strings.HasPrefix(arg, "-") {
						log.Errorf("Flag-like argument '%s' specified after your arguments\r\nBe aware that only flags placed before your arguments are parsed.\r\nSee the help for the command you're calling for invocation examples.\r\n\r\n", arg)
						break
					}
				}
			}

			log.Silent = config.Silent()
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

// GetDebugLevel returns the current debug-level as an integer. This is used throughout the bytemark.co.uk/client library to determine verbosity of output.
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
		return val, nil
	}
	return config.read(name)
}

func (config *Config) GetVirtualMachine() (vm lib.VirtualMachineName) {
	vm.Account = config.GetIgnoreErr("account")
	vm.Group = config.GetIgnoreErr("group")
	vm.VirtualMachine = ""
	//TODO(telyn): make it possible to set a default VM?
	return vm
}

func (config *Config) GetGroup() (group lib.GroupName) {
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
		def := config.GetDefault("user")
		def.Name = "account"
		return def
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
	case "silent":
		return ConfigVar{"silent", "false", "CODE"}
	case "force":
		return ConfigVar{"force", "false", "CODE"}
	}
	return ConfigVar{"", "", ""}
}

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
	delete(config.Memo, name)
	return os.Remove(config.GetPath(name))
}

func (config *Config) Force() bool {
	force, _ := config.GetBool("force")
	return force
}
func (config *Config) Silent() bool {
	silent, _ := config.GetBool("silent")
	return silent
}

func (config *Config) PanelURL() string {
	endpoint := config.EndpointName()
	if strings.EqualFold(endpoint, "uk0.bigv.io") {
		return "https://panel-beta.bytemark.co.uk"
	}
	if strings.EqualFold(endpoint, "int.bigv.io") {
		// worrying leaky code?
		return "https://panel-int.vlan863.bytemark.uk0.bigv.io"
	}
	panel := config.GetIgnoreErr("panel-address")
	if panel == "" {
		panel = "https://your.panel.address"
	}
	return panel
}

func (config *Config) EndpointName() string {
	endpoint := config.GetIgnoreErr("endpoint")
	endpoint = strings.TrimPrefix(endpoint, "https://")
	endpoint = strings.TrimPrefix(endpoint, "http://") // it never hurts to be prepared
	return endpoint
}
