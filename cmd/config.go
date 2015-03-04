package cmd

import (
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
func NewConfig(configDir string) (config *Config) {
	config = new(Config)
	config.Dir = filepath.Join(os.Getenv("HOME"), "/.go-bigv")
	if os.Getenv("BIGV_CONFIG_DIR") != "" {
		config.Dir = os.Getenv("BIGV_CONFIG_DIR")
	}

	if configDir != "" {
		// TODO should probably just try to make config.Dir and if it already exists, bonus.
		// if it can't be made panic()
		stat, err := os.Stat(configDir)

		if os.IsNotExist(err) {
			panic("Specified config directory doesn't exist!")
		}
		if !stat.IsDir() {
			fmt.Printf("%s is not a directory", configDir)
			panic("Cannot continue")
		}
		config.Dir = configDir
	}
	return config
}

// Just to make my code prettier really!
// Joins the string onto the end of the Config.Dir path.
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
		// TODO grab it off the internet
		//		url := config.GetUrl("definitions.json")
	} else {
		_, err := ioutil.ReadFile(config.GetPath("definitions"))
		if err != nil {
			panic("Couldn't load definitions")
		}
		// TODO grab it off the filesystem
	}

}

func (config *Config) Get(name string) string {
	// try to read the Memo
	if val, ok := config.Memo[name]; ok {
		return val
	} else {
		// try to read the file
		val, err := ioutil.ReadFile(config.GetPath(name))
		if err != nil {
			return config.GetDefault(name)
		}
		return string(val)
		// or default
	}
	return ""
}

func (config *Config) GetDefault(name string) string {
	switch name {
	case "user":
		return FirstNotEmpty(os.Getenv("BIGV_USER"), os.Getenv("USER"))
	case "endpoint":
		return "https://uk0.bigv.io"
	case "endpoint_auth":
		return "https://auth.bytemark.co.uk"
	case "account":
		return FirstNotEmpty(os.Getenv("BIGV_ACCOUNT"), os.Getenv("BIGV_USER"), os.Getenv("USER"))
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
