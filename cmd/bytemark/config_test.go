package main

import (
	"flag"
	"fmt"
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/cheekybits/is"
	"github.com/urfave/cli"
)

func TestConfigAccountValidation(t *testing.T) {
	config, client, cliapp := testutil.BaseTestSetup(t, false, commands)

	config.When("GetGroup").Return(lib.GroupName{Group: "default-group", Account: "default-account"})
	config.When("GetIgnoreErr", "account").Return("")
	config.When("GetIgnoreErr", "token").Return("test-token")

	ctx := app.Context{
		Context: app.CliContextWrapper{&cli.Context{
			App: cliapp,
		}},
	}

	runAccountTests(t, &ctx, client, getValidationTests()["account"], validateAccountForConfig)
}

func TestConfigGroupValidation(t *testing.T) {
	config, client, cliapp := testutil.BaseTestSetup(t, false, commands)

	config.When("GetGroup").Return(lib.GroupName{Group: "", Account: ""})
	config.When("GetIgnoreErr", "account").Return("")
	config.When("GetIgnoreErr", "token").Return("test-token")

	flagset := flag.NewFlagSet("TestValidation", flag.ContinueOnError)
	flagset.Bool("force", false, "")

	ctx := app.Context{
		Context: app.CliContextWrapper{cli.NewContext(cliapp, flagset, nil)},
	}
	t.Logf("Testing validateGroupForConfig\r\n")
	runGroupTests(t, &ctx, client, getValidationTests()["group"], validateGroupForConfig)
}

func TestConfigEndpointValidation(t *testing.T) {
	endpoints := getValidationTests()["endpoint"]
	runEndpointTests(t, "endpoint (direct)", endpoints, validateEndpointForConfig)
}

