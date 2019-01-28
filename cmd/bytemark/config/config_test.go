package config

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/cheekybits/is"
)

/*
 =====================
  Environment Helpers
 =====================
*/

func CleanEnv() {
	// apparently os.Unsetenv doesn't exist in the version of go I'm using
	_ = os.Setenv("BM_CONFIG_DIR", "")
	_ = os.Setenv("BM_USER", "")
	_ = os.Setenv("BM_ACCOUNT", "")
	_ = os.Setenv("BM_ENDPOINT", "")
	_ = os.Setenv("BM_AUTH_ENDPOINT", "")
	_ = os.Setenv("BM_DEBUG_LEVEL", "")
}

func JunkEnv() {
	junk := map[string]string{
		"endpoint":      "https://junk.env.localhost.local",
		"user":          "junk-env-user",
		"account":       "junk-env-account",
		"auth-endpoint": "https://junk.env.auth.localhost.local",
		"debug-level":   "junk-env-debug-level",
	}
	MakeEnvFromFixture(junk)
}

func FixtureEnv() (fixture map[string]string) {
	fixture = map[string]string{

		"endpoint":      "https://fixture.env.localhost.local",
		"user":          "fixture-env-user",
		"account":       "fixture-env-account",
		"auth-endpoint": "https://fixture.env.auth.localhost.local",
		"debug-level":   "fixture-env-debug-level",
	}
	return MakeEnvFromFixture(fixture)
}

func MakeEnvFromFixture(fixture map[string]string) (fx map[string]string) {
	_ = os.Setenv("BM_CONFIG_DIR", fixture["config-dir"])
	_ = os.Setenv("BM_USER", fixture["user"])
	_ = os.Setenv("BM_ACCOUNT", fixture["account"])
	_ = os.Setenv("BM_ENDPOINT", fixture["endpoint"])
	_ = os.Setenv("BM_AUTH_ENDPOINT", fixture["auth-endpoint"])
	_ = os.Setenv("BM_DEBUG_LEVEL", fixture["debug-level"])
	return fixture
}

/*
 ===================
  Directory Helpers
 ===================
*/

func CleanDir() (dirName string, err error) {
	return ioutil.TempDir("", "bytemark-client-test")
}

func FixtureDir() (dir string, fixture map[string]string, err error) {
	fixture = map[string]string{
		"endpoint":      "https://fixture.dir.localhost.local",
		"user":          "fixture-dir-user",
		"account":       "fixture-dir-account",
		"auth-endpoint": "https://fixture.dir.auth.localhost.local",
		"debug-level":   "fixture-dir-debug-level",
	}

	return MakeDirFromFixture(fixture)
}

// MakeDirFromFixture takes a map[variable]value and writes it out as a config-dir.
func MakeDirFromFixture(fixture map[string]string) (dir string, fx map[string]string, err error) {
	fx = fixture
	dir, err = ioutil.TempDir("", "bytemark-client-test")
	if err != nil {
		return
	}
	for name, value := range fixture {
		err = ioutil.WriteFile(filepath.Join(dir, name), []byte(value), 0600)
		if err != nil {
			return
		}
	}
	return
}

/*
 =========================
  Environment-based Tests
 =========================
*/

func TestConfigDefaultConfigDir(t *testing.T) {
	is := is.New(t)

	CleanEnv()

	config, err := New("")
	if err != nil {
		t.Fatal(err)
	}
	expected := filepath.Join(os.Getenv("HOME"), "/.bytemark")
	is.Equal(expected, config.ConfigDir())
}

func TestConfigEnvConfigDir(t *testing.T) {
	is := is.New(t)

	CleanEnv()

	expected := "/tmp"
	_ = os.Setenv("BM_CONFIG_DIR", expected)

	config, err := New("")
	if err != nil {
		t.Fatal(err)
	}
	is.Equal(expected, config.ConfigDir())
}

func TestConfigPassedConfigDir(t *testing.T) {
	is := is.New(t)

	JunkEnv()

	expected := "/home"
	config, err := New(expected)
	if err != nil {
		t.Fatal(err)
	}
	is.Equal(expected, config.ConfigDir())
}

/*
 ================
  Defaulting Tests
 ================
*/
// Tests to make sure we get the right defaults given the environment

func TestConfigConfigDefaultsCleanEnv(t *testing.T) {
	is := is.New(t)

	CleanEnv()

	dir, err := CleanDir()
	if err != nil {
		t.Fatal(err)
	}

	config, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}

	is.Equal("https://uk0.bigv.io", config.GetIgnoreErr("endpoint"))
	is.Equal("https://auth.bytemark.co.uk", config.GetIgnoreErr("auth-endpoint"))

	is.Equal(os.Getenv("USER"), config.GetIgnoreErr("user"))
	is.Equal("", config.GetIgnoreErr("account"))

	// nbd if we leave a load of files lying about in a temp folder.
	_ = os.RemoveAll(dir)
}

