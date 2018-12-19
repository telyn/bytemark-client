package brain_test

import (
	"fmt"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	brainRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func TestDeleteAPIKey(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		labelOrID  string
		statusCode int
		shouldErr  bool
	}{
		{
			name:       "success",
			labelOrID:  "9",
			statusCode: 200,
			shouldErr:  false,
		}, {
			name:       "successjumanji",
			id:         229,
			labelOrID:  "jumanji",
			statusCode: 200,
			shouldErr:  false,
		}, {
			name:       "error 500",
			labelOrID:  "25",
			statusCode: 500,
			shouldErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rts := testutil.RequestTestSpec{
				Method:        "DELETE",
				URL:           fmt.Sprintf("/api_keys/%s", test.id),
				Endpoint:      lib.BrainEndpoint,
				StatusCode:    test.statusCode,
				AssertRequest: assert.BodyString(""),
			}
			rts.Run(t, test.name, true, func(client lib.Client) {
				err := brainRequests.DeleteAPIKey(client, test.labelOrID)
				if err != nil && !test.shouldErr {
					t.Errorf("Unexpected error: %v", err)
				} else if err == nil && test.shouldErr {
					t.Error("Error expected but not returned")
				}
			})
		})
	}
}
