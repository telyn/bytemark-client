package auth

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
	"runtime/debug"
	"testing"

	auth3 "github.com/BytemarkHosting/auth-client"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/config"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	mock "github.com/maraino/go-mock"
	"github.com/urfave/cli"
)

type unexpect struct{}

func (u unexpect) Error() string {
	return "this auth call was unexpected."
}

type authState struct {
	// things to return from GetSession* calls.
	user    string
	token   string
	factors []string

	authWithTokenErr       error
	authWithCredentialsErr error
	impersonateErr         error
	expectedCredentials    auth3.Credentials
}

// authInput is the things that are in the Config.
// in real life they'd mostly be set as flags or just typed in to the prompt
type authInput struct {
	user string
	pass string
	// 2fa otp
	otp string
	// yubikey otp
	yubikey         string
	impersonate     string
	token           string
	promptResponses []string
}

func (ai authInput) Credentials() auth3.Credentials {
	credents := auth3.Credentials{
		"username": ai.user,
		"password": ai.pass,
		"validity": "1800",
	}
	if ai.otp != "" {
		credents["2fa"] = ai.otp
	}
	if ai.yubikey != "" {
		credents["yubikey"] = ai.yubikey
	}
	return credents
}

func stubClientAuth(t *testing.T, state *int, c *mocks.Client, input authInput, states []authState) {

	nextState := func(fn string) {
		*state++
		fmt.Printf("Moving to state #%d\n", *state)
		t.Logf("Moving to state #%d", *state)
		if *state >= len(states) {
			panic(fmt.Sprintf("%v was called unexpectedly - run out of states to move to.", fn))
		}
	}

	fmt.Println("stubClientAuth()")
	c.When("GetSessionUser").Call(func() string {
		return states[*state].user
	})
	c.When("GetSessionFactors").Call(func() []string {
		return states[*state].factors
	})
	c.When("GetSessionToken").Call(func() string {
		return states[*state].token
	})
	c.When("AuthWithToken", mock.Any).Call(func(token string) error {
		t.Logf("AuthWithToken called!")

		err := states[*state].authWithTokenErr
		if _, ok := err.(unexpect); ok {
			t.Fatalf("AuthWithToken should not have been called in state %v", *state)
		}

		if token != input.token {
			panic(fmt.Sprintf("token %q != input.token %q", token, input.token))
		}
		nextState("AuthWithToken")
		return err
	})

	c.When("AuthWithCredentials", mock.Any).Call(func(credents auth3.Credentials) error {
		t.Logf("AuthWithCredentials called!")

		err := states[*state].authWithCredentialsErr

		fmt.Printf("begin state %d authWithCredentialsErr: %v\n", *state, err)

		if _, ok := err.(unexpect); ok {
			t.Fatalf("AuthWithCredentials should not have been called in state %v", *state)
		}
		expected := states[*state].expectedCredentials
		if expected == nil {
			expected = input.Credentials()
		}
		if !reflect.DeepEqual(credents, expected) {
			panic(fmt.Sprintf("unexpected credentials. Expecting %#v, got %#v", expected, credents))
		}

		nextState("AuthWithCredentials")
		fmt.Printf("end state %d authWithCredentialsErr: %v\n", (*state)-1, err)
		return err
	})
	c.When("Impersonate", mock.Any).Call(func(user string) error {
		t.Logf("Impersonate called!")

		err := states[*state].impersonateErr

		if _, ok := err.(unexpect); ok {
			t.Fatalf("Impersonate should not have been called in state %v", *state)
		}

		if user != input.impersonate {
			panic(fmt.Sprintf("Impersonate called for unexpected user %q - expected %q", user, input.impersonate))
		}

		nextState("Impersonate")
		return err
	})
}

func stubPromptResponses(t *testing.T, counter *int, prompter *mocks.Prompter, resp []string) {
	nextResponse := func() string {
		if *counter == len(resp) {
			t.Fatal("ran out of prompt responses")
		}
		r := resp[*counter]
		*counter++
		return r
	}
	prompter.When("Prompt", mock.Any).Call(func(_ string) string {
		r := nextResponse()
		t.Logf("Prompt -> %v", r)
		return r
	})
	prompter.When("Ask", mock.Any).Call(func(_ string) (string, error) {
		r := nextResponse()
		t.Logf("Ask -> %v", r)
		if r == "ERR" {
			return "", fmt.Errorf("fake prompt error")
		}
		return r, nil
	})

}

