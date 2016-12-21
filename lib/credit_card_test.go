package lib

import (
	"encoding/json"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/spp"
	"github.com/cheekybits/is"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
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

func TestGetSPPToken(t *testing.T) {
	client, servers, err := mkTestClientAndServers(t, Handlers{
		billing: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/api/v1/accounts/spp_token" || req.Method != "POST" {
				t.Fatalf("Unexpected HTTP request to billing%s", req.URL.Path)
			}
			arr := req.Header["Authorization"]
			if arr == nil || len(arr) == 0 {
				t.Errorf("Call to /accounts/spp_token must be authed in GetSPPToken")
			} else if arr[0] != "Token token=working-auth-token" {
				t.Errorf("GetSPPToken must be authed with a regular token, but got '%s'", arr[0])

			} else {
				_, err := w.Write([]byte("test-spp-token"))
				if err != nil {
					t.Fatal(err)
				}
			}
		}),
	})
	defer servers.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	cc := spp.CreditCard{Number: "4343"}

	token, err := client.GetSPPToken(cc, nil)
	if err != nil {
		t.Fatal(err)
	}
	if token != "test-spp-token" {
		t.Errorf("Expecting test-spp-token, got '%s'", token)
	}

}

func TestGetSPPTokenWithAccount(t *testing.T) {
	is := is.New(t)
	client, servers, err := mkTestClientAndServers(t, Handlers{
		billing: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/api/v1/accounts/spp_token" || req.Method != "POST" {
				t.Fatalf("Unexpected HTTP request to billing%s", req.URL.Path)
			}
			arr := req.Header["Authorization"]
			if arr != nil || len(arr) != 0 {
				t.Errorf("No authentication should be provided by GetSPPTokenWithAccount")
			} else {
				if req.Body == nil {
					t.Error("Request Body must not be nil in GetSPPTokenWithAccount")
				}
				body, err := ioutil.ReadAll(req.Body)
				if err != nil {
					t.Fatal(err)
				}

				var tr sppTokenRequest
				// check that it parses into an account
				err = json.Unmarshal(body, &tr)
				if err != nil {
					t.Fatal(err)
				}
				if tr.Owner == nil {
					t.Fatalf("TestGetSPPTokenWithAccount shouldn't have nil account. for real content: \r\n%s\r\n", string(body))
				}
				is.Equal("Melanie", tr.Owner.FirstName)
				is.Equal("Ownersdottir", tr.Owner.LastName)
				is.Equal("4343", tr.CardEnding)

				_, err = w.Write([]byte("test-spp-token"))
				if err != nil {
					t.Fatal(err)
				}
			}
		}),
	})
	defer servers.Close()

	if err != nil {
		t.Fatal(err)
	}

	cc := spp.CreditCard{Number: "4343"}
	token, err := client.GetSPPToken(cc, &billing.Person{
		FirstName: "Melanie",
		LastName:  "Ownersdottir",
	})
	if err != nil {
		t.Fatal(err)
	}
	if token != "test-spp-token" {
		t.Errorf("Expecting test-spp-token, got '%s'", token)
	}

}

func TestCreateCreditCard(t *testing.T) {
	is := is.New(t)
	client, servers, err := mkTestClientAndServers(t, Handlers{
		billing: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/api/v1/accounts/spp_token" || req.Method != "POST" {
				t.Fatalf("Unexpected HTTP request to billing%s", req.URL.Path)
			}
			arr := req.Header["Authorization"]
			if arr == nil || len(arr) == 0 {
				t.Errorf("Call to /accounts/spp_token must be authed in GetSPPToken")
			} else if arr[0] != "Token token=working-auth-token" {
				t.Errorf("GetSPPToken must be authed with a regular token, but got '%s'", arr[0])

			} else {
				_, err := w.Write([]byte("test-spp-token"))
				if err != nil {
					t.Fatal(err)
				}
			}
		}),
		spp: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/card.ref" || req.Method != "POST" {

				t.Fatalf("Unexpected HTTP request to spp%s", req.URL.Path)
			}
			arr := req.Header["Authorization"]
			if arr == nil || len(arr) == 0 {
				t.Errorf("Call to /card must be authed in GetSPPToken")
			} else if arr[0] != "Token token=test-spp-token" {
				t.Errorf("CreateCreditCard must be authed with an spp token, but got '%s'", arr[0])

			} else {
				if req.Body == nil {
					t.Error("Request Body must not be nil in GetSPPTokenWithAccount")
				}
				body, err := ioutil.ReadAll(req.Body)
				if err != nil {
					t.Fatal(err)
				}

				values, err := url.ParseQuery(string(body))
				if err != nil {
					t.Fatal(err)
				}

				is.Equal(testCard.Name, values["name"][0])
				is.Equal(testCard.Number, values["account_number"][0])
				is.Equal(testCard.Expiry, values["expiry"][0])
				is.Equal(testCard.CVV, values["cvv"][0])
				is.Equal(testCard.Street, values["street"][0])
				is.Equal(testCard.City, values["city"][0])
				is.Equal(testCard.Postcode, values["postcode"][0])
				is.Equal(testCard.Country, values["country"][0])
			}
		}),
	})
	defer servers.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	ref, err := client.CreateCreditCard(&testCard)

	if err != nil {
		t.Fatal(err)
	}

	if ref != "" {
		t.Errorf("card ref was not cool-card-ref, was '%s' instead.", ref)
	}
}
