package brain_test

import (
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	brainMethods "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func TestDeleteDiscByID(t *testing.T) {
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		Method:        "DELETE",
		Endpoint:      lib.BrainEndpoint,
		URL:           "/discs/666",
		AssertRequest: assert.QueryValue("purge", "true"),
	}
	rts.Run(t, testName, true, func(client lib.Client) {
		err := brainMethods.DeleteDiscByID(client, "666")
		if err != nil {
			t.Fatalf("%s err %s", testName, err)
		}
	})

}
