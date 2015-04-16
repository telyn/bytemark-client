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

func TestDefaultConfigDir(t *testing.T) {
	is := is.New(t)

	CleanEnv()

	config := NewConfig("", nil)
	expected := filepath.Join(os.Getenv("HOME"), "/.go-bigv")
	is.Equal(expected, config.Dir)
}

func TestEnvConfigDir(t *testing.T) {
	is := is.New(t)

	CleanEnv()

	expected := "/tmp"
	os.Setenv("BIGV_CONFIG_DIR", expected)

	config := NewConfig("", nil)
	is.Equal(expected, config.Dir)
}

func TestPassedConfigDir(t *testing.T) {
	is := is.New(t)

	JunkEnv()

	expected := "/home"
	config := NewConfig(expected, nil)
	is.Equal(expected, config.Dir)
}

/*
 ================
  Defaulting Tests
 ================
*/
// Tests to make sure we get the right defaults given the environment

func TestConfigDefaultsCleanEnv(t *testing.T) {
	is := is.New(t)

	CleanEnv()
	dir := CleanDir()

	config := NewConfig(dir, nil)

	// TODO(telyn): Update me when we move to api.bigv.io
	is.Equal("https://uk0.bigv.io", config.Get("endpoint"))
	is.Equal("https://auth.bytemark.co.uk", config.Get("auth-endpoint"))

	is.Equal("", config.Get("user"))
	is.Equal("", config.Get("account"))

	os.RemoveAll(dir)
}

func TestConfigDefaultsWithEnvUser(t *testing.T) {
	is := is.New(t)

	CleanEnv()
	dir := CleanDir()

	expected := "test-username"
	os.Setenv("BIGV_USER", expected)

	config := NewConfig(dir, nil)

	is.Equal("https://uk0.bigv.io", config.Get("endpoint"))
	is.Equal("https://auth.bytemark.co.uk", config.Get("auth-endpoint"))
	is.Equal(expected, config.Get("user"))
	is.Equal(expected, config.Get("account"))

	os.RemoveAll(dir)
}

func TestConfigDefaultsFixtureEnv(t *testing.T) {
	is := is.New(t)

	fixture := FixtureEnv()
	dir := CleanDir()

	config := NewConfig(dir, nil)

	is.Equal(fixture["endpoint"], config.Get("endpoint"))
	is.Equal(fixture["auth-endpoint"], config.Get("auth-endpoint"))
	is.Equal(fixture["user"], config.Get("user"))
	is.Equal(fixture["account"], config.Get("account"))
	is.Equal(fixture["debug-level"], config.Get("debug-level"))
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

	config := NewConfig(dir, nil)

	is.Equal(fixture["endpoint"], config.Get("endpoint"))
	is.Equal(fixture["auth-endpoint"], config.Get("auth-endpoint"))
	is.Equal(fixture["user"], config.Get("user"))
	is.Equal(fixture["account"], config.Get("account"))
	is.Equal(fixture["debug-level"], config.Get("debug-level"))

	os.RemoveAll(dir)
}

/*
 ============
  Flag Tests
 ============
*/

// TODO(telyn): no i don't want to write flag tests today

/*
 ===========
  Set Tests
 ===========
*/

func TestSet(t *testing.T) {
	is := is.New(t)

	CleanEnv()
	dir, fixture := FixtureDir()
	config := NewConfig(dir, nil)

	is.Equal(fixture["endpoint"], config.Get("endpoint"))
	is.Equal(fixture["auth-endpoint"], config.Get("auth-endpoint"))
	is.Equal(fixture["user"], config.Get("user"))
	is.Equal(fixture["account"], config.Get("account"))
	is.Equal(fixture["debug-level"], config.Get("debug-level"))

	config.Set("user", "test-user")
	fixture["user"] = "test-user"

	is.Equal(fixture["endpoint"], config.Get("endpoint"))
	is.Equal(fixture["auth-endpoint"], config.Get("auth-endpoint"))
	is.Equal(fixture["user"], config.Get("user"))
	is.Equal(fixture["account"], config.Get("account"))
	is.Equal(fixture["debug-level"], config.Get("debug-level"))
	os.RemoveAll(dir)
}

func TestSetPersistent(t *testing.T) {
	is := is.New(t)

	CleanEnv()
	dir, fixture := FixtureDir()
	config := NewConfig(dir, nil)

	is.Equal(fixture["endpoint"], config.Get("endpoint"))
	is.Equal(fixture["auth-endpoint"], config.Get("auth-endpoint"))
	is.Equal(fixture["user"], config.Get("user"))
	is.Equal(fixture["account"], config.Get("account"))
	is.Equal(fixture["debug-level"], config.Get("debug-level"))

	config.SetPersistent("user", "test-user")
	fixture["user"] = "test-user"

	is.Equal(fixture["endpoint"], config.Get("endpoint"))
	is.Equal(fixture["auth-endpoint"], config.Get("auth-endpoint"))
	is.Equal(fixture["user"], config.Get("user"))
	is.Equal(fixture["account"], config.Get("account"))
	is.Equal(fixture["debug-level"], config.Get("debug-level"))

	CleanEnv() // in case for some wacky reason I write to the environment
	//create a new config (blanking the memo) to test the file in the directory has changed.
	config2 := NewConfig(dir, nil)

	is.Equal(fixture["endpoint"], config2.Get("endpoint"))
	is.Equal(fixture["auth-endpoint"], config2.Get("auth-endpoint"))
	is.Equal(fixture["user"], config2.Get("user"))
	is.Equal(fixture["account"], config2.Get("account"))
	is.Equal(fixture["debug-level"], config2.Get("debug-level"))

	os.RemoveAll(dir)
}
