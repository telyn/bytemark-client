package brain_test

import (
	"net/http"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	brainRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func TestDeleteUserAuthorizedKey(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		keyToDelete string

		// user returned from GET /users/username
		user         brain.User
		expectedKeys string

		shouldErr bool
	}{
		{
			name:        "no keys",
			username:    "jeff",
			keyToDelete: "ssh-rsa v-real-key",
			user: brain.User{
				Username:       "jeff",
				Email:          "jeff@jeff.jeff",
				AuthorizedKeys: brain.Keys{},
			},
			expectedKeys: "",
		}, {
			name:        "full key is same",
			username:    "grease",
			keyToDelete: "ssh-rsa v-real-key",
			user: brain.User{
				Username: "grease",
				Email:    "grease@jeff.jeff",
				AuthorizedKeys: brain.Keys{
					brain.Key{Key: "the old key\nssh-rsa v-real-key"},
				},
			},
			expectedKeys: "the old key",
		}, {
			name:        "delete by key comment",
			username:    "grease",
			keyToDelete: "comment",
			user: brain.User{
				Username: "grease",
				Email:    "grease@jeff.jeff",
				AuthorizedKeys: brain.Keys{
					brain.Key{Key: "the old key\nssh-rsa v-real-key comment"},
				},
			},
			expectedKeys: "the old key",
		}, {
			name:        "ambiguous comment",
			username:    "grease",
			keyToDelete: "comment",
			user: brain.User{
				Username: "grease",
				Email:    "grease@jeff.jeff",
				AuthorizedKeys: brain.Keys{
					brain.Key{Key: "the old key\nssh-rsa v-real-key comment\nssh-rsa another-key comment"},
				},
			},
			shouldErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user := make(map[string]interface{})
			url := "/users/" + test.username
			rts := testutil.RequestTestSpec{
				MuxHandlers: &testutil.MuxHandlers{
					Brain: testutil.Mux{
						url: func(wr http.ResponseWriter, r *http.Request) {
							switch r.Method {
							case "GET":
								testutil.WriteJSON(t, wr, test.user)
							case "PUT":
								assert.BodyUnmarshal(&user, func(_t *testing.T, _testName string) {
									assert.Equal(t, test.name, test.expectedKeys, user["authorized_keys"])
									assert.Equal(t, test.name, test.user.Username, user["username"])
									assert.Equal(t, test.name, test.user.Email, user["email"])
								})
							default:
								t.Errorf("Unexpected %s request to /users/jeff", r.Method)
							}
						},
					},
				},
			}
			// empty string for test name - t.Run takes care of that.
			rts.Run(t, "", true, func(client lib.Client) {
				err := brainRequests.DeleteUserAuthorizedKey(client, test.username, test.keyToDelete)
				if err != nil && !test.shouldErr {
					t.Errorf("Unexpected error: %v", err)
				} else if err == nil && test.shouldErr {
					t.Errorf("Error expected but not returned")
				}
			})
		})
	}
}
