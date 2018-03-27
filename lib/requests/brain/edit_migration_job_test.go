package brain_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	brainMethods "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
	"github.com/BytemarkHosting/bytemark-client/lib/util"
)

func TestUpdateMigration(t *testing.T) {
	tests := []struct {
		id            int
		modifications brain.MigrationJobModification
		// expected      map[string]interface{}
		shouldErr bool
	}{
		{
			id: 1,
			modifications: brain.MigrationJobModification{
				Cancel: brain.MigrationJobLocations{
					Discs: []util.NumberOrString{"disc.1234"},
					Pools: []util.NumberOrString{"tail1-sata4"},
					Tails: []util.NumberOrString{"tail2"},
				},
				Options: brain.MigrationJobOptions{
					Priority: 10,
				},
			},
		}, //end of the test
	}
	for i, test := range tests {
		testName := testutil.Name(i)
		rts := testutil.RequestTestSpec{
			Method:   "PUT",
			Endpoint: lib.BrainEndpoint,
			URL:      fmt.Sprintf("/admin/migration_jobs/%d", test.id),
			AssertRequest: assert.Body(func(t *testing.T, testName string, body string) {
				var req brain.MigrationJobModification
				err := json.Unmarshal([]byte(body), &req)
				if err != nil {
					t.Fatalf("failed to unmarshal request body: %v", err)
				}
				assert.Equal(t, testName, req, test.modifications)
			}),
		}
		rts.Run(t, testName, true, func(client lib.Client) {
			err := brainMethods.EditMigrationJob(client, test.id, test.modifications)
			if test.shouldErr {
				assert.NotEqual(t, testName, nil, err)
			} else {
				assert.Equal(t, testName, nil, err)
			}

		})
	}

}

func TestCancelMigrationJob(t *testing.T) {
	tests := []struct {
		id        int
		expected  map[string]interface{}
		shouldErr bool
	}{
		{
			id: 1,
			expected: map[string]interface{}{
				"cancel": map[string]interface{}{
					"all": true,
				},
			},
		},
	}
	for i, test := range tests {
		testName := testutil.Name(i)
		rts := testutil.RequestTestSpec{
			Method:        "PUT",
			Endpoint:      lib.BrainEndpoint,
			URL:           fmt.Sprintf("/admin/migration_jobs/%d", test.id),
			AssertRequest: assert.BodyUnmarshalEqual(test.expected),
		}
		rts.Run(t, testName, true, func(client lib.Client) {
			err := brainMethods.CancelMigrationJob(client, test.id)
			if test.shouldErr {
				assert.NotEqual(t, testName, nil, err)
			} else {
				assert.Equal(t, testName, nil, err)
			}
		})
	}
}
