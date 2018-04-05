package config

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/util/log"
)

// Config determines the configuration of the Bytemark client.
// It's responsible for handling things like the credentials to use and what endpoints to talk to.
//
// Each configuration item is read from the following places, falling back to successive places:
//
// Per-command command-line flags, global command-line flags, environment variables, configuration directory, hard-coded defaults
//
//The location of the configuration directory is read from global command-line flags, or is otherwise ~/.bytemark
//
type config struct {
	debugLevel  int
	Dir         string
	Memo        map[string]Var
	Definitions map[string]string
}

// New sets up a new config struct. Pass in an empty string to default to ~/.bytemark
func New(configDir string) (manager Manager, err error) {
	conf := new(config)
	conf.Memo = make(map[string]Var)
	home := os.Getenv("HOME")
	if runtime.GOOS == "windows" {
		home = os.Getenv("APPDATA")

	}

	conf.Dir = filepath.Join(home, "/.bytemark")
	if os.Getenv("BM_CONFIG_DIR") != "" {
		conf.Dir = os.Getenv("BM_CONFIG_DIR")
	}

	if configDir != "" {
		conf.Dir = strings.Replace(configDir, "~", os.Getenv("HOME"), -1)
	}

	err = os.MkdirAll(conf.Dir, 0700)
	if err != nil {

		return nil, err
	}

	stat, err := os.Stat(conf.Dir)
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return nil, &DirInvalidError{conf.Dir}
	}

	dbgLog := conf.GetPath("debug.log")
	// if there's already a debug.log, rename it debug.log.1
	_, err = os.Stat(dbgLog)
	if err == nil {
		_, err = os.Stat(dbgLog + ".1")
		if err == nil {
			// if debug.log.1 exists, remove it
			err = os.Remove(dbgLog + ".1") // we don't truly care if we couldn't clean up
			if err != nil {
				return nil, errors.New("Couldn't remove debug.log.1: " + err.Error())
			}
		}
		err = os.Rename(dbgLog, dbgLog+".1")
		if err != nil {
			return nil, errors.New("Couldn't rename debug.log to debug.log.1: " + err.Error())
		}
	}

	log.LogFile, err = os.Create(dbgLog)
	if err != nil {
		log.Errorf("Couldn't open %s for writing\r\n", conf.GetPath("debug.log"))
	}

	strDL, err := conf.Get("debug-level")
	if err != nil {
		log.Error(err)
	} else {
		debugLevel, err := strconv.ParseInt(strDL, 10, 0)
		if err == nil {
			conf.debugLevel = int(debugLevel)
			log.DebugLevel = int(debugLevel)
		}
	}
	return conf, nil
}

// ImportFlags reads all the flags from the passed FlagSet that have the same name as a valid configVar, and sets the configVar to that.
func (config *config) ImportFlags(flags *flag.FlagSet) (args []string) {
	if flags != nil {
		if flags.Parsed() {
			// dump all the flags into the memo
			// should be reet...reet?
			flags.Visit(func(f *flag.Flag) {
				val := config.massageFlagValue(f.Name, f.Value.String())
				config.Memo[f.Name] = Var{
					f.Name,
					val,
					"FLAG " + f.Name,
				}
			})

			strDL := config.GetIgnoreErr("debug-level")
			debugLevel, err := strconv.ParseInt(strDL, 10, 0)
			if err == nil {
				config.debugLevel = int(debugLevel)
				log.DebugLevel = int(debugLevel)
			}

			args = flags.Args()
		}
	}
	config.endpointOverrides()
	return
}

func (config *config) massageFlagValue(name string, val string) string {
	switch name {
	case "account":
		defAccount := config.GetDefault(val)
		return lib.ParseAccountName(val, defAccount.Value)
	}
	return val
}

func (config *config) endpointOverrides() {
	url, err := url.Parse(config.GetIgnoreErr("endpoint"))
	if err != nil {
		// we don't actually _care_ if there is an err - it's not int.
		return
	}
	if url.Host == "int.bigv.io" {
		config.Set("billing-endpoint", "", "CODE nullify billing-endpoint when using int")
		config.Set("spp-endpoint", "", "CODE nullify spp-endpoint when using int")
		config.Set("account", "bytemark", "CODE use bytemark account as default on int")
	}
}

