package lib_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/spp"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
	"github.com/cheekybits/is"
)

var testCard = spp.CreditCard{
	Number:   "4444444444444444",
	Name:     "Henrik v.Karthaltersburg",
	Expiry:   "4491",
	CVV:      "906",
	Street:   "Karthalterstrasse",
	City:     "Karthaltersburg",
	Postcode: "34098",
	Country:  "DE",
}

type sppTokenRequestWithRawOwner struct {
	Owner      json.RawMessage `json:"owner"`
	CardEnding string          `json:"card_ending"`
}

func TestGetSPPToken(t *testing.T) {
	var tokenRequest sppTokenRequestWithRawOwner
	rts := testutil.RequestTestSpec{
		Method:   "POST",
		Endpoint: lib.BillingEndpoint,
		URL:      "/api/v1/accounts/spp_token",
		Response: json.RawMessage(`{"token":"test-spp-token"}`),
		AssertRequest: assert.BodyUnmarshal(&tokenRequest, func(_ *testing.T, testName string) {
			t.Logf("%#v", tokenRequest)
			if len(tokenRequest.Owner) > 0 {
				t.Errorf("%s sent an owner to the server when it shouldn't have. Owner: %s", testName, string(tokenRequest.Owner))
			}
			if tokenRequest.CardEnding != "4343" {
				t.Errorf("%s incorrect card ending %s", testName, tokenRequest.CardEnding)
			}
		}),
	}
	cc := spp.CreditCard{Number: "52524343"}

	rts.Run(t, testutil.Name(0), true, func(client lib.Client) {
		token, err := client.GetSPPToken(cc, billing.Person{})
		if err != nil {
			t.Fatal(err)
		}
		if token != "test-spp-token" {
			t.Errorf("Expecting test-spp-token, got '%s'", token)
		}
	})
}

func TestGetSPPTokenWithAccount(t *testing.T) {
	is := is.New(t)
	var tokenRequest billing.SPPTokenRequest

	rts := testutil.RequestTestSpec{
		Method:   "POST",
		Endpoint: lib.BillingEndpoint,
		URL:      "/api/v1/accounts/spp_token",
		Response: json.RawMessage(`{"token":"test-spp-token"}`),
		AssertRequest: assert.BodyUnmarshal(&tokenRequest, func(_ *testing.T, testName string) {
			is.Equal(9020, tokenRequest.Owner.ID)
			is.Equal("melanie", tokenRequest.Owner.Username)
			is.Equal("Melanie", tokenRequest.Owner.FirstName)
			is.Equal("Ownersdottir", tokenRequest.Owner.LastName)
		}),
	}

	cc := spp.CreditCard{Number: "4343"}
	rts.Run(t, testutil.Name(0), false, func(client lib.Client) {
		token, err := client.GetSPPToken(cc, billing.Person{
			ID:        9020,
			Username:  "melanie",
			FirstName: "Melanie",
			LastName:  "Ownersdottir",
		})
		if err != nil {
			t.Fatal(err)
		}
		if token != "test-spp-token" {
			t.Errorf("Expecting test-spp-token, got '%s'", token)
		}
	})

}

func TestCreateCreditCard(t *testing.T) {
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		MuxHandlers: &testutil.MuxHandlers{
			Billing: testutil.Mux{
				"/api/v1/accounts/spp_token": func(w http.ResponseWriter, r *http.Request) {
					assert.All(
						assert.Method("POST"),
						assert.Auth(lib.TokenType(lib.BillingEndpoint)),
					)(t, testutil.Name(0), r)
					_, err := w.Write([]byte(`{"token":"test-spp-token"}`))
					if err != nil {
						t.Fatal(err)
					}
				},
			},
			SPP: testutil.Mux{
				"/card.ref": func(w http.ResponseWriter, r *http.Request) {
					assert.All(
						assert.Method("POST"),
						assert.Unauthed(),
						assert.BodyFormValue("token", "test-spp-token"),
						assert.BodyFormValue("name", testCard.Name),
						assert.BodyFormValue("account_number", testCard.Number),
						assert.BodyFormValue("expiry", testCard.Expiry),
						assert.BodyFormValue("cvv", testCard.CVV),
						assert.BodyFormValue("street", testCard.Street),
						assert.BodyFormValue("city", testCard.City),
						assert.BodyFormValue("postcode", testCard.Postcode),
						assert.BodyFormValue("country", testCard.Country),
					)(t, testName, r)
					w.Write([]byte("cool-card-ref"))
				},
			},
		},
	}

	rts.Run(t, testutil.Name(0), true, func(client lib.Client) {
		ref, err := client.CreateCreditCard(testCard)

		if err != nil {
			t.Fatal(err)
		}

		if ref != "cool-card-ref" {
			t.Errorf("card ref was not 'cool-card-ref', was '%s' instead.", ref)
		}
	})
}
