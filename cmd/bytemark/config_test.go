package main

import (
	"flag"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/cheekybits/is"
	"github.com/urfave/cli"
	"strings"
	"testing"
)

func TestConfigAccountValidation(t *testing.T) {
	config, client := baseTestSetup(t, false)

	config.When("GetGroup").Return(lib.GroupName{Group: "default-group", Account: "default-account"})
	config.When("GetIgnoreErr", "account").Return("")
	config.When("GetIgnoreErr", "token").Return("test-token")

	ctx := Context{}

	runAccountTests(t, &ctx, client, getValidationTests()["account"], validateAccountForConfig)
}

func TestConfigGroupValidation(t *testing.T) {
	config, client := baseTestSetup(t, false)

	config.When("GetGroup").Return(lib.GroupName{Group: "", Account: ""})
	config.When("GetIgnoreErr", "account").Return("")
	config.When("GetIgnoreErr", "token").Return("test-token")

	flagset := flag.NewFlagSet("TestValidation", flag.ContinueOnError)
	flagset.Bool("force", false, "")

	ctx := Context{
		Context: cliContextWrapper{cli.NewContext(global.App, flagset, nil)},
	}
	t.Logf("Testing validateGroupForConfig\r\n")
	runGroupTests(t, &ctx, client, getValidationTests()["group"], validateGroupForConfig)
}

func TestConfigEndpointValidation(t *testing.T) {
	endpoints := getValidationTests()["endpoint"]
	runEndpointTests(t, "endpoint (direct)", endpoints, validateEndpointForConfig)
}

func TestConfigValidations(t *testing.T) {
	config, client := baseTestSetup(t, false)

	config.When("GetGroup").Return(lib.GroupName{Group: "", Account: ""})
	config.When("GetIgnoreErr", "account").Return("")
	config.When("GetIgnoreErr", "token").Return("test-token")

	flagset := flag.NewFlagSet("TestValidationWithForce", flag.ContinueOnError)
	flagset.Bool("force", false, "")

	ctx := Context{
		Context: cliContextWrapper{cli.NewContext(global.App, flagset, nil)},
	}

	tests := getValidationTests()

	for _, varname := range []string{"endpoint", "api-endpoint", "auth-endpoint", "billing-endpoint", "spp-endpoint"} {
		runEndpointTests(t, varname, tests["endpoint"], func(endpoint string) error {
			return validateConfigValue(&ctx, varname, endpoint)
		})
	}
	runGroupTests(t, &ctx, client, tests["group"], func(c *Context, groupName string) error {
		return validateConfigValue(c, "group", groupName)
	})
	runAccountTests(t, &ctx, client, tests["account"], func(c *Context, accountName string) error {

		return validateConfigValue(c, "account", accountName)
	})
	runDebugLevelTests(t, tests["debug-level"], func(level string) error {
		return validateConfigValue(&ctx, "debug-level", level)
	})

	if ok, vErr := config.Verify(); !ok {
		t.Fatal(vErr)
	}
}

func TestCommandConfigSet(t *testing.T) {
	is := is.New(t)
	config, client := baseTestSetup(t, false)

	// setup sets up all the necessary config defaulty stuff for all our tests

	setup := func() {
		config.Reset()
		client.Reset()
		config.When("GetV", "user").Return(util.ConfigVar{"user", "old-test-user", "config"})
		config.When("GetV", "account").Return(util.ConfigVar{"account", "", ""})
		config.When("GetV", "endpoint").Return(util.ConfigVar{"endpoint", "", ""})
		config.When("GetV", "group").Return(util.ConfigVar{"group", "", ""})
		config.When("GetV", "debug-level").Return(util.ConfigVar{"debug-level", "", ""})
		config.When("Get", "token").Return("test-token", nil)
		config.When("GetIgnoreErr", "user").Return("old-test-user")
		config.When("GetIgnoreErr", "yubikey").Return("")
		config.When("GetIgnoreErr", "2fa-otp").Return("")
		config.When("GetIgnoreErr", "account").Return("")
		config.When("GetGroup").Return(lib.GroupName{})
		client.When("AuthWithToken", "test-token").Return(nil)
	}

	setup()

	config.When("SetPersistent", "user", "test-user", "CMD set").Times(1)

	err := global.App.Run(strings.Split("bytemark config set user test-user", " "))
	is.Nil(err)

	if ok, vErr := config.Verify(); !ok {
		t.Fatal(vErr)
	}

	err = global.App.Run(strings.Split("bytemark config set flimflam test-user", " "))
	is.NotNil(err)

	// test setting all the other variables

	tests := getValidationTests()
	for varname := range tests {
		for value, errSpec := range tests[varname] {
			setup()
			switch varname {
			case "account":
				setupAccountTest(client, value, errSpec.err)
			case "group":
				setupGroupTest(client, value, errSpec.err)
			}
			if !errSpec.shouldErr {
				config.When("GetIgnoreErr", varname).Return("")
				config.When("SetPersistent", varname, value, "CMD set").Times(1)
			}

			err = global.App.Run([]string{"bytemark", "config", "set", varname, value})
			if errSpec.shouldErr && err == nil {
				t.Errorf("bytemark config set %s %s should've errored but didn't.\r\n", varname, value)
			} else if !errSpec.shouldErr && err != nil {
				t.Errorf("bytemark config set %s %s should've suceeded but didn't: %s\r\n", varname, value, err.Error())
			}
			if ok, vErr := config.Verify(); !ok {
				t.Errorf("bytemark config set %s %s - config.Verify error: %s\r\n", varname, value, vErr)
			}
		}
	}

	// now test that --force works.

	for varname := range tests {
		for value := range tests[varname] {
			setup()
			config.When("GetIgnoreErr", varname).Return("")
			config.When("SetPersistent", varname, value, "CMD set").Times(1)

			err = global.App.Run([]string{"bytemark", "config", "set", "--force", varname, value})
			if err != nil {
				t.Errorf("bytemark config set %s %s should've suceeded but didn't: %s\r\n", varname, value, err.Error())
			}
			if ok, vErr := config.Verify(); !ok {
				t.Errorf("bytemark config set %s %s - config.Verify error: %s\r\n", varname, value, vErr)
			}
		}
	}

}

