package brain

import (
	"fmt"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func TestRevokeAPIKey(t *testing.T) {
	tests := []struct {
		id              int
		requestExpected map[string]interface{}
		statusCode      int
		shouldErr       error
	}{
		{
			id: 9,
			requestExpected: map[string]interface{}{
				"id":         9,
				"expires_at": "00:00:00",
			},
			statusCode: 200,
			shouldErr:  false,
		}, {
			id: 25,
			requestExpected: map[string]interface{}{
				"id":         25,
				"expires_at": "00:00:00",
			},
			statusCode: 500,
			shouldErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rts := testutil.RequestTestSpec{
				Method:   "PUT",
				URL:      fmt.Sprintf("/api_keys/%d", test.id),
				Endpoint: lib.BrainEndpoint,
				Response: brain.APIKey{
					ID:        test.id,
					Label:     "jeff's cool api key for arctic exploration",
					APIKey:    "fake-api-key-whatever",
					ExpiresAt: "2018-11-12T00:00:00Z",
				},
				StatusCode:    test.statusCode,
				AssertRequest: assert.BodyUnmarshalEqual(test.requestExpected),
			}
			rts.Run(t, test.Name, true, func(client lib.Client) {
				err := brainRequests.RevokeAPIKey(client, test.id)
				if err != nil && !test.shouldErr {
					t.Errorf("Unexpected error: %v", err)
				} else if err == nil && test.shouldErr {
					t.Error("Error expected but not returned")
				}
			})
		})
	}
}
