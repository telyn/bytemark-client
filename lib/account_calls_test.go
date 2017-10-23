package lib_test

import (
	"net/http"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
	"github.com/cheekybits/is"
)

func TestGetAccount(t *testing.T) {
	is := is.New(t)
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		MuxHandlers: &testutil.MuxHandlers{
			Brain: testutil.Mux{
				"/accounts/account": func(w http.ResponseWriter, r *http.Request) {
					assert.All(
						assert.Method("GET"),
						assert.Auth(lib.TokenType(lib.BrainEndpoint)),
					)(t, testName, r)
					_, err := w.Write([]byte(`{
						"name": "account",
						"id": 1
					}`))
					if err != nil {
						t.Fatal(err)
					}
				},
			},
			Billing: testutil.Mux{
				"/api/v1/accounts": func(w http.ResponseWriter, r *http.Request) {
					assert.All(
						assert.Method("GET"),
						assert.Auth(lib.TokenType(lib.BillingEndpoint)),
					)(t, testName, r)
					_, err := w.Write([]byte(`[
						{
							"bigv_account_subscription": "account"
						},
						{ "bigv_account_subscription": "wrong-account" }
					]`))
					if err != nil {
						t.Fatal(err)
					}
				},
			},
		},
	}
	rts.Run(t, testutil.Name(0), true, func(client lib.Client) {
		t.Log("Testing an invalid account!")
		acc, err := client.GetAccount("invalid-account")
		is.NotNil(err)

		t.Log("Testing the default account!")
		acc, err = client.GetAccount("")
		is.Nil(err)
		is.Equal("account", acc.Name)
		is.Equal(1, acc.BrainID)

		t.Log("Testing a named account!")
		acc, err = client.GetAccount("account")
		is.Nil(err)
		if !acc.IsValid() {
			t.Fatal("account isn't valid")
		}
		is.Equal("account", acc.Name)
	})

}

func TestGetAccounts(t *testing.T) {
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		MuxHandlers: &testutil.MuxHandlers{
			Brain: testutil.Mux{
				"/accounts": func(w http.ResponseWriter, r *http.Request) {
					assert.All(
						assert.Method("GET"),
						assert.Auth(lib.TokenType(lib.BrainEndpoint)),
					)(t, testName, r)
					testutil.WriteJSON(t, w, []brain.Account{
						{ID: 1, Name: "account"},
						{ID: 10, Name: "dr-evil", Suspended: true},
					})
				},
			},
			Billing: testutil.Mux{
				"/api/v1/accounts": func(w http.ResponseWriter, r *http.Request) {
					assert.All(
						assert.Method("GET"),
						assert.Auth(lib.TokenType(lib.BillingEndpoint)),
					)(t, testName, r)
					testutil.WriteJSON(t, w, []billing.Account{
						{ID: 4032, Name: "dr-evil"},
					})
				},
			},
		},
	}

	rts.Run(t, testName, true, func(client lib.Client) {
		accs, err := client.GetAccounts()
		if err != nil {
			t.Fatalf("%s: %s", testName, err)
		}
		if 2 != len(accs) {
			t.Errorf("%s didn't return 2 accounts. returned %#v", testName, accs)
		}
		seenDrEvil := false
		seenAccount := false
		for _, a := range accs {
			if a.Name == "dr-evil" {
				seenDrEvil = true
				assert.Equal(t, testName, a.BillingID, 4032)
				assert.Equal(t, testName, a.BrainID, 10)
			} else if a.Name == "account" {
				seenAccount = true
				assert.Equal(t, testName, 0, a.BillingID)
				assert.Equal(t, testName, 1, a.BrainID)
			}
		}
		if !seenDrEvil {
			t.Errorf("%s didn't see dr-evil account", testName)
		}
		if !seenAccount {
			t.Errorf("%s didn't see account account", testName)
		}
	})

}

