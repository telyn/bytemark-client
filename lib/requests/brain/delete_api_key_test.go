package brain_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
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
		apiKeys    brain.APIKeys
	}{
		{
			name:       "success",
			id:         9,
			labelOrID:  "9",
			statusCode: 200,
			shouldErr:  false,
		}, {
			name:       "successjumanji",
			id:         229,
			labelOrID:  "jumanji",
			statusCode: 200,
			shouldErr:  false,
			apiKeys: brain.APIKeys{{
				ID:    442,
				Label: "not-jumanji",
			}, {
				ID:    229,
				Label: "jumanji",
			}},
		}, {
			name:       "error 500",
			id:         25,
			labelOrID:  "25",
			statusCode: 500,
			shouldErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rts := testutil.RequestTestSpec{
				MuxHandlers: &testutil.MuxHandlers{
					Brain: testutil.Mux{
						fmt.Sprintf("/api_keys/%d", test.id): func(wr http.ResponseWriter, r *http.Request) {
							assert.BodyString("")(t, "", r)
							assert.Method("DELETE")(t, "", r)
							wr.WriteHeader(test.statusCode)
						},
						"/api_keys": func(wr http.ResponseWriter, r *http.Request) {
							assert.Method("GET")(t, "", r)
							bs, err := json.Marshal(test.apiKeys)
							if err != nil {
								t.Fatalf("json marshal failed: %s", err)
							}
							_, err = wr.Write(bs)
							if err != nil {
								t.Fatalf("http write failed %s", err)
							}
						},
					},
				},
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
