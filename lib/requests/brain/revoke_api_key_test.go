package brain_test

import (
	"fmt"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	brainRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func TestRevokeAPIKey(t *testing.T) {
	tests := []struct {
		name            string
		id              int
		requestExpected map[string]interface{}
		statusCode      int
		shouldErr       bool
	}{
		{
			name: "success",
			id:   9,
			requestExpected: map[string]interface{}{
				"expires_at": "00:00:00",
			},
			statusCode: 200,
			shouldErr:  false,
		}, {
			name: "error 500",
			id:   25,
			requestExpected: map[string]interface{}{
				"expires_at": "00:00:00",
			},
			statusCode: 500,
			shouldErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rts := testutil.RequestTestSpec{
				Method:        "PUT",
				URL:           fmt.Sprintf("/api_keys/%d", test.id),
				Endpoint:      lib.BrainEndpoint,
				StatusCode:    test.statusCode,
				AssertRequest: assert.BodyUnmarshalEqual(test.requestExpected),
			}
			rts.Run(t, test.name, true, func(client lib.Client) {
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