func TestDefaultAccount(t *testing.T) {
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		MuxHandlers: &testutil.MuxHandlers{
			Billing: testutil.Mux{
				"/api/v1/accounts": func(w http.ResponseWriter, r *http.Request) {
					assert.All(
						assert.Method("GET"),
						assert.Auth(lib.TokenType(lib.BillingEndpoint)),
					)(t, testName, r)
					testutil.WriteJSON(t, w, []billing.Account{
						{Name: "default-account"},
						{Name: "not-default-account"},
						{Name: "also-not-default-account"},
					})
				},
			},
			Brain: testutil.Mux{
				"/accounts/default-account": func(w http.ResponseWriter, r *http.Request) {
					assert.All(
						assert.Method("GET"),
						assert.Auth(lib.TokenType(lib.BrainEndpoint)),
					)(t, testName, r)
					testutil.WriteJSON(t, w, brain.Account{
						ID:   2402,
						Name: "default-account",
					})
				},
			},
		},
	}

	rts.Run(t, testName, true, func(client lib.Client) {
		acc, err := client.GetDefaultAccount()
		if err != nil {
			t.Fatalf("%#v\r\n", err)
		}

		assert.Equal(t, testName, "default-account", acc.Name)
		assert.Equal(t, testName, 2402, acc.BrainID)
	})
}

// TestDefaultAccountHasNoBigVSubscription relates to open-source/bytemark-client#33
func TestDefaultAccountHasNoBigVSubscription(t *testing.T) {
	testName := testutil.Name(0)

	rts := testutil.RequestTestSpec{
		MuxHandlers: &testutil.MuxHandlers{
			Brain: testutil.Mux{
				"/accounts/default-account": func(w http.ResponseWriter, r *http.Request) {
					assert.All(
						assert.Method("GET"),
						assert.Auth(lib.TokenType(lib.BrainEndpoint)),
					)(t, testName, r)
					testutil.WriteJSON(t, w, brain.Account{ID: 2402, Name: "default-account"})
				},
			},
			Billing: testutil.Mux{
				"/api/v1/accounts": func(w http.ResponseWriter, r *http.Request) {
					assert.All(
						assert.Method("GET"),
						assert.Auth(lib.TokenType(lib.BillingEndpoint)),
					)(t, testName, r)
					testutil.WriteJSON(t, w, []billing.Account{
						{ID: 469},
						{Name: "not-default-account"},
					})
				},
			},
		},
	}

	rts.Run(t, testName, true, func(client lib.Client) {
		acc, err := client.GetDefaultAccount()
		if err == nil {
			t.Fatalf("%s was expecting an err but didn't get one. Account received: %#v", testName, acc)
		}
		_, ok := err.(lib.NoDefaultAccountError)
		if !ok {
			t.Fatalf("%s got the wrong kind of error: %T %v", testName, err, err)
		}
	})
}

func TestRegisterNewAccount(t *testing.T) {
	testName := testutil.Name(0)

	ownertech := billing.Person{
		ID:        249385,
		FirstName: "Test",
		LastName:  "User",
		Username:  "test-user",
		Email:     "test@example.com",
		Address:   "Testing Street",
		City:      "Testropolis",
		Postcode:  "TE57 7ES",
		Country:   "TE",
		Phone:     "735773577357",
	}
	rts := testutil.RequestTestSpec{
		Method:        "POST",
		Endpoint:      lib.BillingEndpoint,
		URL:           "/api/v1/accounts",
		AssertRequest: assert.All(),
		Response: billing.Account{
			ID:               324567,
			Owner:            ownertech,
			TechnicalContact: ownertech,
			InvoiceTerms:     0,
			Name:             "test-user",
			PaymentMethod:    "Credit Card",
			CardReference:    "testxq12e",
			EarliestActivity: "2016-09-18",
		},
	}

	// ready to test!
	person := billing.Person{
		Username:  "test-user",
		Password:  "aaaa",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Address:   "Testing Street",
		City:      "Testropolis",
		Postcode:  "TE57 7ES",
		Country:   "TE",
		Phone:     "735773577357",
	}

	rts.Run(t, testName, false, func(client lib.Client) {
		newAcc, err := client.RegisterNewAccount(lib.Account{
			Owner:         person,
			CardReference: "testxq12e",
		})
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, testName, "test-user", newAcc.Owner.Username)
		assert.Equal(t, testName, "", newAcc.Owner.Password)
		assert.Equal(t, testName, "Test", newAcc.Owner.FirstName)
		assert.Equal(t, testName, "User", newAcc.Owner.LastName)
		assert.Equal(t, testName, "Testing Street", newAcc.Owner.Address)
		assert.Equal(t, testName, "Testropolis", newAcc.Owner.City)
		assert.Equal(t, testName, "TE57 7ES", newAcc.Owner.Postcode)
		assert.Equal(t, testName, "TE", newAcc.Owner.Country)
		assert.Equal(t, testName, "735773577357", newAcc.Owner.Phone)
	})

}