// GetDebugLevel returns the current debug-level as an integer. This is used throughout the github.com/BytemarkHosting/bytemark-client library to determine verbosity of output.
func (config *config) GetDebugLevel() int {
	return config.debugLevel
}

// GetPath joins the given string onto the end of the Config.Dir path
func (config *config) GetPath(name string) string {
	return filepath.Join(config.Dir, name)
}

// Get returns the value of a Var. Used to simplify code when the source is unnecessary.
func (config *config) Get(name string) (string, error) {
	v, err := config.GetV(name)
	return v.Value, err
}

// GetSessionValidity returns the configured session validity or the default, if the configured one is not a valid int between 0 and infinity
func (config *config) GetSessionValidity() (validity int, err error) {
	validity = DefaultSessionValidity
	v, err := config.Get("session-validity")
	if err != nil {
		return
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return
	}
	// if the configured session validity is a negative number, return default without error
	// the brain will happily clamp the validity to whatever the maximum is if it's more than that, so we don't need to worry on that score
	// TODO(telyn): print a warning to cmd/bytemark.global.App.Writer
	if n < 0 {
		return
	}
	validity = n
	return

}

// GetIgnoreErr returns the value of a Var or an empty string , if it was unable to read it for whatever reason.
func (config *config) GetIgnoreErr(name string) string {
	s, _ := config.Get(name)
	return s
}

// GetV returns the Var for the given key.
func (config *config) GetV(name string) (Var, error) {
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
func (config *config) GetVirtualMachine() (vm lib.VirtualMachineName) {
	vm.Account = config.GetIgnoreErr("account")
	vm.Group = config.GetIgnoreErr("group")
	vm.VirtualMachine = ""
	return vm
}

// GetGroup returns a GroupName with the config's default group and account
func (config *config) GetGroup() (group lib.GroupName) {
	group.Account = config.GetIgnoreErr("account")
	group.Group = config.GetIgnoreErr("group")
	return group
}

// GetAll returns all of the available Vars in the Config.
func (config *config) GetAll() (vars []Var, err error) {
	vars = make([]Var, len(configVars))
	for i, v := range configVars {
		vars[i], err = config.GetV(v)
		if err != nil {
			return
		}
	}
	return
}

// GetDefault returns the default Var for the given key.
func (config *config) GetDefault(name string) Var {
	// ideally most of these should just be	os.Getenv("BM_"+name.Upcase().Replace("-","_"))
	switch name {
	case "user":
		if os.Getenv("BM_USER") == "" {
			return Var{"user", os.Getenv("USER"), "ENV USER"}
		}
		return Var{"user", os.Getenv("BM_USER"), "ENV BM_USER"}
	case "endpoint":
		v := Var{"endpoint", lib.DefaultURLs().Brain, "CODE"}

		val := os.Getenv("BM_ENDPOINT")
		if val != "" {
			v.Value = val
			v.Source = "ENV BM_ENDPOINT"
		}
		return v
	case "billing-endpoint":
		// this is here instead of in endpointOverrides bc u can't override endpointOverrides - and maybe staging will get a bmbilling. maybe it already has one? bmbilling.dev.bytemark.co.uk?
		if config.GetIgnoreErr("endpoint") == "https://staging.bigv.io" {
			return Var{"billing-endpoint", "", "CODE nullify billing-endpoint on bigv-staging"}
		}
		v := Var{"billing-endpoint", lib.DefaultURLs().Billing, "CODE"}
		if val := os.Getenv("BM_BILLING_ENDPOINT"); val != "" {
			v.Value = val
			v.Source = "ENV BM_BILLING_ENDPOINT"
		}
		return v
	case "spp-endpoint":
		v := Var{"spp-endpoint", lib.DefaultURLs().SPP, "CODE"}
		if val := os.Getenv("BM_SPP_ENDPOINT"); val != "" {
			v.Value = val
			v.Source = "ENV BM_SPP_ENDPOINT"
		}
		return v
	case "auth-endpoint":
		v := Var{"auth-endpoint", lib.DefaultURLs().Auth, "CODE"}

		val := os.Getenv("BM_AUTH_ENDPOINT")
		if val != "" {
			v.Value = val
			v.Source = "ENV BM_AUTH_ENDPOINT"
		}
		return v
	case "api-endpoint":
		v := Var{"api-endpoint", lib.DefaultURLs().API, "CODE"}

		val := os.Getenv("BM_API_ENDPOINT")
		if val != "" {
			v.Value = val
			v.Source = "ENV BM_API_ENDPOINT"
		}
	case "account":
		val := os.Getenv("BM_ACCOUNT")
		if val != "" {
			return Var{
				"account",
				val,
				"ENV BM_ACCOUNT",
			}
		}
		return Var{
			"account",
			"",
			"CODE",
		}
	case "group":
		val := os.Getenv("BM_GROUP")
		if val != "" {
			return Var{
				"group",
				val,
				"ENV BM_GROUP",
			}
		}
		return Var{"group", "default", "CODE"}
	case "debug-level":
		v := Var{"debug-level", "0", "CODE"}
		if val := os.Getenv("BM_DEBUG_LEVEL"); val != "" {
			v.Value = val
		}
		return v
	case "force":
		return Var{"force", "false", "CODE"}
	case "output-format":
		return Var{"output-format", "human", "CODE"}
	case "session-validity":
		return Var{"session-validity", fmt.Sprintf("%d", DefaultSessionValidity), "CODE"}
	}
	return Var{name, "", "UNSET"}
}

// GetBool returns the given configvar as a bool - true if it is set, not blank, and not equal to "false". false otherwise.
func (config *config) GetBool(name string) (bool, error) {
	v, err := config.Get(name)
	if err != nil {
		return false, err
	}
	return !(v == "" || v == "false"), nil
}

func (config *config) read(name string) (Var, error) {
	path := config.GetPath(name)
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return config.GetDefault(name), nil
		}

		return config.GetDefault(name), &ReadError{Name: name, Path: path, Err: err}
	}

	return Var{name, strings.TrimSpace(string(contents)), "FILE " + path}, nil
}

