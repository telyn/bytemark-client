package lib

import (
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/cheekybits/is"
	"net/http"
	"testing"
)

func TestGetAccount(t *testing.T) {
	is := is.New(t)
	client, authServer, brain, billing, err := mkTestClientAndServers(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/accounts/account" {
			_, err := w.Write([]byte(`{
			    "name": "account",
			    "id": 1
			}`))
			if err != nil {
				t.Fatal(err)
			}
		} else if req.URL.Path == "/accounts/invalid-account" {
			http.NotFound(w, req)
		} else {
			t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
		}

	}), http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/api/v1/accounts" {
			_, err := w.Write([]byte(`[
				{
				    "bigv_account_subscription": "account"
				},
				{ "bigv_account_subscription": "wrong-account" }
			]`))
			if err != nil {
				t.Fatal(err)
			}
		} else {
			t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
		}
	}))
	defer authServer.Close()
	defer brain.Close()
	defer billing.Close()

	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	acc, err := client.GetAccount("invalid-account")
	is.NotNil(err)

	acc, err = client.GetAccount("")
	is.Nil(err)
	is.Equal("account", acc.Name)
	is.Equal(1, acc.BrainID)

	acc, err = client.GetAccount("account")
	is.Nil(err)
	is.Equal("account", acc.Name)

}

func TestGetAccounts(t *testing.T) {
	is := is.New(t)
	client, authServer, brain, billing, err := mkTestClientAndServers(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/accounts" {
			_, err := w.Write([]byte(`[
			{
			    "name": "account",
			    "id": 1
			}, {
			    "name": "dr-evil",
			    "suspended": true,
			    "id": 10
			}
			]`))
			if err != nil {
				t.Fatal(err)
			}
		} else {
			t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
		}

	}), http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/api/v1/accounts" {
			_, err := w.Write([]byte(`[]`))
			if err != nil {
				t.Fatal(err)
			}
		} else {
			t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
		}
	}))
	defer authServer.Close()
	defer brain.Close()
	defer billing.Close()

	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	acc, err := client.GetAccounts()
	fmt.Print(err)
	is.Nil(err)
	is.Equal(2, len(acc))
	seenDrEvil := false
	seenAccount := false
	for _, a := range acc {
		if a.Name == "dr-evil" {
			seenDrEvil = true
		} else if a.Name == "account" {
			seenAccount = true
		}
	}
	is.Equal(true, seenDrEvil)
	is.Equal(true, seenAccount)

}

func TestDefaultAccount(t *testing.T) {
	is := is.New(t)
	client, authServer, brain, billing, err := mkTestClientAndServers(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/accounts/default-account" {
			_, err := w.Write([]byte(`
			{ "id": 2402, "suspended": false, "name": "default-account" }
			`))
			if err != nil {
				t.Fatal(err)
			}
		} else {
			t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
		}

	}),
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/api/v1/accounts" {
				_, err := w.Write([]byte(`[
				{ "bigv_account_subscription": "default-account" },
				{ "bigv_account_subscription": "not-default-account" },
				{ "bigv_account_subscription": "also-not-default-account" }
				]`))
				if err != nil {
					t.Fatal(err)
				}
			} else {
				t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
			}
		}))
	defer authServer.Close()
	defer brain.Close()
	defer billing.Close()

	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	acc, err := client.GetDefaultAccount()
	if err != nil {
		t.Fatalf("%#v\r\n", err)
	}

	is.Equal("default-account", acc.Name)
	is.Equal(2402, acc.BrainID)
}

// TestDefaultAccountHasNoBigVSubscription relates to open-source/bytemark-client#33
func TestDefaultAccountHasNoBigVSubscription(t *testing.T) {
	client, authServer, brain, billing, err := mkTestClientAndServers(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/accounts/default-account" {
			_, err := w.Write([]byte(`
			{ "id": 2402, "suspended": false, "name": "default-account" }
			`))
			if err != nil {
				t.Fatal(err)
			}
		} else {
			t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
		}

	}),
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/api/v1/accounts" {
				_, err := w.Write([]byte(`[
				{ },
				{ "bigv_account_subscription": "not-default-account" }
				]`))
				if err != nil {
					t.Fatal(err)
				}
			} else {
				t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
			}
		}))
	defer authServer.Close()
	defer brain.Close()
	defer billing.Close()

	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.GetDefaultAccount()
	_, ok := err.(NoDefaultAccountError)
	if !ok {
		t.Fatal(err)
	}
}

func TestRegisterNewAccount(t *testing.T) {
	is := is.New(t)
	client, authServer, brain, billingServer, err := mkTestClientAndServers(mkNilHandler(t), http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/api/v1/accounts" {
			// TODO check there's no auth header
			// TODO check it's a POST
			person := `{
					"id":249385,
					"firstname":"Test",
					"surname":"User",
					"username":"test-user",
					"email":"test@example.com",
					"email_backup":null,
					"address":"Testing Street",
					"city":"Testropolis",
					"statecounty":null,
					"postcode":"TE57 7ES",
					"country":"TE",
					"phone":"735773577357",
					"phonemobile":null,
					"organization":null,
					"division":null,
					"vatnumber":null
				}`
			w.Write([]byte(`
			{
				"id":324567,
				"owner": ` + person + `,
				"tech": ` + person + `,
				"invoice_terms":0,
				"bigv_account_subscription":"test-user",
				"payment_method":"Credit Card",
				"card_reference":"testxq12e",
				"earliest_activity":"2016-09-18"
			}
			`))
		} else {
			t.Fatalf("Unexpected HTTP request to %s", req.URL.Path)
		}

	}))
	defer authServer.Close()
	defer brain.Close()
	defer billingServer.Close()

	if err != nil {
		t.Fatal(err)
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

	newAcc, err := client.RegisterNewAccount(&Account{
		Owner:         &person,
		CardReference: "testxq12e",
	})
	if err != nil {
		t.Fatal(err)
	}
	is.NotNil(newAcc)
	is.Equal("test-user", newAcc.Owner.Username)
	is.Equal("", newAcc.Owner.Password)
	is.Equal("Test", newAcc.Owner.FirstName)
	is.Equal("User", newAcc.Owner.LastName)
	is.Equal("Testing Street", newAcc.Owner.Address)
	is.Equal("Testropolis", newAcc.Owner.City)
	is.Equal("TE57 7ES", newAcc.Owner.Postcode)
	is.Equal("TE", newAcc.Owner.Country)
	is.Equal("735773577357", newAcc.Owner.Phone)
}