func TestConfigDefaultsWithEnvUser(t *testing.T) {
	is := is.New(t)

	CleanEnv()
	dir, err := CleanDir()
	if err != nil {
		t.Fatal(err)
	}

	expected := "test-username"
	_ = os.Setenv("BM_USER", expected)

	config, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}

	is.Equal("https://uk0.bigv.io", config.GetIgnoreErr("endpoint"))
	is.Equal("https://auth.bytemark.co.uk", config.GetIgnoreErr("auth-endpoint"))

	v, err := config.GetV("user")
	is.Nil(err)
	is.Equal("user", v.Name)
	is.Equal(expected, v.Value)
	is.Equal("ENV BM_USER", v.Source)

	v, err = config.GetV("account")
	is.Nil(err)
	is.Equal("account", v.Name)
	is.Equal("", v.Value)
	is.Equal("CODE", v.Source)

	_ = os.RemoveAll(dir)
}

func TestConfigDefaultsFixtureEnv(t *testing.T) {
	is := is.New(t)

	fixture := FixtureEnv()
	dir, err := CleanDir()
	if err != nil {
		t.Fatal(err)
	}

	config, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}

	is.Equal(fixture["endpoint"], config.GetIgnoreErr("endpoint"))
	is.Equal(fixture["auth-endpoint"], config.GetIgnoreErr("auth-endpoint"))
	is.Equal(fixture["user"], config.GetIgnoreErr("user"))
	is.Equal(fixture["account"], config.GetIgnoreErr("account"))
	is.Equal(fixture["debug-level"], config.GetIgnoreErr("debug-level"))
	_ = os.RemoveAll(dir)
}

/*
 =================
  Directory Tests
 =================
*/
// Tests to ensure that we get given the value in the directory if we didn't specify it as a flag

func TestConfigDir(t *testing.T) {
	is := is.New(t)

	JunkEnv()
	dir, fixture, err := FixtureDir()
	if err != nil {
		t.Fatal(err)
	}

	config, err := New(dir)
	if err != nil {
		t.Fatal(err.Error())
	}

	is.Equal(fixture["endpoint"], config.GetIgnoreErr("endpoint"))
	is.Equal(fixture["auth-endpoint"], config.GetIgnoreErr("auth-endpoint"))
	is.Equal(fixture["user"], config.GetIgnoreErr("user"))
	is.Equal(fixture["account"], config.GetIgnoreErr("account"))
	is.Equal(fixture["debug-level"], config.GetIgnoreErr("debug-level"))

	_ = os.RemoveAll(dir)
}

/*
 ===========
  Set Tests
 ===========
*/

func TestConfigSet(t *testing.T) {
	is := is.New(t)

	CleanEnv()
	dir, fixture, err := FixtureDir()
	if err != nil {
		t.Fatal(err)
	}
	config, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}

	is.Equal(fixture["endpoint"], config.GetIgnoreErr("endpoint"))
	is.Equal(fixture["auth-endpoint"], config.GetIgnoreErr("auth-endpoint"))
	is.Equal(fixture["user"], config.GetIgnoreErr("user"))
	is.Equal(fixture["account"], config.GetIgnoreErr("account"))
	is.Equal(fixture["debug-level"], config.GetIgnoreErr("debug-level"))

	config.Set("user", "test-user", "TEST")
	fixture["user"] = "test-user"

	is.Equal(fixture["endpoint"], config.GetIgnoreErr("endpoint"))
	is.Equal(fixture["auth-endpoint"], config.GetIgnoreErr("auth-endpoint"))
	is.Equal(fixture["user"], config.GetIgnoreErr("user"))
	is.Equal(fixture["account"], config.GetIgnoreErr("account"))
	is.Equal(fixture["debug-level"], config.GetIgnoreErr("debug-level"))
	_ = os.RemoveAll(dir)
}

