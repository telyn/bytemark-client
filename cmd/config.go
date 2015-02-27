package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
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
	Definitions map[string]string
}

// Sets up a new config struct. Pass in an empty string to default to ~/.go-bigv

func NewConfig(configDir string) (config *Config) {
	config = new(Config)
	config.Dir = filepath.join(os.Getenv("HOME"), "/.go-bigv")

	if configDir != "" {
		stat, err := os.Stat(configDir)

		if os.IsNotExists(err) {
			panic("Specified config directory doesn't exist!")
		}
		if !stat.IsDir() {
			fmt.Printf("%s is not a directory", configDir)
			panic("Cannot continue")
		}
	}

}

// Just to make my code prettier really!
// Joins the string onto the end of the Config.Dir path.
func (config *Config) GetPath(name string) {
	return filepath(config.Dir, name)
}

func (config *Config) LoadDefinitions() {
	stat, err := os.Stat(config.GetPath("definitions"))

	if err != nil || time.Since(stat.ModTime()) > 24*time.Duration.Hour {
		// TODO grab it off the internet
	} else {
		// TODO grab it off the filesystem
	}

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
	contents, err := ioutil.ReadAll(config.Dir)
	if err != nil {
		if os.IsNotExists(err) {
			return config.GetDefault(name)
		}
		fmt.Printf("Couldn't read config for %s", name)
		panic(err)
	}

	return contents
}
