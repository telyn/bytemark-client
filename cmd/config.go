package cmd

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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
	Dir         string
	Memo        map[string]string
	Definitions map[string]string
}

// Do I really need to have the flags passed in here?
// Yes. Doing commands will be sorted out in a different place, and I don't want to touch it here.

// NewConfig sets up a new config struct. Pass in an empty string to default to ~/.go-bigv
func NewConfig(configDir string, flags *flag.FlagSet) (config *Config) {
	config = new(Config)
	config.Memo = make(map[string]string)
	config.Dir = filepath.Join(os.Getenv("HOME"), "/.go-bigv")
	if os.Getenv("BIGV_CONFIG_DIR") != "" {
		config.Dir = os.Getenv("BIGV_CONFIG_DIR")
	}

	if configDir != "" {
		err := os.MkdirAll(configDir, 0600)
		if err != nil {
			// TODO(telyn): Better error handling here

			panic(err)
		}

		stat, err := os.Stat(configDir)
		if err != nil {
			// TODO(telyn): Better error handling here
			panic(err)
		}

		if !stat.IsDir() {
			fmt.Printf("%s is not a directory", configDir)
			panic("Cannot continue")
		}
		config.Dir = configDir
	}

	if flags != nil {
		// dump all the flags into the memo
		// should be reet...reet?
		flags.Visit(func(f *flag.Flag) {
			config.Memo[f.Name] = f.Value.String()
		})
	}
	return config
}

// GetPath joins the given string onto the end of the Config.Dir path
func (config *Config) GetPath(name string) string {
	return filepath.Join(config.Dir, name)
}
func (config *Config) GetUrl(path ...string) *url.URL {
	url, err := url.Parse(config.Get("endpoint"))
	if err != nil {
		panic("Endpoint is not a valid URL")
	}
	url.Parse("/" + strings.Join([]string(path), "/"))
	return url
}

func (config *Config) LoadDefinitions() {
	stat, err := os.Stat(config.GetPath("definitions"))

	if err != nil || time.Since(stat.ModTime()) > 24*time.Hour {
		// TODO(telyn): grab it off the internet
		//		url := config.GetUrl("definitions.json")
	} else {
		_, err := ioutil.ReadFile(config.GetPath("definitions"))
		if err != nil {
			panic("Couldn't load definitions")
		}
	}

}

func (config *Config) Get(name string) string {
	// try to read the Memo
	if val, ok := config.Memo[name]; ok {
		return val
	} else {
		return config.Read(name)
	}
	return ""
}

func (config *Config) GetDefault(name string) string {
	// ideally most of these should just be	os.Getenv("BIGV_"+name.Upcase().Replace("-","_"))
	switch name {
	case "user":
		return FirstNotEmpty(os.Getenv("BIGV_USER"), os.Getenv("USER"))
	case "endpoint":
		return FirstNotEmpty(os.Getenv("BIGV_ENDPOINT"), "https://uk0.bigv.io")
	case "auth-endpoint":
		return FirstNotEmpty(os.Getenv("BIGV_AUTH_ENDPOINT"), "https://auth.bytemark.co.uk")
	case "account":
		return FirstNotEmpty(os.Getenv("BIGV_ACCOUNT"), os.Getenv("BIGV_USER"), os.Getenv("USER"))
	case "debug-level":
		return FirstNotEmpty(os.Getenv("BIGV_DEBUG_LEVEL"), "0")
	}
	return ""
}

func (config *Config) Read(name string) string {
	contents, err := ioutil.ReadFile(config.GetPath(name))
	if err != nil {
		if os.IsNotExist(err) {
			return config.GetDefault(name)
		}
		fmt.Printf("Couldn't read config for %s", name)
		panic(err)
	}

	return string(contents)
}

// Set stores the given key-value pair in config's Memo. This storage does not persist once the program terminates.
func (config *Config) Set(name string, value string) {
	config.Memo[name] = value
}

// SetPersistent writes a file to the config directory for the given key-value pair.
func (config *Config) SetPersistent(name, value string) {
	config.Set(name, value)
	err := ioutil.WriteFile(config.GetPath("user"), []byte(value), 0600)
	if err != nil {
		fmt.Println("Couldn't write to config directory " + config.Dir)
		panic(err)
	}
}
