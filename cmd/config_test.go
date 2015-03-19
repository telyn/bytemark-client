package cmd

import (
	"github.com/stretchr/testify/assert"
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
 =====================
  Environment Helpers
 =====================
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
 ===================
  Assertion Helpers
 ===================
*/

func AssertConfigMatchesFixture(t *testing.T, fixture map[string]string, config *Config) {
	assert.Equal(t, fixture["endpoint"], config.Get("endpoint"))
	assert.Equal(t, fixture["auth-endpoint"], config.Get("auth-endpoint"))
	assert.Equal(t, fixture["user"], config.Get("user"))
	assert.Equal(t, fixture["account"], config.Get("account"))
	assert.Equal(t, fixture["debug-level"], config.Get("debug-level"))
}

/*
 =========================
  Environment-based Tests
 =========================
*/

func TestDefaultConfigDir(t *testing.T) {
	CleanEnv()

	config := NewConfig("", nil)
	expected := filepath.Join(os.Getenv("HOME"), "/.go-bigv")
	assert.Equal(t, expected, config.Dir)
}

func TestEnvConfigDir(t *testing.T) {
	CleanEnv()

	expected := "/tmp"
	os.Setenv("BIGV_CONFIG_DIR", expected)

	config := NewConfig("", nil)
	assert.Equal(t, expected, config.Dir)
}

func TestPassedConfigDir(t *testing.T) {
	JunkEnv()

	expected := "/home"
	config := NewConfig(expected, nil)
	assert.Equal(t, expected, config.Dir)
}

/*
 ================
  Defaulting Tests
 ================
*/
// Tests to make sure we get the right defaults given the environment

func TestConfigDefaultsCleanEnv(t *testing.T) {
	CleanEnv()
	dir := CleanDir()

	config := NewConfig(dir, nil)

	// TODO(telyn): Update me when we move to api.bigv.io
	assert.Equal(t, "https://uk0.bigv.io", config.Get("endpoint"))
	assert.Equal(t, "https://auth.bytemark.co.uk", config.Get("auth-endpoint"))

	assert.Equal(t, os.Getenv("USER"), config.Get("user"))
	assert.Equal(t, os.Getenv("USER"), config.Get("account"))

	os.RemoveAll(dir)
}

func TestConfigDefaultsWithEnvUser(t *testing.T) {
	CleanEnv()
	dir := CleanDir()

	expected := "test-username"
	os.Setenv("BIGV_USER", expected)

	config := NewConfig(dir, nil)

	assert.Equal(t, "https://uk0.bigv.io", config.Get("endpoint"))
	assert.Equal(t, "https://auth.bytemark.co.uk", config.Get("auth-endpoint"))
	assert.Equal(t, expected, config.Get("user"))
	assert.Equal(t, expected, config.Get("account"))

	os.RemoveAll(dir)
}

func TestConfigDefaultsFixtureEnv(t *testing.T) {
	fixture := FixtureEnv()
	dir := CleanDir()

	config := NewConfig(dir, nil)

	AssertConfigMatchesFixture(t, fixture, config)
}

/*
 =================
  Directory Tests
 =================
*/
// Tests to ensure that we get given the value in the directory if we didn't specify it as a flag

func TestConfigDir(t *testing.T) {
	JunkEnv()
	dir, fixture := FixtureDir()

	config := NewConfig(dir, nil)

	AssertConfigMatchesFixture(t, fixture, config)

	os.RemoveAll(dir)
}

/*
 ============
  Flag Tests
 ============
*/

// TODO(telyn): no i don't want to write flag tests today
