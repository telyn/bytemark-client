package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var configVars = [...]string{
	"endpoint",
	"auth-endpoint",
	"user",
	"account",
	"token",
	"debug-level",
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
	Get(string) string
	GetBool(string) bool
	GetV(string) ConfigVar
	GetAll() []ConfigVar
	Set(string, string, string)
	SetPersistent(string, string, string)
	Unset(string) error
	GetDebugLevel() int

	ImportFlags(*flag.FlagSet) []string
}

// Params currently used:
// token - an OAuth 2.0 bearer token to use when authenticating
// username - the default username to use - if not present, $USER
// endpoint - the default endpoint to use - if not present, https://uk0.bigv.io
// auth-endpoint - the default auth API endpoint to use - if not present, https://auth.bytemark.co.uk

// A Config determines the configuration of the bigv client.
// It's responsible for handling things like the credentials to use and what endpoints to talk to.
//
// Each configuration item is read from the following places, falling back to successive places:
//
// Per-command command-line flags, global command-line flags, environment variables, configuration directory, hard-coded defaults
//
//The location of the configuration directory is read from global command-line flags, or is otherwise ~/.go-bigv
//
type Config struct {
	debugLevel  int
	Dir         string
	Memo        map[string]ConfigVar
	Definitions map[string]string
}

// Do I really need to have the flags passed in here?
// Yes. Doing commands will be sorted out in a different place, and I don't want to touch it here.

// NewConfig sets up a new config struct. Pass in an empty string to default to ~/.go-bigv
func NewConfig(configDir string, flags *flag.FlagSet) (config *Config) {
	config = new(Config)
	config.Memo = make(map[string]ConfigVar)
	config.Dir = filepath.Join(os.Getenv("HOME"), "/.go-bigv")
	if os.Getenv("BIGV_CONFIG_DIR") != "" {
		config.Dir = os.Getenv("BIGV_CONFIG_DIR")
	}

	if configDir != "" {
		config.Dir = configDir
	}

	err := os.MkdirAll(config.Dir, 0700)
	if err != nil {

		exit(err)
	}

	stat, err := os.Stat(config.Dir)
	if err != nil {
		exit(err)
	}

	if !stat.IsDir() {
		exit(nil, fmt.Sprintf("%s is not a directory", config.Dir))
	}

	debugLevel, err := strconv.ParseInt(config.Get("debug-level"), 10, 0)
	if err == nil {
		config.debugLevel = int(debugLevel)
	}

	config.ImportFlags(flags)
	return config
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
			return flags.Args()
		}
	}
	return nil
}

// GetDebugLevel returns the current debug-level as an integer. This is used throughout the bigv.io/client library to determine verbosity of output.
func (config *Config) GetDebugLevel() int {
	return config.debugLevel
}

// GetPath joins the given string onto the end of the Config.Dir path
func (config *Config) GetPath(name string) string {
	return filepath.Join(config.Dir, name)
}

// LoadDefinitions reads the local copy of the definitions json file, or downloads it from the endpoint if it's too old or nonexistant.
// Eventually this will be used to provide information on various things throughout the application
func (config *Config) LoadDefinitions() {
	stat, err := os.Stat(config.GetPath("definitions"))

	if err != nil || time.Since(stat.ModTime()) > 24*time.Hour {
		// TODO(telyn): grab it off the internet
		c := &http.Client{}
		req, err := http.NewRequest("GET", config.Get("endpoint")+"/definitions.json", nil)
		if err != nil {
			exit(err)
		}
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		res, err := c.Do(req)
		if err != nil {
			exit(err)
		}
		if res.StatusCode == 200 {
			responseBody, err := ioutil.ReadAll(res.Body)
			if err != nil {
				exit(err)
			}
			ioutil.WriteFile(config.GetPath("definitions"), responseBody, 0660)
		}
	} else {
		_, err := ioutil.ReadFile(config.GetPath("definitions"))
		if err != nil {
			exit(err, "Couldn't load definitions")
		}
	}

}

// Get returns the value of a ConfigVar. Used to simplify code when the source is unnecessary.
func (config *Config) Get(name string) string {
	return config.GetV(name).Value
}

// GetV returns the ConfigVar for the given key.
func (config *Config) GetV(name string) ConfigVar {
	// try to read the Memo
	name = strings.ToLower(name)
	if val, ok := config.Memo[name]; ok {
		return val
	}
	return config.read(name)
}

// GetAll returns all of the available ConfigVars in the Config.
func (config *Config) GetAll() []ConfigVar {
	vars := make([]ConfigVar, len(configVars))
	for i, v := range configVars {
		vars[i] = config.GetV(v)
	}
	return vars
}

// GetDefault returns the default ConfigVar for the given key.
func (config *Config) GetDefault(name string) ConfigVar {
	// ideally most of these should just be	os.Getenv("BIGV_"+name.Upcase().Replace("-","_"))
	switch name {
	case "user":
		// we don't actually want to default to USER - that will happen during Dispatcher's PromptForCredentials so it can be all "Hey you should bigv config set user <youruser>"
		return ConfigVar{"user", os.Getenv("BIGV_USER"), "ENV BIGV_USER"}
	case "endpoint":
		v := ConfigVar{"endpoint", "https://uk0.bigv.io", "CODE"}

		val := os.Getenv("BIGV_ENDPOINT")
		if val != "" {
			v.Value = val
			v.Source = "ENV BIGV_ENDPOINT"
		}
		return v
	case "auth-endpoint":
		v := ConfigVar{"auth-endpoint", "https://auth.bytemark.co.uk", "CODE"}

		val := os.Getenv("BIGV_AUTH_ENDPOINT")
		if val != "" {
			v.Value = val
			v.Source = "ENV BIGV_AUTH_ENDPOINT"
		}
		return v
	case "account":
		val := os.Getenv("BIGV_ACCOUNT")
		if val != "" {
			return ConfigVar{
				"account",
				val,
				"ENV BIGV_AUTH_ENDPOINT",
			}
		}
		def := config.GetDefault("user")
		def.Name = "account"
		return def
	case "debug-level":
		v := ConfigVar{"debug-level", "0", "CODE"}
		if val := os.Getenv("BIGV_DEBUG_LEVEL"); val != "" {
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

func (config *Config) GetBool(name string) bool {
	return !(config.Get(name) == "" || config.Get(name) == "false")
}

func (config *Config) read(name string) ConfigVar {
	contents, err := ioutil.ReadFile(config.GetPath(name))
	if err != nil {
		if os.IsNotExist(err) {
			return config.GetDefault(name)
		}

		exit(err, fmt.Sprintf("Couldn't read config for %s", name))
	}

	return ConfigVar{name, string(contents), "FILE " + config.GetPath(name)}
}

// Set stores the given key-value pair in config's Memo. This storage does not persist once the program terminates.
func (config *Config) Set(name, value, source string) {
	config.Memo[name] = ConfigVar{name, value, source}
}

// SetPersistent writes a file to the config directory for the given key-value pair.
func (config *Config) SetPersistent(name, value, source string) {
	config.Set(name, value, source)
	err := ioutil.WriteFile(config.GetPath(name), []byte(value), 0600)
	if err != nil {
		exit(err, fmt.Sprintf("Couldn't write to config directory "+config.Dir))
	}
}

// Unset removes the named key from both config's Memo and the user's config directory.
func (config *Config) Unset(name string) error {
	delete(config.Memo, name)
	return os.Remove(config.GetPath(name))
}
