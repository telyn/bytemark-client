package brain_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	brainRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func TestCreateAPIKey(t *testing.T) {
	tests := []struct {
		name string
		user string
		spec brain.APIKey
		// if the user ID needs to be retrieved from the api (when user == "" and spec.UserID == 0)
		// then specify the user id to return here
		userID          int
		requestExpected map[string]interface{}
		response        brain.APIKey
		responseErr     error

		shouldErr bool
	}{
		{
			name: "neither user nor UserID set",
			spec: brain.APIKey{
				ExpiresAt: "2018-03-03T03:03:03Z",
			},
			shouldErr: true,
		}, {
			name: "both user and UserID set",
			user: "jeff",
			spec: brain.APIKey{
				UserID:    4123.0,
				ExpiresAt: "2018-03-03T03:03:03Z",
			},
			shouldErr: true,
		}, {
			name: "user set",
			user: "jeff",
			spec: brain.APIKey{
				ExpiresAt: "2019-04-04T04:04:04Z",
			},
			userID: 1111,
			requestExpected: map[string]interface{}{
				"user_id":    1111.0,
				"expires_at": "2019-04-04T04:04:04Z",
			},
			response: brain.APIKey{
				ID:        9,
				APIKey:    "fake-api-key",
				UserID:    1111,
				ExpiresAt: "2019-04-04T04:04:04Z",
			},
		}, {
			name: "UserID set",
			spec: brain.APIKey{
				UserID:    4123,
				ExpiresAt: "2018-05-05T05:05:05Z",
				Label:     "jeffs-cool-key-for-arctic-exploration",
			},
			requestExpected: map[string]interface{}{
				"user_id":    4123.0,
				"expires_at": "2018-05-05T05:05:05Z",
				"label":      "jeffs-cool-key-for-arctic-exploration",
			},
			response: brain.APIKey{
				ID:        12435,
				APIKey:    "this-key-be-fake",
				UserID:    4123,
				Label:     "jeffs-cool-key-for-arctic-exploration",
				ExpiresAt: "2018-05-05T05:05:05Z",
			},
		}, {
			name: "error!!",
			spec: brain.APIKey{
				UserID:    4123,
				ExpiresAt: "2018-05-05T05:05:05Z",
			},
			requestExpected: map[string]interface{}{
				"user_id":    4123.0,
				"expires_at": "2018-05-05T05:05:05Z",
			},
			shouldErr:   true,
			responseErr: errors.New("the turboencabulator couldn't effectively prevent side fumbling"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rts := testutil.RequestTestSpec{
				MuxHandlers: &testutil.MuxHandlers{
					Brain: testutil.Mux{
						"/api_keys": func(wr http.ResponseWriter, r *http.Request) {
							assert.Equal(t, "method", "POST", r.Method)
							assert.BodyUnmarshalEqual(test.requestExpected)(t, "request", r)

							if test.responseErr != nil {
								wr.WriteHeader(500)
								wr.Write([]byte(test.responseErr.Error()))
							}
							bytes, err := json.Marshal(test.response)
							if err != nil {
								t.Fatal(err)
							}
							wr.Write(bytes)
						},
					},
				},
			}
			if test.userID != 0 {
				rts.MuxHandlers.Brain["/users/"+test.user] = func(wr http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "", "GET", r.Method)
					bytes, err := json.Marshal(brain.User{ID: test.userID})
					if err != nil {
						t.Fatal(err)
					}
					wr.Write(bytes)
				}
			}

			rts.Run(t, "", true, func(client lib.Client) {
				apikey, err := brainRequests.CreateAPIKey(client, test.user, test.spec)
				if err != nil && !test.shouldErr {
					t.Errorf("Unexpected error: %v", err)
				} else if err == nil && test.shouldErr {
					t.Error("Error expected but not returned")
				}

				if !test.shouldErr {
					assert.Equal(t, "response", test.response, apikey)
				}
			})
		})
	}
}
