package brain

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib/brain"
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
		userID      int
		expect      map[string]interface{}
		response    brain.APIKey
		responseErr error

		shouldErr bool
	}{
		{
			name: "neither user nor UserID set",
			spec: brain.APIKey{
				ExpiresAt: "2018-03-03T03:03:03Z",
			},
			shouldErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rts := testutil.RequestTestSpec{
				MuxHandlers: &testutil.MuxHandlers{
					Brain: testutil.Mux{
						"/api_keys": func(wr http.ResponseWriter, r *http.Request) {
							assert.Equal(t, test.name, "POST", r.Method)
							assert.BodyUnmarshalEqual(test.expect)(t, test.name, r)
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
					assert.Equal(t, test.name, "GET", r.Method)
					bytes, err := json.Marshal(brain.User{ID: test.userID})
					if err != nil {
						t.Fatal(err)
					}
					wr.Write(bytes)
				}
			}

		})
	}
}
