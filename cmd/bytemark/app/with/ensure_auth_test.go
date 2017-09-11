package with

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
)

func TestEnsureAuth(t *testing.T) {
	tt := []struct {
		InputUsername             string
		InputPassword             string
		Input2FA                  string
		AuthWithCredentialsErrors []error
		Factors                   []string
		ExpectedError             bool
	}{
		{
			InputUsername:             "input-user",
			InputPassword:             "input-pass",
			AuthWithCredentialsErrors: []error{nil},
			ExpectedError:             false,
		},
		{
			InputUsername:             "input-user",
			InputPassword:             "input-pass",
			AuthWithCredentialsErrors: []error{fmt.Errorf("{}")},
			ExpectedError:             true,
		},
		{
			InputUsername:             "input-user",
			InputPassword:             "input-pass",
			Input2FA:                  "123456",
			AuthWithCredentialsErrors: []error{fmt.Errorf("Missing 2FA"), nil}, // 2nd error as nil tests success with 2FA login
			Factors:                   []string{"2fa"},
			ExpectedError:             false,
		},
		{
			InputUsername:             "input-user",
			InputPassword:             "input-pass",
			Input2FA:                  "123456",
			AuthWithCredentialsErrors: []error{fmt.Errorf("Missing 2FA"), fmt.Errorf("Invalid token")}, // 2nd error tests failure with 2FA token
			ExpectedError:             true,
		},
		{
			InputUsername:             "input-user",
			InputPassword:             "input-pass",
			Input2FA:                  "123456",
			AuthWithCredentialsErrors: []error{fmt.Errorf("Missing 2FA"), nil}, // 2nd error as nil means success with 2FA token
			Factors:                   []string{"missing-2fa-factor"},
			ExpectedError:             true,
		},
	}

	for _, test := range tt {
		_, c, _ := baseTestSetup(t, false)

		configDir, err := ioutil.TempDir("", "")
		if err != nil {
			t.Errorf("Unexpected error when setting up config temp directory: %v", err)
		}
		defer func() {
			removeErr := os.RemoveAll(configDir)
			if removeErr != nil {
				t.Errorf("Could not clean up config dir: %v", removeErr)
			}
		}()

		config, err := util.NewConfig(configDir)
		if err != nil {
			t.Errorf("Unexpected error when setting up config temp directory: %v", err)
		}

		// Pretending the input comes from terminal
		config.Set("user", test.InputUsername, "INTERACTION")
		config.Set("pass", test.InputPassword, "TESTING")
		config.Set("2fa-otp", test.Input2FA, "TESTING")

		c.When("AuthWithToken", "").Return(fmt.Errorf("Not logged in")).Times(1)

		credentials := auth3.Credentials{
			"username": test.InputUsername,
			"password": test.InputPassword,
			"validity": "1800",
		}

		c.When("AuthWithCredentials", credentials).Return(test.AuthWithCredentialsErrors[0]).Times(1)

		// We are supplying a 2FA token, so we want to test that flow
		if test.Input2FA != "" {
			credentials := auth3.Credentials{
				"username": test.InputUsername,
				"password": test.InputPassword,
				"validity": "1800",
				"2fa":      test.Input2FA,
			}
			c.When("AuthWithCredentials", credentials).Return(test.AuthWithCredentialsErrors[1]).Times(1) // Returns nil means success
		}

		// Only called if the login succeeded, so always return a token
		c.When("GetSessionToken").Return("test-token")

		c.When("GetSessionFactors").Return(test.Factors)

		err = EnsureAuth(c, config)
		if test.ExpectedError && err == nil {
			t.Error("Expecting EnsureAuth to error, but it didn't")
		} else if !test.ExpectedError && err != nil {
			t.Errorf("Not expecting EnsureAuth to error, but got %v", err)
		}

		if ok, err := c.Verify(); !ok {
			t.Fatal(err)
		}
	}
}
