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

func TestCommandConfigSet(t *testing.T) {
	is := is.New(t)
	config, _ := baseTestSetup(t, false)

	config.When("GetV", "user").Return(util.ConfigVar{"user", "old-test-user", "config"})
	config.When("GetIgnoreErr", "user").Return("old-test-user")

	config.When("SetPersistent", "user", "test-user", "CMD set").Times(1)

	err := global.App.Run(strings.Split("bytemark config set user test-user", " "))
	is.Nil(err)

	if ok, vErr := config.Verify(); !ok {
		t.Fatal(vErr)
	}

	err = global.App.Run(strings.Split("bytemark config set flimflam test-user", " "))
	is.NotNil(err)
}

func endpointTests(t *testing.T, varname string, fnUnderTest func(string) error) {
	endpoints := map[string]bool{
		"https://uk0.bigv.io":      false,
		"https://uk0.bigv.io/":     false,
		"test":                     true,
		"/bivoac":                  true,
		"http://insecure-endpoint": false,
		"gopher://really":          true,
	}

	for e, shouldError := range endpoints {
		err := fnUnderTest(e)
		if shouldError && err == nil {
			t.Errorf("testing set %s %s: should error, but didn't\r\n", varname, e)
		} else if !shouldError && err != nil {
			t.Errorf("testing set %s %s: should not error, but did with '%s'\r\n", varname, e, err.Error())
		}
	}
}

func accountTests(t *testing.T, ctx *Context, c *mocks.Client, fnUnderTest validationFn) {
	accounts := map[string]bool{
		"":                    false,
		"existent-account":    true,
		"nonexistent-account": false,
	}

	for a, shouldError := range accounts {
		setupAccountTest(c, a, shouldError)
		err := fnUnderTest(ctx, a)
		if shouldError && err == nil {
			t.Errorf("testing set account %s: should error, but didn't\r\n", a)
		} else if !shouldError && err != nil {
			t.Errorf("testing set account %s: should not error, but did with '%s'\r\n", a, err.Error())
		}
		if ok, vErr := c.Verify(); !ok {
			t.Fatal(vErr)
		}
		c.Reset()
		ctx.Reset()

	}
}

func setupAccountTest(c *mocks.Client, name string, shouldError bool) {
	c.When("ParseAccountName", name, []string{""}).Return(name)

	if shouldError {
		c.When("GetAccount", name).Return(nil, &lib.NotFoundError{})
	} else {
		c.When("GetAccount", name).Return(&lib.Account{}, nil)
	}
}

func setupGroupTest(c *mocks.Client, name string, shouldError bool) {
	groupName := lib.GroupName{Group: name}
	c.When("ParseGroupName", name, []*lib.GroupName{{}}).Return(&groupName)
	if shouldError {
		c.When("GetGroup", &groupName).Return(nil, &lib.NotFoundError{}).Times(1)
	} else {
		c.When("GetGroup", &groupName).Return(&brain.Group{}, nil).Times(1)
	}
}

type validationFn func(ctx *Context, name string) error

func groupTests(t *testing.T, c *Context, client *mocks.Client, fnUnderTest validationFn) {
	groups := map[string]bool{
		"":                  true,
		"existent-group":    false,
		"nonexistent-group": true,
	}

	for g, shouldError := range groups {
		t.Log(g)
		setupGroupTest(client, g, shouldError)
		err := fnUnderTest(c, g)
		if shouldError && err == nil {
			t.Errorf("testing set group %s: should error, but didn't\r\n", g)
		} else if !shouldError && err != nil {
			t.Errorf("testing set group %s: should not error, but did with '%s'\r\n", g, err.Error())
		}
		if ok, vErr := client.Verify(); !ok {
			t.Fatalf("testing set group %s: %v\r\n", g, vErr)
		}
		client.Reset()
		c.Reset()
	}
}

func TestEndpointValidation(t *testing.T) {
	endpointTests(t, "endpoint (direct)", validateEndpointForConfig)
}

func TestValidation(t *testing.T) {
	config, client := baseTestSetup(t, false)

	config.When("GetGroup").Return(&lib.GroupName{Group: "", Account: ""})
	config.When("GetIgnoreErr", "account").Return("")

	flagset := flag.NewFlagSet("TestValidation", flag.ContinueOnError)
	flagset.Bool("force", false, "")

	ctx := Context{
		Context: cli.NewContext(global.App, flagset, nil),
	}

	t.Logf("Testing validateConfigValue\r\n")
	for _, varname := range []string{"endpoint", "api-endpoint", "auth-endpoint", "billing-endpoint", "spp-endpoint"} {
		endpointTests(t, varname, func(endpoint string) error {
			return validateConfigValue(&ctx, varname, endpoint)
		})
	}
	groupTests(t, &ctx, client, func(c *Context, groupName string) error {
		return validateConfigValue(c, "group", groupName)
	})
	accountTests(t, &ctx, client, func(c *Context, accountName string) error {
		return validateConfigValue(c, "account", accountName)
	})

	if ok, vErr := config.Verify(); !ok {
		t.Fatal(vErr)
	}
}

func TestGroupValidation(t *testing.T) {
	config, client := baseTestSetup(t, false)

	config.When("GetGroup").Return(&lib.GroupName{Group: "", Account: ""})
	config.When("GetIgnoreErr", "account").Return("")

	flagset := flag.NewFlagSet("TestValidation", flag.ContinueOnError)
	flagset.Bool("force", false, "")

	ctx := Context{
		Context: cli.NewContext(global.App, flagset, nil),
	}
	t.Logf("Testing validateGroupForConfig\r\n")
	groupTests(t, &ctx, client, validateGroupForConfig)
}

func TestAccountValidation(t *testing.T) {
	config, client := baseTestSetup(t, false)

	config.When("GetGroup").Return(lib.GroupName{Group: "default-group", Account: "default-account"})
	config.When("GetIgnoreErr", "account").Return("")

	ctx := Context{}

	accountTests(t, &ctx, client, validateAccountForConfig)
}