// not using a mocked config cause Authenticator will sometimes unset config and i
// need that to actually happen, so that it can later do the right thing
// alternative would be keeping a fake config around with multiple states
// like the auth states but that seems like work
func setupAuthConfig(t *testing.T, input authInput) (conf config.Manager) {
	configDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Errorf("Unexpected error when setting up config temp directory: %v", err)
	}

	conf, err = config.New(configDir)
	if err != nil {
		t.Errorf("Unexpected error when setting up config temp directory: %v", err)
	}

	// Pretending the input comes from terminal so as not to cause a prompt
	conf.Set("user", input.user, "INTERACTION")
	conf.Set("pass", input.pass, "TESTING")
	conf.Set("impersonate", input.impersonate, "TESTING")
	conf.Set("2fa-otp", input.otp, "TESTING")
	if input.yubikey != "" {
		conf.Set("yubikey", "true", "TESTING")
		conf.Set("yubikey-otp", input.yubikey, "INTERACTION")
	}
	err = conf.SetPersistent("token", input.token, "TESTING")
	if err != nil {
		t.Errorf(fmt.Sprintf("Unexpected error when setting up config temp directory: %v", err))
	}
	return
}

func TestAuthenticate(t *testing.T) {
	// situations to detect:
	// no token:
	// - password auth succeed
	// - password auth fail 2fa-missing, 2fa works
	// - password auth fail 2fa-missing, 2fa does not work
	// - password auth fail
	// - yubikey auth succeed
	// - yubikey auth fail 2fa-missing, 2fa works
	// - yubikey auth fail 2fa-missing, 2fa does not work
	// - yubikey auth fail
	// authWithToken errors:
	// - password auth succeed
	// - password auth fail 2fa-missing, 2fa works
	// - password auth fail 2fa-missing, 2fa does not work
	// - password auth fail
	// - yubikey auth succeed
	// - yubikey auth fail 2fa-missing, 2fa works
	// - yubikey auth fail 2fa-missing, 2fa does not work
	// - yubikey auth fail
	// authWithToken succeeds:
	// - yubikey needed but not given
	//   - yubikey auth succeed
	//   - yubikey auth fail 2fa-missing, 2fa works
	//   - yubikey auth fail 2fa-missing, 2fa does not work
	//   - yubikey auth fail

	// and then everything again with impersonation succeeding
	// and then everything again with impersonation failing

	// + 1 more:
	// authWithToken succeeds:
	// - impersonate factor without asking for it causes complete retry
	//   - password auth succeed
	//   - password auth fail 2fa-missing, 2fa works
	//   - password auth fail 2fa-missing, 2fa does not work
	//   - password auth fail
	//   - yubikey needed but not given
	//     - yubikey auth succeed
	//     - yubikey auth fail 2fa-missing, 2fa works
	//     - yubikey auth fail 2fa-missing, 2fa does not work
	//     - yubikey auth fail

	tt := []struct {
		// name of the test in t.Run
		// cause there's loads to keep track of and lots of similar names,
		// name should be prefixed with some letters:
		// N - no token B - bad token G - good token
		// Y - yubikey provided
		// 2 - 2fa
		// I - impersonation requested
		name string

		input authInput

		// states is a list of SessionData - token, username, factors - which the client
		// will return from its GetSession* methods, and an error each for
		// Client.AuthWithToken, AuthWithCredentials and Impersonate methods.
		// If the error is unexpect{} then calls to that method will cause the test to fail,
		// otherwise it is returned from that method.
		// the initial state of each test is the first in the list.
		// every time AuthWithToken, AuthWithCredentials or Impersonate is called, the state counter increments by one
		states         []authState
		expectingError bool
	}{
		{
			name: "G ok when token valid",
			input: authInput{
				user:  "input-user",
				token: "valid-token",
			},
			states: []authState{
				{ // state 0
					authWithTokenErr:       nil,
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				}, {
					user:                   "input-user",
					factors:                []string{"password"},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: false,
		}, {
			name: "G ok when user in config is not same as user in session",
			input: authInput{
				user:  "input-user",
				token: "valid-token",
			},
			states: []authState{
				{ // state 0
					authWithTokenErr:       nil,
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				}, {
					user:                   "fox-mulder",
					factors:                []string{"password"},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: false,
		}, {
			name: "N credentials auth tries 3 times",
			input: authInput{
				user: "input-user",
				pass: "input-pass",
				promptResponses: []string{
					"input-user",
					"input-pass",
					"input-user",
					"input-pass",
				},
			},
			// start with a blank authSession
			states: []authState{
				{ // state 0 - credential login fails
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: fmt.Errorf("Authentication failed"),
					impersonateErr:         unexpect{},
				}, { // state 1 - credential login fails
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: fmt.Errorf("Authentication failed"),
					impersonateErr:         unexpect{},
				}, { // state 2 - credential login fails complete
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: fmt.Errorf("Authentication failed"),
					impersonateErr:         unexpect{},
				}, {
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: true,
		}, {
			name: "N credentials auth is ok when 3rd attempt succeeds",
			input: authInput{
				user: "input-user",
				pass: "input-pass",
				promptResponses: []string{
					"input-user",
					"input-pass",
					"input-user",
					"input-pass",
				},
			},
			// start with a blank authSession
			states: []authState{
				{ // state 0 - credential login fails
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: fmt.Errorf("Authentication failed"),
					impersonateErr:         unexpect{},
				}, { // state 1 - credential login fails
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: fmt.Errorf("Authentication failed"),
					impersonateErr:         unexpect{},
				}, { // state 2 - credential login fails complete
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: nil,
					impersonateErr:         unexpect{},
				}, {
					user:                   "input-user",
					factors:                []string{"password"},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: false,
		}, {
			name: "N ok when user and pass used",
			input: authInput{
				user: "input-user",
				pass: "input-pass",
			},
			// start with a blank authSession
			states: []authState{
				{ // state 0 - ahWithToken failed, now expecting credential login
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: nil,
					impersonateErr:         unexpect{},
				}, { // state 1 - authentication complete
					user:                   "input-user",
					token:                  "valid-token",
					factors:                []string{"password"},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: false,
		}, {
			name: "N stop when credentials returns some error",
			input: authInput{
				user: "input-user",
				pass: "input-pass",
			},
			states: []authState{
				{
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: fmt.Errorf("{}"),
					impersonateErr:         unexpect{},
				}, {
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: true,
		}, {
			name: "N2 ok when log in with 2fa flow",
			input: authInput{
				user: "input-user",
				pass: "input-pass",
				otp:  "input-2fa",
			},
			states: []authState{
				{ // state 0 - trying AuthWithCredentials without 2fa
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: fmt.Errorf("Missing 2FA"),
					impersonateErr:         unexpect{},
					expectedCredentials: auth3.Credentials{
						"username": "input-user",
						"password": "input-pass",
						"validity": "1800",
					},
				}, { // state 1 - trying AuthWithCredentials with 2fa
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: nil,
					impersonateErr:         unexpect{},
				}, { // state 2 - login successful, now we can get stuff!
					user:  "input-user",
					token: "valid-token",
					factors: []string{
						"password",
						"2fa",
					},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: false,
		}, {
			name: "N2 stop when invalid 2fa otp",
			input: authInput{
				user: "input-user",
				pass: "input-pass",
				otp:  "input-2fa",
			},
			states: []authState{
				{ // state 0 - trying AuthWithCredentials without 2fa
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: fmt.Errorf("Missing 2FA"),
					impersonateErr:         unexpect{},
					expectedCredentials: auth3.Credentials{
						"username": "input-user",
						"password": "input-pass",
						"validity": "1800",
					},
				}, { // state 1 - trying AuthWithCredentials with 2fa
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: fmt.Errorf("Invalid token"),
					impersonateErr:         unexpect{},
				}, { // state 2 - login failed
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: true,
		}, {
			name: "N2 stop when 2fa login gives missing-2fa-factor",
			input: authInput{
				user: "input-user",
				pass: "input-pass",
				otp:  "input-2fa",
			},
			states: []authState{
				{ // state 0 - trying AuthWithCredentials without 2fa
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: fmt.Errorf("Missing 2FA"),
					impersonateErr:         unexpect{},
					expectedCredentials: auth3.Credentials{
						"username": "input-user",
						"password": "input-pass",
						"validity": "1800",
					},
				}, { // state 1 - trying AuthWithCredentials with 2fa
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: nil,
					impersonateErr:         unexpect{},
					factors:                []string{"password"},
				}, { // state 2 - login failed
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: true,
		}, {
			name: "NY ok when login with yubikey",
			input: authInput{
				user:    "input-user",
				pass:    "input-pass",
				yubikey: "input-yubikey",
			},
			// start with a blank authSession
			states: []authState{
				{ // state 0 - now expecting credential login
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: nil,
					impersonateErr:         unexpect{},
				}, { // state 1 - authentication complete
					user:                   "input-user",
					token:                  "valid-token",
					factors:                []string{"password", "yubikey"},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: false,
		}, {
			name: "B ok when token fails, user and pass used",
			input: authInput{
				user:  "input-user",
				pass:  "input-pass",
				token: "invalid-token",
			},
			// start with a blank authSession
			states: []authState{
				{ // state 0
					authWithTokenErr:       fmt.Errorf("invalid token"),
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				}, { // state 1 - authWithToken failed, now expecting credential login
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: nil,
					impersonateErr:         unexpect{},
				}, { // state 2 - authentication complete
					user:                   "input-user",
					token:                  "valid-token",
					factors:                []string{"password"},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: false,
		}, {
			name: "stop when AuthWithToken fails with url.Error",
			input: authInput{
				user:  "input-user",
				pass:  "input-pass",
				token: "token",
			},
			states: []authState{
				{
					authWithCredentialsErr: unexpect{},
					authWithTokenErr: &auth3.Error{
						Err: &url.Error{
							Op:  "GET",
							URL: "fake-url",
							Err: fmt.Errorf("pretending that there was some kind of URL parsing / transport error"),
						},
					},
					impersonateErr: unexpect{},
				}, {
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: true,
		}, {
			name: "B stop when bad token, credentials returns some error",
			input: authInput{
				user:  "input-user",
				pass:  "input-pass",
				token: "invalid-token",
			},
			states: []authState{
				{
					authWithTokenErr:       fmt.Errorf("invalid token"),
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				}, {
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: fmt.Errorf("{}"),
					impersonateErr:         unexpect{},
				}, {
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: true,
		}, {
			name: "B2 ok when bad token, log in with credentials using 2fa flow",
			input: authInput{
				user:  "input-user",
				pass:  "input-pass",
				otp:   "input-2fa",
				token: "invalid-token",
			},
			states: []authState{
				{ // state 0 - trying AuthWithToken
					authWithTokenErr:       fmt.Errorf("invalid token"),
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				}, { // state 1 - trying AuthWithCredentials without 2fa
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: fmt.Errorf("Missing 2FA"),
					impersonateErr:         unexpect{},
					expectedCredentials: auth3.Credentials{
						"username": "input-user",
						"password": "input-pass",
						"validity": "1800",
					},
				}, { // state 2 - trying AuthWithCredentials with 2fa
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: nil,
					impersonateErr:         unexpect{},
				}, { // state 3 - login successful, now we can get stuff!
					user:  "input-user",
					token: "valid-token",
					factors: []string{
						"password",
						"2fa",
					},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: false,
		}, {
			name: "B2 stop when invalid token, invalid 2fa otp",
			input: authInput{
				user:  "input-user",
				pass:  "input-pass",
				otp:   "input-2fa",
				token: "invalid-token",
			},
			states: []authState{
				{ // state 0 - trying AuthWithToken
					authWithTokenErr:       fmt.Errorf("invalid token"),
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				}, { // state 1 - trying AuthWithCredentials without 2fa
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: fmt.Errorf("Missing 2FA"),
					impersonateErr:         unexpect{},
					expectedCredentials: auth3.Credentials{
						"username": "input-user",
						"password": "input-pass",
						"validity": "1800",
					},
				}, { // state 2 - trying AuthWithCredentials with 2fa
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: fmt.Errorf("Invalid token"),
					impersonateErr:         unexpect{},
				}, { // state 3 - login failed
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: true,
		}, {
			name: "B2 stop when invalid token, 2fa login gives missing-2fa-factor",
			input: authInput{
				user:  "input-user",
				pass:  "input-pass",
				otp:   "input-2fa",
				token: "invalid-token",
			},
			states: []authState{
				{ // state 0 - trying AuthWithToken
					authWithTokenErr:       fmt.Errorf("invalid token"),
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				}, { // state 1 - trying AuthWithCredentials without 2fa
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: fmt.Errorf("Missing 2FA"),
					impersonateErr:         unexpect{},
					expectedCredentials: auth3.Credentials{
						"username": "input-user",
						"password": "input-pass",
						"validity": "1800",
					},
				}, { // state 2 - trying AuthWithCredentials with 2fa
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: nil,
					impersonateErr:         unexpect{},
					factors:                []string{"password"},
				}, { // state 3 - login failed
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: true,
			// END 2fa
		}, {
			// BEGIN yubikey
			name: "BY ok when invalid token, login with yubikey",
			input: authInput{
				user:    "input-user",
				pass:    "input-pass",
				yubikey: "input-yubikey",
				token:   "invalid-token",
			},
			// start with a blank authSession
			states: []authState{
				{ // state 0
					authWithTokenErr:       fmt.Errorf("invalid token"),
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				}, { // state 1 - authWithToken failed, now expecting credential login
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: nil,
					impersonateErr:         unexpect{},
				}, { // state 2 - authentication complete
					user:                   "input-user",
					token:                  "valid-token",
					factors:                []string{"password", "yubikey"},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: false,
		}, {
			name: "GY ok when yubikey requested, valid token has yubikey",
			input: authInput{
				user:    "input-user",
				yubikey: "input-yubikey",
				token:   "valid-token",
			},
			states: []authState{
				{ // state 0
					authWithTokenErr:       nil,
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				}, {
					user:                   "input-user",
					factors:                []string{"password", "yubikey"},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: false,
		}, {
			name: "GY ok when yubikey requested, valid token but has no yubikey, login with yubikey",
			input: authInput{
				user:    "input-user",
				pass:    "input-pass",
				yubikey: "input-yubikey",
				token:   "valid-token",
			},
			states: []authState{
				{ // state 0
					authWithTokenErr:       nil,
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
					factors:                []string{"password"},
				}, { // state 1 - authWithToken succeeded but had wrong factors
					// now expecting credential login
					user:                   "input-user",
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: nil,
					impersonateErr:         unexpect{},
				}, {
					user:                   "input-user",
					factors:                []string{"password", "yubikey"},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: false,
			// END yubikey set
		}, {
			// BEGIN no token yubikey + 2fa
			name: "NY2 ok when log in with 2fa flow",
			input: authInput{
				user:    "input-user",
				pass:    "input-pass",
				yubikey: "input-yubikey",
				otp:     "input-2fa",
			},
			states: []authState{
				{ // state 0 - trying AuthWithCredentials without 2fa
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: fmt.Errorf("Missing 2FA"),
					impersonateErr:         unexpect{},
					expectedCredentials: auth3.Credentials{
						"username": "input-user",
						"password": "input-pass",
						"yubikey":  "input-yubikey",
						"validity": "1800",
					},
				}, { // state 1 - trying AuthWithCredentials with 2fa
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: nil,
					impersonateErr:         unexpect{},
				}, { // state 2 - login successful, now we can get stuff!
					user:  "input-user",
					token: "valid-token",
					factors: []string{
						"password",
						"yubikey",
						"2fa",
					},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: false,
		}, {
			name: "NY2 stop when invalid 2fa otp",
			input: authInput{
				user:    "input-user",
				pass:    "input-pass",
				yubikey: "input-yubikey",
				otp:     "input-2fa",
			},
			states: []authState{
				{ // state 0 - trying AuthWithCredentials without 2fa
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: fmt.Errorf("Missing 2FA"),
					impersonateErr:         unexpect{},
					expectedCredentials: auth3.Credentials{
						"username": "input-user",
						"password": "input-pass",
						"yubikey":  "input-yubikey",
						"validity": "1800",
					},
				}, { // state 1 - trying AuthWithCredentials with 2fa
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: fmt.Errorf("Invalid token"),
					impersonateErr:         unexpect{},
				}, { // state 2 - login failed
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: true,
		}, {
			name: "NY2 stop when 2fa login gives missing-2fa-factor",
			input: authInput{
				user:    "input-user",
				pass:    "input-pass",
				yubikey: "input-yubikey",
				otp:     "input-2fa",
			},
			states: []authState{
				{ // state 0 - trying AuthWithCredentials without 2fa
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: fmt.Errorf("Missing 2FA"),
					impersonateErr:         unexpect{},
					expectedCredentials: auth3.Credentials{
						"username": "input-user",
						"password": "input-pass",
						"yubikey":  "input-yubikey",
						"validity": "1800",
					},
				}, { // state 1 - trying AuthWithCredentials with 2fa
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: nil,
					impersonateErr:         unexpect{},
					factors:                []string{"password"},
				}, { // state 2 - login failed
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: true,
		}, {
			name: "BY2 ok when bad token, log in with credentials using 2fa flow",
			input: authInput{
				user:    "input-user",
				pass:    "input-pass",
				yubikey: "input-yubikey",
				otp:     "input-2fa",
				token:   "invalid-token",
			},
			states: []authState{
				{ // state 0 - trying AuthWithToken
					authWithTokenErr:       fmt.Errorf("invalid token"),
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				}, { // state 1 - trying AuthWithCredentials without 2fa
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: fmt.Errorf("Missing 2FA"),
					impersonateErr:         unexpect{},
					expectedCredentials: auth3.Credentials{
						"username": "input-user",
						"password": "input-pass",
						"yubikey":  "input-yubikey",
						"validity": "1800",
					},
				}, { // state 2 - trying AuthWithCredentials with 2fa
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: nil,
					impersonateErr:         unexpect{},
				}, { // state 3 - login successful, now we can get stuff!
					user:  "input-user",
					token: "valid-token",
					factors: []string{
						"password",
						"yubikey",
						"2fa",
					},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: false,
		}, {
			name: "BY2 stop when invalid token, invalid 2fa otp",
			input: authInput{
				user:    "input-user",
				pass:    "input-pass",
				yubikey: "input-yubikey",
				otp:     "input-2fa",
				token:   "invalid-token",
			},
			states: []authState{
				{ // state 0 - trying AuthWithToken
					authWithTokenErr:       fmt.Errorf("invalid token"),
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				}, { // state 1 - trying AuthWithCredentials without 2fa
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: fmt.Errorf("Missing 2FA"),
					impersonateErr:         unexpect{},
					expectedCredentials: auth3.Credentials{
						"username": "input-user",
						"password": "input-pass",
						"yubikey":  "input-yubikey",
						"validity": "1800",
					},
				}, { // state 2 - trying AuthWithCredentials with 2fa
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: fmt.Errorf("Invalid token"),
					impersonateErr:         unexpect{},
				}, { // state 3 - login failed
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: true,
		}, {
			name: "BY2 stop when invalid token, 2fa login gives missing-2fa-factor",
			input: authInput{
				user:    "input-user",
				pass:    "input-pass",
				yubikey: "input-yubikey",
				otp:     "input-2fa",
				token:   "invalid-token",
			},
			states: []authState{
				{ // state 0 - trying AuthWithToken
					authWithTokenErr:       fmt.Errorf("invalid token"),
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				}, { // state 1 - trying AuthWithCredentials without 2fa
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: fmt.Errorf("Missing 2FA"),
					impersonateErr:         unexpect{},
					expectedCredentials: auth3.Credentials{
						"username": "input-user",
						"password": "input-pass",
						"yubikey":  "input-yubikey",
						"validity": "1800",
					},
				}, { // state 2 - trying AuthWithCredentials with 2fa
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: nil,
					impersonateErr:         unexpect{},
					factors:                []string{"password", "yubikey"},
				}, { // state 3 - login failed
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: true,
		}, {
			name: "GI impersonation already done",
			input: authInput{
				user:        "input-user",
				impersonate: "identity-theft-michael",
				token:       "valid-token",
			},
			states: []authState{
				{ // state 0
					authWithTokenErr:       nil,
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				}, {
					user:                   "identity-theft-michael",
					factors:                []string{"password", "impersonated"},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: false,
		}, {
			name: "G impersonation already done but didn't want to impersonate someone",
			input: authInput{
				user:  "input-user",
				pass:  "input-pass",
				token: "valid-token",
			},
			states: []authState{
				{ // state 0 - logging in with the token
					authWithTokenErr:       nil,
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				}, { // state 1 - huh, we're identity theft michael. relog with credents
					user:                   "identity-theft-michael",
					factors:                []string{"password", "impersonated"},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: nil,
					impersonateErr:         unexpect{},
				}, { // state 2 - success
					user:                   "input-user",
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: false,
		}, {
			name: "GI impersonation successful",
			input: authInput{
				user:        "input-user",
				impersonate: "identity-theft-michael",
				token:       "valid-token",
			},
			states: []authState{
				{ // state 0
					authWithTokenErr:       nil,
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				}, { // state 1 - authenticated as input user, impersonating as identity-theft-michael
					factors:                []string{"password"},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         nil,
				}, { // state 2 - authenticated as identity-theft-michael
					user:                   "identity-theft-michael",
					factors:                []string{"password", "impersonated"},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: false,
		}, {
			name: "GI impersonation fails causes error",
			input: authInput{
				user:        "input-user",
				impersonate: "identity-theft-michael",
				token:       "valid-token",
			},
			states: []authState{
				{ // state 0
					authWithTokenErr:       nil,
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				}, { // state 1 - authenticated as input user, impersonating as identity-theft-michael
					factors:                []string{"password"},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         fmt.Errorf("idk i guess you don't have permissions or something"),
				}, { // state 2 - authenticated as identity-theft-michael
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: true,
		}, {
			name: "GI impersonation already done but wanted to impersonate someone else",
			input: authInput{
				user:        "input-user",
				pass:        "input-pass",
				impersonate: "identity-theft-michael",
				token:       "valid-token",
			},
			states: []authState{
				{ // state 0 - logging in with the token
					authWithTokenErr:       nil,
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				}, { // state 1 - huh, we're identity theft svetlana. relog with credents
					user:                   "identity-theft-svetlana",
					factors:                []string{"password", "impersonated"},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: nil,
					impersonateErr:         unexpect{},
				}, { // state 2 - now we impersonate michael
					factors:                []string{"password"},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         nil,
				}, { // state 3 - now we are michael
					user:                   "identity-theft-michael",
					factors:                []string{"password", "impersonated"},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: false,
		}, {
			name: "GI impersonation already done but wanted to impersonate someone else and auth is kinda broken causes error",
			input: authInput{
				user:        "input-user",
				pass:        "input-pass",
				impersonate: "identity-theft-michael",
				token:       "valid-token",
			},
			states: []authState{
				{ // state 0 - logging in with the token
					authWithTokenErr:       nil,
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				}, { // state 1 - huh, we're identity theft svetlana. relog with credents
					user:                   "identity-theft-svetlana",
					factors:                []string{"password", "impersonated"},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: nil,
					impersonateErr:         unexpect{},
				}, { // state 2 - now we impersonate michael
					factors:                []string{"password"},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         nil,
				}, { // state 3 - now we are kayfabe because auth is being RATHER worrying
					user:                   "identity-theft-kayfabe",
					factors:                []string{"password", "impersonated"},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: true,
		}, {
			name: "GYI impersonation successful",
			input: authInput{
				user:        "input-user",
				pass:        "input-pass",
				yubikey:     "yubikey",
				impersonate: "identity-theft-michael",
			},
			states: []authState{
				{ // state 0
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: nil,
					impersonateErr:         unexpect{},
				}, { // state 1 - authenticated as input user, try impersonating as identity-theft-michael
					factors:                []string{"password", "yubikey"},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         nil,
				}, { // state 2 - authenticated as identity-theft-michael
					user:                   "identity-theft-michael",
					factors:                []string{"password", "impersonated"},
					authWithTokenErr:       unexpect{},
					authWithCredentialsErr: unexpect{},
					impersonateErr:         unexpect{},
				},
			},
			expectingError: false,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if p := recover(); p != nil {
					t.Fatalf("something panicked: %v\n%s", p, debug.Stack())
				}
			}()
			// ignoring the config from this cause we make a real config later.
			_, client, _ := testutil.BaseTestSetup(t, false, []cli.Command{})

			stateCounter := 0
			stubClientAuth(t, &stateCounter, client, test.input, test.states)

			config := setupAuthConfig(t, test.input)
			defer func() {
				removeErr := os.RemoveAll(config.ConfigDir())
				if removeErr != nil {
					t.Fatalf("Could not clean up config dir: %v", removeErr)
				}
			}()

			fmt.Println("attempting authentication")

			prompter := mocks.Prompter{}
			promptCounter := 0
			stubPromptResponses(t, &promptCounter, &prompter, test.input.promptResponses)

			err := Authenticator{
				client:       client,
				config:       config,
				prompter:     prompter,
				passPrompter: prompter,
			}.Authenticate()
			if test.expectingError && err == nil {
				t.Error("Expecting Authenticate to error, but it didn't")
			} else if !test.expectingError && err != nil {
				t.Errorf("Not expecting Authenticate to error, but got %v", err)
			}

			if stateCounter+1 < len(test.states) {
				t.Errorf("Only went through first %v of %v states", stateCounter+1, len(test.states))
			}

			if ok, err := client.Verify(); !ok {
				t.Fatal(err)
			}
		})
	}
}
