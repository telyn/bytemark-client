package cmd

import (
	"github.com/cheekybits/is"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

/*
 =====================
  Environment Helpers
 =====================
*/

func CleanEnv() {
	// apparently os.Unsetenv doesn't exist in the version of go I'm using
	os.Setenv("BIGV_CONFIG_DIR", "")
	os.Setenv("BIGV_USER", "")
	os.Setenv("BIGV_ACCOUNT", "")
	os.Setenv("BIGV_ENDPOINT", "")
	os.Setenv("BIGV_AUTH_ENDPOINT", "")
	os.Setenv("BIGV_DEBUG_LEVEL", "")
}

func JunkEnv() {
	os.Setenv("BIGV_CONFIG_DIR", "junk-env-config-dir")
	os.Setenv("BIGV_USER", "junk-env-user")
	os.Setenv("BIGV_ACCOUNT", "junk-env-account")
	os.Setenv("BIGV_ENDPOINT", "junk-env-endpoint")
	os.Setenv("BIGV_AUTH_ENDPOINT", "junk-env-auth-endpoint")
	os.Setenv("BIGV_DEBUG_LEVEL", "junk-env-debug-level")
}

func FixtureEnv() (fixture map[string]string) {
	fixture = map[string]string{

		"endpoint":      "https://fixture.env.localhost.local",
		"user":          "fixture-env-user",
		"account":       "fixture-env-account",
		"auth-endpoint": "https://fixture.env.auth.localhost.local",
		"debug-level":   "fixture-env-debug-level",
	}
	os.Setenv("BIGV_CONFIG_DIR", fixture["config-dir"])
	os.Setenv("BIGV_USER", fixture["user"])
	os.Setenv("BIGV_ACCOUNT", fixture["account"])
	os.Setenv("BIGV_ENDPOINT", fixture["endpoint"])
	os.Setenv("BIGV_AUTH_ENDPOINT", fixture["auth-endpoint"])
	os.Setenv("BIGV_DEBUG_LEVEL", fixture["debug-level"])
	return fixture
}

/*
 ===================
  Directory Helpers
 ===================
*/

func CleanDir() (name string) {
	dir, err := ioutil.TempDir("", "bigv-client-test")
	if err != nil {
		panic("Couldn't create test dir.")
	}

	return dir

}

func JunkDir() (name string) {
	junk := map[string]string{
		"endpoint":      "https://junk.dir.localhost.local",
		"user":          "junk-dir-user",
		"account":       "junk-dir-account",
		"auth-endpoint": "https://junk.dir.auth.localhost.local",
		"debug-level":   "junk-dir-debug-level",
	}

	dir, err := ioutil.TempDir("", "bigv-client-test")
	if err != nil {
		panic("Couldn't create test dir.")
	}

	for name, value := range junk {
		ioutil.WriteFile(filepath.Join(dir, name), []byte(value), 0600)
	}
	return dir
}

func FixtureDir() (dir string, fixture map[string]string) {
	fixture = map[string]string{
		"endpoint":      "https://fixture.dir.localhost.local",
		"user":          "fixture-dir-user",
		"account":       "fixture-dir-account",
		"auth-endpoint": "https://fixture.dir.auth.localhost.local",
		"debug-level":   "fixture-dir-debug-level",
	}

	dir, err := ioutil.TempDir("", "bigv-client-test")
	if err != nil {
		panic("Couldn't create test dir.")
	}

	for name, value := range fixture {
		ioutil.WriteFile(filepath.Join(dir, name), []byte(value), 0600)
	}
	return dir, fixture
}

/*
 =========================
  Environment-based Tests
 =========================
*/

func TestConfigDefaultConfigDir(t *testing.T) {
	is := is.New(t)

	CleanEnv()

	config, err := NewConfig("", nil)
	if err != nil {
		t.Fatalf(err.Error())
	}
	expected := filepath.Join(os.Getenv("HOME"), "/.go-bigv")
	is.Equal(expected, config.Dir)
}

func TestConfigEnvConfigDir(t *testing.T) {
	is := is.New(t)

	CleanEnv()

	expected := "/tmp"
	os.Setenv("BIGV_CONFIG_DIR", expected)

	config, err := NewConfig("", nil)
	if err != nil {
		t.Fatalf(err.Error())
	}
	is.Equal(expected, config.Dir)
}

func TestConfigPassedConfigDir(t *testing.T) {
	is := is.New(t)

	JunkEnv()

	expected := "/home"
	config, err := NewConfig(expected, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}
	is.Equal(expected, config.Dir)
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
	dir := CleanDir()

	config, err := NewConfig(dir, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	is.Equal("https://uk0.bigv.io", config.GetIgnoreErr("endpoint"))
	is.Equal("https://auth.bytemark.co.uk", config.GetIgnoreErr("auth-endpoint"))

	is.Equal("", config.GetIgnoreErr("user"))
	is.Equal("", config.GetIgnoreErr("account"))

	os.RemoveAll(dir)
}

func TestConfigDefaultsWithEnvUser(t *testing.T) {
	is := is.New(t)

	CleanEnv()
	dir := CleanDir()

	expected := "test-username"
	os.Setenv("BIGV_USER", expected)

	config, err := NewConfig(dir, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	is.Equal("https://uk0.bigv.io", config.GetIgnoreErr("endpoint"))
	is.Equal("https://auth.bytemark.co.uk", config.GetIgnoreErr("auth-endpoint"))

	v, err := config.GetV("user")
	is.Nil(err)
	is.Equal("user", v.Name)
	is.Equal(expected, v.Value)
	is.Equal("ENV BIGV_USER", v.Source)

	v, err = config.GetV("account")
	is.Nil(err)
	is.Equal("account", v.Name)
	is.Equal(expected, v.Value)
	is.Equal("ENV BIGV_USER", v.Source)

	os.RemoveAll(dir)
}

func TestConfigDefaultsFixtureEnv(t *testing.T) {
	is := is.New(t)

	fixture := FixtureEnv()
	dir := CleanDir()

	config, err := NewConfig(dir, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	is.Equal(fixture["endpoint"], config.GetIgnoreErr("endpoint"))
	is.Equal(fixture["auth-endpoint"], config.GetIgnoreErr("auth-endpoint"))
	is.Equal(fixture["user"], config.GetIgnoreErr("user"))
	is.Equal(fixture["account"], config.GetIgnoreErr("account"))
	is.Equal(fixture["debug-level"], config.GetIgnoreErr("debug-level"))
	os.RemoveAll(dir)
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
	dir, fixture := FixtureDir()

	config, err := NewConfig(dir, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	is.Equal(fixture["endpoint"], config.GetIgnoreErr("endpoint"))
	is.Equal(fixture["auth-endpoint"], config.GetIgnoreErr("auth-endpoint"))
	is.Equal(fixture["user"], config.GetIgnoreErr("user"))
	is.Equal(fixture["account"], config.GetIgnoreErr("account"))
	is.Equal(fixture["debug-level"], config.GetIgnoreErr("debug-level"))

	os.RemoveAll(dir)
}

/*
 ============
  Flag Tests
 ============
*/

func testFlagsWithArgs(args []string) (*Config, error) {
	CleanEnv()
	dir := CleanDir()
	flags := MakeCommonFlagSet()

	flags.Parse(args)
	return NewConfig(dir, flags)

}

func TestMainFlags(t *testing.T) {
	is := is.New(t)

	config, err := testFlagsWithArgs([]string{"--help", "--force"})
	is.Nil(err)

	for _, v := range configVars {
		// if this line ever fails then either configVars is out of date, GetDefault is out of date, or something weird has happened

		vv, err := config.GetV(v)
		is.Nil(err)
		is.Equal(config.GetDefault(v), vv)
	}

	config, err = testFlagsWithArgs([]string{"--user=test-user", "-account=test-account", "--endpoint", "example.com"})

	v, err := config.GetV("user")
	is.Nil(err)
	is.Equal(ConfigVar{"user", "test-user", "FLAG user"}, v)

	v, err = config.GetV("account")
	is.Nil(err)
	is.Equal(ConfigVar{"account", "test-account", "FLAG account"}, v)

	v, err = config.GetV("endpoint")
	is.Nil(err)
	is.Equal(ConfigVar{"endpoint", "example.com", "FLAG endpoint"}, v)

}

/*
 ===========
  Set Tests
 ===========
*/

func TestConfigSet(t *testing.T) {
	is := is.New(t)

	CleanEnv()
	dir, fixture := FixtureDir()
	config, err := NewConfig(dir, nil)
	if err != nil {
		t.Fatalf(err.Error())
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
	os.RemoveAll(dir)
}

func TestConfigSetPersistent(t *testing.T) {
	is := is.New(t)

	CleanEnv()
	dir, fixture := FixtureDir()
	config, err := NewConfig(dir, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	is.Equal(fixture["endpoint"], config.GetIgnoreErr("endpoint"))
	is.Equal(fixture["auth-endpoint"], config.GetIgnoreErr("auth-endpoint"))
	is.Equal(fixture["user"], config.GetIgnoreErr("user"))
	is.Equal(fixture["account"], config.GetIgnoreErr("account"))
	is.Equal(fixture["debug-level"], config.GetIgnoreErr("debug-level"))

	config.SetPersistent("user", "test-user", "TEST")
	fixture["user"] = "test-user"

	is.Equal(fixture["endpoint"], config.GetIgnoreErr("endpoint"))
	is.Equal(fixture["auth-endpoint"], config.GetIgnoreErr("auth-endpoint"))
	is.Equal(fixture["user"], config.GetIgnoreErr("user"))
	is.Equal(fixture["account"], config.GetIgnoreErr("account"))
	is.Equal(fixture["debug-level"], config.GetIgnoreErr("debug-level"))

	CleanEnv() // in case for some wacky reason I write to the environment
	//create a new config (blanking the memo) to test the file in the directory has changed.
	config2, err := NewConfig(dir, nil)
	is.Nil(err)

	is.Equal(fixture["endpoint"], config2.GetIgnoreErr("endpoint"))
	is.Equal(fixture["auth-endpoint"], config2.GetIgnoreErr("auth-endpoint"))
	is.Equal(fixture["user"], config2.GetIgnoreErr("user"))
	is.Equal(fixture["account"], config2.GetIgnoreErr("account"))
	is.Equal(fixture["debug-level"], config2.GetIgnoreErr("debug-level"))

	os.RemoveAll(dir)
}