func setupAccountTest(c *mocks.Client, name string, err error) {

	if err != nil {
		c.When("GetAccount", name).Return(nil, err)
	} else {
		c.When("GetAccount", name).Return(&lib.Account{}, nil)
	}
}

func setupGroupTest(c *mocks.Client, name string, err error) {
	groupName := lib.GroupName{Group: name}
	if err != nil {
		c.When("GetGroup", groupName).Return(nil, err).Times(1)
	} else {
		c.When("GetGroup", groupName).Return(&brain.Group{}, nil).Times(1)
	}
}

type validationFn func(ctx *Context, name string) error

func runAccountTests(t *testing.T, ctx *Context, c *mocks.Client, accounts map[string]errSpec, fnUnderTest validationFn) {
	for a, spec := range accounts {
		setupAccountTest(c, a, spec.err)
		err := fnUnderTest(ctx, a)
		if spec.shouldErr && err == nil {
			t.Errorf("testing set account %s: should error, but didn't\r\n", a)
		} else if !spec.shouldErr && err != nil {
			t.Errorf("testing set account %s: should not error, but did with '%s'\r\n", a, err.Error())
		}
		if ok, vErr := c.Verify(); !ok {
			t.Fatal(vErr)
		}
		c.Reset()
		ctx.Reset()
	}
}

func runDebugLevelTests(t *testing.T, debugLvls map[string]errSpec, fnUnderTest func(string) error) {
	for level, spec := range debugLvls {
		err := fnUnderTest(level)
		if spec.shouldErr && err == nil {
			t.Errorf("testing set debug-level %s: should error, but didn't\r\n", level)
		} else if !spec.shouldErr && err != nil {
			t.Errorf("testing set debug-level %s: should not error, but did with '%s'\r\n", level, err.Error())
		}
	}

}

func runEndpointTests(t *testing.T, varname string, endpoints map[string]errSpec, fnUnderTest func(string) error) {
	for e, spec := range endpoints {
		err := fnUnderTest(e)
		if spec.shouldErr && err == nil {
			t.Errorf("testing set %s %s: should error, but didn't\r\n", varname, e)
		} else if !spec.shouldErr && err != nil {
			t.Errorf("testing set %s %s: should not error, but did with '%s'\r\n", varname, e, err.Error())
		}
	}
}

func runGroupTests(t *testing.T, c *Context, client *mocks.Client, groups map[string]errSpec, fnUnderTest validationFn) {
	for g, spec := range groups {
		t.Log(g)
		setupGroupTest(client, g, spec.err)
		err := fnUnderTest(c, g)
		if spec.shouldErr && err == nil {
			t.Errorf("testing set group %s: should error, but didn't\r\n", g)
		} else if !spec.shouldErr && err != nil {
			t.Errorf("testing set group %s: should not error, but did with '%s'\r\n", g, err.Error())
		}
		if ok, vErr := client.Verify(); !ok {
			t.Fatalf("testing set group %s: %v\r\n", g, vErr)
		}
		client.Reset()
		c.Reset()
	}
}

type errSpec struct {
	// true if there should be an error
	shouldErr bool
	// the error to output from the lib mock
	err error
}

func getValidationTests() map[string]map[string]errSpec {
	return map[string]map[string]errSpec{
		"endpoint": map[string]errSpec{
			"https://uk0.bigv.io":          {shouldErr: false},
			"https://uk0.bigv.io/":         {shouldErr: false},
			"test":                         {shouldErr: true},
			"/bivoac":                      {shouldErr: true},
			"http://insecure-endpoint":     {shouldErr: false},
			"gopher://really":              {shouldErr: true},
			"http://hamlet's uncle did it": {shouldErr: true},
			"http:///no-hostname":          {shouldErr: true},
		},
		"account": map[string]errSpec{
			"":                    {shouldErr: true, err: lib.NotFoundError{}},
			"existent-account":    {shouldErr: false},
			"nonexistent-account": {shouldErr: true, err: lib.NotFoundError{}},
		},

		"group": map[string]errSpec{
			"":                  {shouldErr: true, err: lib.NotFoundError{}},
			"existent-group":    {shouldErr: false},
			"nonexistent-group": {shouldErr: true, err: lib.NotFoundError{}},
		},

		"debug-level": map[string]errSpec{
			"":         {shouldErr: true},
			"1":        {shouldErr: false},
			"9000":     {shouldErr: false},
			"-5":       {shouldErr: true},
			"barfbags": {shouldErr: true},
		},
	}
}