func TestConfigValidations(t *testing.T) {
	config, client, cliapp := testutil.BaseTestSetup(t, false, commands)

	config.When("GetGroup").Return(lib.GroupName{Group: "", Account: ""})
	config.When("GetIgnoreErr", "account").Return("")
	config.When("GetIgnoreErr", "token").Return("test-token")

	flagset := flag.NewFlagSet("TestValidationWithForce", flag.ContinueOnError)
	flagset.Bool("force", false, "")

	ctx := app.Context{
		Context: app.CliContextWrapper{cli.NewContext(cliapp, flagset, nil)},
	}

	tests := getValidationTests()

	for _, varname := range []string{"endpoint", "api-endpoint", "auth-endpoint", "billing-endpoint", "spp-endpoint"} {
		runEndpointTests(t, varname, tests["endpoint"], func(endpoint string) error {
			return validateConfigValue(&ctx, varname, endpoint)
		})
	}
	runGroupTests(t, &ctx, client, tests["group"], func(c *app.Context, groupName string) error {
		return validateConfigValue(c, "group", groupName)
	})
	runAccountTests(t, &ctx, client, tests["account"], func(c *app.Context, accountName string) error {

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

	// setup sets up all the necessary config defaulty stuff for all our tests

	//	setup := func() {
	//		config.Reset()
	//		client.Reset()
	//		config.When("GetV", "user").Return(util.ConfigVar{"user", "old-test-user", "config"})
	//		config.When("GetV", "account").Return(util.ConfigVar{"account", "", ""})
	//		config.When("GetV", "endpoint").Return(util.ConfigVar{"endpoint", "", ""})
	//		config.When("GetV", "group").Return(util.ConfigVar{"group", "", ""})
	//		config.When("GetV", "debug-level").Return(util.ConfigVar{"debug-level", "", ""})
	//		config.When("GetIgnoreErr", "token").Return("test-token")
	//		config.When("GetIgnoreErr", "user").Return("old-test-user")
	//		config.When("GetIgnoreErr", "yubikey").Return("")
	//		config.When("GetIgnoreErr", "2fa-otp").Return("")
	//		config.When("GetIgnoreErr", "account").Return("")
	//		config.When("GetGroup").Return(lib.GroupName{})
	//		client.When("AuthWithToken", "test-token").Return(nil)
	//	}

	t.Run("set user", func(t *testing.T) {
		config, _, app := testutil.BaseTestSetup(t, false, commands)
		config.When("GetV", "user").Return(util.ConfigVar{"user", "old-test-user", "config"})
		config.When("GetIgnoreErr", "user").Return("old-test-user")

		config.When("SetPersistent", "user", "test-user", "CMD set").Times(1)

		err := app.Run(strings.Split("bytemark config set user test-user", " "))
		is.Nil(err)

		if ok, vErr := config.Verify(); !ok {
			t.Fatal(vErr)
		}
	})

	t.Run("set flimflam", func(t *testing.T) {
		config, _, app := testutil.BaseTestSetup(t, false, commands)

		err := app.Run(strings.Split("bytemark config set flimflam test-user", " "))
		is.NotNil(err)
		if ok, vErr := config.Verify(); !ok {
			t.Fatal(vErr)
		}
	})

	// test setting all the other variables

	tests := getValidationTests()
	for varname := range tests {
		for value, errSpec := range tests[varname] {
			t.Run(fmt.Sprintf("set %s %s", varname, value), func(t *testing.T) {
				config, client, app := testutil.BaseTestSetup(t, false, commands)

				switch varname {
				case "account":
					config, client, app = testutil.BaseTestAuthSetup(t, false, commands)
					client.When("GetAccount", value).Return(&lib.Account{}, errSpec.err)
				case "group":
					config, client, app = testutil.BaseTestAuthSetup(t, false, commands)
					groupName := lib.GroupName{Group: value}
					config.When("GetGroup").Return(lib.GroupName{Group: "not-real"})
					client.When("GetGroup", groupName).Return(&brain.Group{}, errSpec.err).Times(1)
				}

				config.When("GetV", varname).Return(util.ConfigVar{varname, "old-test-" + varname, "config"})
				if !errSpec.shouldErr {
					config.When("GetIgnoreErr", varname).Return("")
					config.When("SetPersistent", varname, value, "CMD set").Times(1)
				}

				err := app.Run([]string{"bytemark", "config", "set", varname, value})
				if errSpec.shouldErr && err == nil {
					t.Errorf("bytemark config set %s %s should've errored but didn't.\r\n", varname, value)
				} else if !errSpec.shouldErr && err != nil {
					t.Errorf("bytemark config set %s %s should've succeeded but didn't: %s\r\n", varname, value, err.Error())
				}
				if ok, vErr := config.Verify(); !ok {
					t.Errorf("bytemark config set %s %s - config.Verify error: %s\r\n", varname, value, vErr)
				}
			})
		}
	}

	// now test that --force works.
	for varname := range tests {
		for value, errSpec := range tests[varname] {
			t.Run(fmt.Sprintf("force set %s %s", varname, value), func(t *testing.T) {
				config, client, app := testutil.BaseTestSetup(t, false, commands)

				switch varname {
				case "account":
					config, client, app = testutil.BaseTestAuthSetup(t, false, commands)
					client.When("GetAccount", value).Return(&lib.Account{}, errSpec.err)
				case "group":
					config, client, app = testutil.BaseTestAuthSetup(t, false, commands)
					groupName := lib.GroupName{Group: value}
					config.When("GetGroup").Return(lib.GroupName{Group: "not-real"})
					client.When("GetGroup", groupName).Return(&brain.Group{}, errSpec.err).Times(1)
				}
				config.When("GetV", varname).Return(util.ConfigVar{varname, "old-test-" + varname, "config"})
				config.When("GetIgnoreErr", varname).Return("test-old-" + varname)
				config.When("SetPersistent", varname, value, "CMD set").Times(1)

				err := app.Run([]string{"bytemark", "config", "set", "--force", varname, value})
				if err != nil {
					t.Errorf("bytemark config set %s %s should've succeeded but didn't: %s\r\n", varname, value, err.Error())
				}
				if ok, vErr := config.Verify(); !ok {
					t.Errorf("bytemark config set %s %s - config.Verify error: %s\r\n", varname, value, vErr)
				}
			})
		}
	}

}

func setupGroupTest(c *mocks.Client, name string, err error) {
	groupName := lib.GroupName{Group: name}
	c.When("GetGroup", groupName).Return(&brain.Group{}, err).Times(1)
}

type validationFn func(ctx *app.Context, name string) error

func runAccountTests(t *testing.T, ctx *app.Context, c *mocks.Client, accounts map[string]errSpec, fnUnderTest validationFn) {
	for a, spec := range accounts {
		c.When("GetAccount", a).Return(&lib.Account{}, spec.err)
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

func runGroupTests(t *testing.T, c *app.Context, client *mocks.Client, groups map[string]errSpec, fnUnderTest validationFn) {
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
