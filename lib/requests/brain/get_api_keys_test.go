package brain_test

import (
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	brainRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func TestGetAPIKeys(t *testing.T) {
	tests := []struct {
		name       string
		response   []brain.APIKey
		statusCode int
		shouldErr  bool
	}{
		{
			name: "empty array",
		},
		{
			name:       "http 500",
			statusCode: 500,
			shouldErr:  true,
		},
		{
			name: "some keys",
			response: []brain.APIKey{
				{
					ID:        6,
					UserID:    2152,
					Label:     "gitlab-autoscaling",
					APIKey:    "extremelyrandomdatahereuhhhhh7",
					ExpiresAt: "2019-03-21T00:24:45.0312Z",
					Privileges: brain.Privileges{
						{
							Username: "dr-gitlabotopis",
							APIKeyID: 6,
							GroupID:  2433,
							Level:    "group_admin",
						},
					},
				}, {
					ID:     912,
					UserID: 2505,
					Label:  "kubernetes-cloud-controller-manager",
					APIKey: "keykeykeykeykeynoonesgonnaguessthat",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rts := testutil.RequestTestSpec{
				Method:     "GET",
				Endpoint:   lib.BrainEndpoint,
				URL:        "/api_keys",
				StatusCode: test.statusCode,
				Response:   test.response,
			}
			rts.Run(t, test.name, true, func(client lib.Client) {
				keys, err := brainRequests.GetAPIKeys(client)
				if err != nil && !test.shouldErr {
					t.Errorf("Unexpected error: %v", err)
				} else if err == nil && test.shouldErr {
					t.Error("Error expected but not returned")
				}
				assert.Equal(t, "api keys", test.response, keys)
			})
		})
	}
}
