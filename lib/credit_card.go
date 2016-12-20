package lib

import (
	"bytes"
	"encoding/json"
	"github.com/BytemarkHosting/bytemark-client/lib/spp"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"net/url"
)

// GetSPPTokenWithAccount requests a token to use with bmbilling, passing an account object.
func (c *bytemarkClient) GetSPPTokenWithAccount(account Account) (token string, err error) {
	r, err := c.BuildRequestNoAuth("POST", BillingEndpoint, "/accounts/spp_token")
	if err != nil {
		return
	}

	js, err := json.Marshal(account.billingAccount())
	if err != nil {
		return "", err
	}
	_, res, err := r.Run(bytes.NewBuffer(js), nil)
	token = string(res)
	return
}

// GetSPPToken requests a token to use with SPP from bmbilling.
// If account is nil, authenticates against bmbilling.
func (c *bytemarkClient) GetSPPToken() (token string, err error) {
	r, err := c.BuildRequest("POST", BillingEndpoint, "/accounts/spp_token")
	if err != nil {
		return
	}
	_, res, err := r.Run(nil, nil)
	if err != nil {
		return
	}
	token = string(res)
	return
}

// CreateCreditCard creates a credit card on SPP using the given token. Tokens must be acquired by using GetSPPToken or GetSPPTokenWithAccount first.
func (c *bytemarkClient) CreateCreditCardWithToken(cc *spp.CreditCard, token string) (ref string, err error) {
	req, err := c.BuildRequestNoAuth("POST", SPPEndpoint, "/card.ref")
	req.sppToken = token
	if err != nil {
		return
	}
	values := url.Values{}
	values.Add("account_number", cc.Number)
	values.Add("name", cc.Name)
	values.Add("expiry", cc.Expiry)
	values.Add("cvv", cc.CVV)
	if cc.Street != "" {
		values.Add("street", cc.Street)
		values.Add("city", cc.City)
		values.Add("postcode", cc.Postcode)
		values.Add("country", cc.Country)
	}
	// prevent CC details and card reference being written to log
	// this is a bit of a sledgehammer
	// TODO make it not a sledgehammer somehow
	oldfile := log.LogFile
	log.LogFile = nil
	_, response, err := req.Run(bytes.NewBufferString(values.Encode()), nil)
	log.LogFile = oldfile

	return string(response), err
}

// CreateCreditCard creates a credit card on SPP. It uses GetSPPToken to get a token.
func (c *bytemarkClient) CreateCreditCard(cc *spp.CreditCard) (ref string, err error) {
	token, err := c.GetSPPToken()
	if err != nil {
		return
	}
	return c.CreateCreditCardWithToken(cc, token)

}