func TestConfigSetPersistent(t *testing.T) {
	is := is.New(t)

	CleanEnv()
	dir, fixture, err := FixtureDir()
	if err != nil {
		t.Fatal(err)
	}
	config, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}

	is.Equal(fixture["endpoint"], config.GetIgnoreErr("endpoint"))
	is.Equal(fixture["auth-endpoint"], config.GetIgnoreErr("auth-endpoint"))
	is.Equal(fixture["user"], config.GetIgnoreErr("user"))
	is.Equal(fixture["account"], config.GetIgnoreErr("account"))
	is.Equal(fixture["debug-level"], config.GetIgnoreErr("debug-level"))

	err = config.SetPersistent("user", "test-user", "TEST")
	if err != nil {
		t.Errorf("Couldn't SetPersistent(user, test-user, TEST): %v", err)
		t.Fail()
	}
	fixture["user"] = "test-user"

	is.Equal(fixture["endpoint"], config.GetIgnoreErr("endpoint"))
	is.Equal(fixture["auth-endpoint"], config.GetIgnoreErr("auth-endpoint"))
	is.Equal(fixture["user"], config.GetIgnoreErr("user"))
	is.Equal(fixture["account"], config.GetIgnoreErr("account"))
	is.Equal(fixture["debug-level"], config.GetIgnoreErr("debug-level"))

	CleanEnv()
	//create a new config (blanking the memo) to test the file in the directory has changed.
	config2, err := New(dir)
	is.Nil(err)

	is.Equal(fixture["endpoint"], config2.GetIgnoreErr("endpoint"))
	is.Equal(fixture["auth-endpoint"], config2.GetIgnoreErr("auth-endpoint"))
	is.Equal(fixture["user"], config2.GetIgnoreErr("user"))
	is.Equal(fixture["account"], config2.GetIgnoreErr("account"))
	is.Equal(fixture["debug-level"], config2.GetIgnoreErr("debug-level"))

	_ = os.RemoveAll(dir)
}

func TestConfigCorrectDefaultingAccountAndUserBug14038(t *testing.T) {
	is := is.New(t)

	envfixture := MakeEnvFromFixture(map[string]string{
		"user": "test-env-user",
	})
	dir, dirfixture, err := MakeDirFromFixture(map[string]string{
		"account": "test-fixture-account",
	})
	if err != nil {
		t.Fatal(err)
	}

	config, err := New(dir)
	is.Nil(err)

	v, err := config.GetV("account")
	is.Nil(err)
	is.Equal(v.Value, dirfixture["account"])

	fmt.Printf("%v\r\n", hex.EncodeToString([]byte(v.Source)))
	is.Equal(v.SourceType(), "FILE")
	is.Equal(v.SourceBaseName(), "account")
	is.NotEqual(v.Value, envfixture["user"])

}

func TestConfigUnset(t *testing.T) {
	is := is.New(t)

	dir, fx, err := FixtureDir()
	if err != nil {
		t.Fatal(err)
	}

	is.NotEqual("", fx["endpoint"])

	config, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}

	v, err := config.GetV("endpoint")
	is.Nil(err)
	is.Equal(fx["endpoint"], v.Value)
	is.Equal("FILE", v.SourceType())

	err = config.Unset("endpoint")
	is.Nil(err)

	v, err = config.GetV("endpoint")
	is.Nil(err)
	is.Equal("CODE", v.SourceType())

	err = config.Unset("endpoint")
	is.Nil(err)
}

func TestConfigEndpointOverrides(t *testing.T) {
	tests := []struct {
		Args             []string
		ExpectedEndpoint string
		ExpectedBilling  string
	}{{
		Args:             []string{},
		ExpectedEndpoint: "https://uk0.bigv.io",
		ExpectedBilling:  "https://bmbilling.bytemark.co.uk",
	}, {
		Args:             []string{"--endpoint", "https://int.bigv.io"},
		ExpectedEndpoint: "https://int.bigv.io",
		ExpectedBilling:  "",
	}}

	for i, test := range tests {
		fs := flag.NewFlagSet("test-config-endpoint-overrides", flag.ContinueOnError)
		fs.String("endpoint", "", "")
		fs.String("billing-endpoint", "", "")
		fs.String("auth-endpoint", "", "")
		fs.String("spp-endpoint", "", "")

		err := fs.Parse(test.Args)
		if err != nil {
			t.Fatal(err)
		}
		CleanEnv()
		dir, err := CleanDir()
		if err != nil {
			t.Fatal(err)
		}
		config, err := New(dir)
		if err != nil {
			t.Fatal(err)
		}

		config.ImportFlags(fs)
		if test.ExpectedEndpoint != config.GetIgnoreErr("endpoint") {
			t.Errorf("%d %q != %q", i, test.ExpectedEndpoint, config.GetIgnoreErr("endpoint"))
		}
		if test.ExpectedBilling != config.GetIgnoreErr("billing-endpoint") {
			t.Errorf("%d %q != %q", i, test.ExpectedBilling, config.GetIgnoreErr("billing-endpoint"))
		}

	}
}
