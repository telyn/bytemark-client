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

// Reads the config (usually in ~/.go-bigv)
// which is a set of flat files
// token - an OAuth 2.0 bearer token to use when authenticating
// username - the default username to use - if not present, $USER
// endpoint - the default endpoint to use - if not present, https://uk0.bigv.io
// endpoint_auth - the default auth_endpoint to use - if not present, https://auth.bytemark.co.uk
type Config struct {
	Dir         string
	Memo        []string
	Definitions map[string]string
}

// Sets up a new config struct. Pass in an empty string to default to ~/.go-bigv

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
	// try to read the file
	// or default
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