// Set stores the given key-value pair in config's Memo. This storage does not persist once the program terminates.
func (config *config) Set(name, value, source string) {
	config.Memo[name] = Var{name, value, source}
}

// SetPersistent writes a file to the config directory for the given key-value pair.
func (config *config) SetPersistent(name, value, source string) error {
	found := false
	for _, v := range configVars {
		if v == name {
			found = true
		}
	}
	if !found {
		return InvalidVarError{name}
	}
	path := config.GetPath(name)
	config.Set(name, value, source)
	err := ioutil.WriteFile(path, []byte(value), 0600)
	if err != nil {
		return &WriteError{Name: name, Path: path, Err: err}
	}
	return nil
}

// Unset removes the named key from both config's Memo and the user's config directory.
func (config *config) Unset(name string) (err error) {
	found := false
	for _, v := range configVars {
		if v == name {
			found = true
		}
	}
	if !found {
		return InvalidVarError{name}
	}
	delete(config.Memo, name)
	err = os.Remove(config.GetPath(name))
	if err != nil {
		info, statErr := os.Stat(config.Dir)
		if statErr != nil {
			if !info.IsDir() {
				return &DirInvalidError{config.Dir} // config dir is not a dir.
			}
			return nil // file didn't exist, so was already unset => success
		}
		return statErr // config dir couldn't be read for whatever reason
	}
	return // success
}

// PanelURL returns config's best guess at the correct URL for the bytemark panel for the cluster with the endpoint we're using. Basically it flips between panel.bytemark and panel-int.
func (config *config) PanelURL() string {
	endpoint := config.EndpointName()
	if strings.EqualFold(endpoint, "https://uk0.bigv.io") {
		return "https://panel.bytemark.co.uk"
	}
	if strings.EqualFold(endpoint, "https://int.bigv.io") {
		// am i leaking a secret?
		return "https://panel-int.admin.bytemark.co.uk"
	}
	panel := config.GetIgnoreErr("panel-address")
	if panel == "" {
		panel = "https://your.panel.address.example.com"
	}
	return panel
}

// EndpointName trims the URL scheme off the beginning of the endpoint.
// TODO(telyn): Why?
func (config *config) EndpointName() string {
	endpoint := config.GetIgnoreErr("endpoint")
	endpoint = strings.TrimPrefix(endpoint, "https://")
	endpoint = strings.TrimPrefix(endpoint, "http://") // it never hurts to be prepared
	return endpoint
}

// ConfigDir returns the path of the directory used to read config.
func (config *config) ConfigDir() string {
	return config.Dir
}
