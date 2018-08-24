package brain_test

import (
	"encoding/json"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	brainMethods "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func TestGetDiscsOnTail(t *testing.T) {
	testName := testutil.Name(0)

	testDiscs := brain.Discs{{
		ID:    4,
		Label: "test_disc",
	}}

	rts := testutil.RequestTestSpec{
		Method:   "GET",
		URL:      "/admin/tails/123/discs",
		Endpoint: lib.BrainEndpoint,
		Response: json.RawMessage(`[{
		"id": 4,
		"label": "test_disc"
	    }]`),
	}

	rts.Run(t, testName, true, func(client lib.Client) {
		discs, err := brainMethods.GetDiscsOnTail(client, "123", "")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, testName, discs, testDiscs)
	})

	rts.AssertRequest = assert.QueryValue("at", "2018-08-21T15:00:00+0000")

	rts.Run(t, testName, true, func(client lib.Client) {
		discs, err := brainMethods.GetDiscsOnTail(client, "123", "2018-08-21T15:00:00+0000")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, testName, discs, testDiscs)
	})
}

func TestGetDiscsOnStoragePool(t *testing.T) {
	testName := testutil.Name(0)

	testDiscs := brain.Discs{{
		ID:    4,
		Label: "test_disc",
	}}

	rts := testutil.RequestTestSpec{
		Method:   "GET",
		URL:      "/admin/storage_pools/123/discs",
		Endpoint: lib.BrainEndpoint,
		Response: json.RawMessage(`[{
		"id": 4,
		"label": "test_disc"
	    }]`),
	}

	rts.Run(t, testName, true, func(client lib.Client) {
		discs, err := brainMethods.GetDiscsOnStoragePool(client, "123", "")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, testName, discs, testDiscs)
	})

	rts.AssertRequest = assert.QueryValue("at", "2018-08-21T15:00:00+0000")

	rts.Run(t, testName, true, func(client lib.Client) {
		discs, err := brainMethods.GetDiscsOnStoragePool(client, "123", "2018-08-21T15:00:00+0000")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, testName, discs, testDiscs)
	})
}
